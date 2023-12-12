package etcd

import (
	"context"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"runtime/debug"
	libconsts "social/lib/consts"
	liblog "social/lib/log"
	libutil "social/lib/util"
	"sync"
)

// Mgr 管理器
type Mgr struct {
	client                        *clientv3.Client
	kv                            clientv3.KV
	lease                         clientv3.Lease
	leaseGrantResponse            *clientv3.LeaseGrantResponse
	leaseKeepAliveResponseChannel <-chan *clientv3.LeaseKeepAliveResponse

	cancelFunc context.CancelFunc
	waitGroup  sync.WaitGroup // Stop 等待信号

	options *Options
}

// Handler etcd 处理数据
func (p *Mgr) Handler(key string, val string) error {
	return p.options.onFunc(key, val)
}

// Start 开始
func (p *Mgr) Start(ctx context.Context, opts ...*Options) error {
	p.options = mergeOptions(opts...)
	err := configure(p.options)
	if err != nil {
		return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}

	p.client, err = clientv3.New(clientv3.Config{
		Endpoints:   p.options.addrs,
		DialTimeout: *p.options.dialTimeout,
	})
	if err != nil {
		return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	// 获得kv api子集
	p.kv = clientv3.NewKV(p.client)
	// 申请一个lease 租约
	p.lease = clientv3.NewLease(p.client)
	// 申请一个ttl秒的租约
	p.leaseGrantResponse, err = p.lease.Grant(ctx, *p.options.ttl)
	if err != nil {
		return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	// 先删除
	for _, v := range p.options.kvSlice {
		_, err = p.Del(ctx, v.Key)
		if err != nil {
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
	}
	// 再添加
	for _, v := range p.options.kvSlice {
		_, err = p.PutWithLease(ctx, v.Key, v.Value)
		if err != nil {
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
	}
	return nil
}

// Run 租约续约
func (p *Mgr) Run(ctx context.Context) error {
	var err error
	p.leaseKeepAliveResponseChannel, err = p.lease.KeepAlive(ctx, p.leaseGrantResponse.ID)
	if err != nil {
		return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}

	p.waitGroup.Add(1)

	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc

	go func(ctx context.Context) {
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					liblog.PrintErr(libconsts.GoroutinePanic, err, debug.Stack())
				}
			}
			p.waitGroup.Done()
			liblog.PrintInfo(libconsts.GoroutineDone)
		}()
		for {
			select {
			case <-ctx.Done():
				liblog.PrintInfo(libconsts.GoroutineDone)
				return
			case leaseKeepAliveResponse, ok := <-p.leaseKeepAliveResponseChannel:
				liblog.PrintInfo(leaseKeepAliveResponse, ok)
				if leaseKeepAliveResponse != nil {
					continue
				}
				if ok {
					continue
				}
				p.abnormal(ctx)
				return
			}
		}
	}(ctxWithCancel)
	// 关注 服务
	if err = p.WatchPrefixSendIntoChan(ctxWithCancel, *p.options.watchServicePrefix); err != nil {
		return errors.Errorf("WatchPrefix err:%v %v", err, libutil.GetCodeLocation(1).String())
	}
	// 获取 服务
	if err = p.GetPrefixSendIntoChan(ctx, *p.options.watchServicePrefix); err != nil {
		return errors.Errorf("GetPrefix err:%v %v", err, libutil.GetCodeLocation(1).String())
	}
	// 关注 命令
	if err = p.WatchPrefixSendIntoChan(ctxWithCancel, *p.options.watchCommandPrefix); err != nil {
		return errors.Errorf("WatchPrefix err:%v %v", err, libutil.GetCodeLocation(1).String())
	}
	return nil
}

// Stop 停止
func (p *Mgr) Stop() error {
	if p.client != nil { // 删除
		for _, v := range p.options.kvSlice {
			_, err := p.Del(context.Background(), v.Key)
			if err != nil {
				liblog.PrintErr(err, libutil.GetCodeLocation(1).String())
				//	return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
			}
		}
		err := p.client.Close()
		if err != nil {
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
		p.client = nil
	}

	if p.cancelFunc != nil {
		p.cancelFunc()
		// 等待 goroutine退出.
		p.waitGroup.Wait()
		p.cancelFunc = nil
	}
	return nil
}

// Put 将一个键值对放入etcd中
func (p *Mgr) Put(ctx context.Context, key string, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	putResponse, err := p.kv.Put(ctx, key, value, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	return putResponse, nil
}

// PutWithLease 将一个键值对放入etcd中 WithLease 带ttl
func (p *Mgr) PutWithLease(ctx context.Context, key string, value string) (*clientv3.PutResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithLease(p.leaseGrantResponse.ID),
	}
	return p.Put(ctx, key, value, opts...)
}

// Del 删除
func (p *Mgr) Del(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	deleteResponse, err := p.kv.Delete(ctx, key, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	return deleteResponse, nil
}

// DelWithPrefix 删除键值 匹配的键值
func (p *Mgr) DelWithPrefix(ctx context.Context, keyPrefix string) (*clientv3.DeleteResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
	}
	return p.Del(ctx, keyPrefix, opts...)
}

// DelRange 按选项删除范围内的键值
func (p *Mgr) DelRange(ctx context.Context, startKeyPrefix string, endKeyPrefix string) (*clientv3.DeleteResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithFromKey(),
		clientv3.WithRange(endKeyPrefix),
	}
	return p.Del(ctx, startKeyPrefix, opts...)
}

// Watch 监视key
func (p *Mgr) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return p.client.Watch(ctx, key, opts...)
}

// WatchPrefix 监视以key为前缀的所有 key value
func (p *Mgr) WatchPrefix(ctx context.Context, key string) clientv3.WatchChan {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
	}
	return p.Watch(ctx, key, opts...)
}

// WatchPrefixSendIntoChan 监听key变化,放入 chan 中
func (p *Mgr) WatchPrefixSendIntoChan(ctx context.Context, preFix string) error {
	eventChan := p.WatchPrefix(ctx, preFix)
	go func(ctx context.Context) {
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					liblog.PrintErr(libconsts.GoroutinePanic, err, debug.Stack())
				}
			}
			liblog.PrintInfo(libconsts.GoroutineDone)
		}()
		for {
			select {
			case <-ctx.Done():
				liblog.PrintInfo(libconsts.GoroutineDone)
				return
			case v, ok := <-eventChan:
				liblog.PrintInfo(v, ok)
				Key := string(v.Events[0].Kv.Key)
				Value := string(v.Events[0].Kv.Value)
				p.options.outgoingEventChan <- &KV{
					Key:   Key,
					Value: Value,
				}
			}
		}
		//for v := range eventChan {
		//	Key := string(v.Events[0].Kv.Key)
		//	Value := string(v.Events[0].Kv.Value)
		//	p.options.outgoingEventChan <- &KV{
		//		Key:   Key,
		//		Value: Value,
		//	}
		//}
	}(ctx)
	return nil
}

// Get 检索键
func (p *Mgr) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	getResponse, err := p.kv.Get(ctx, key, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	return getResponse, nil
}

// GetPrefix 查找以key为前缀的所有 key value
func (p *Mgr) GetPrefix(ctx context.Context, key string) (*clientv3.GetResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
	}
	return p.Get(ctx, key, opts...)
}

// GetPrefixSendIntoChan  取得关心的前缀,放入 chan 中
func (p *Mgr) GetPrefixSendIntoChan(ctx context.Context, preFix string) error {
	getResponse, err := p.GetPrefix(ctx, preFix)
	if err != nil {
		return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	for _, v := range getResponse.Kvs {
		p.options.outgoingEventChan <- &KV{
			Key:   string(v.Key),
			Value: string(v.Value),
		}
	}
	return nil
}
