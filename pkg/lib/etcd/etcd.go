package etcd

import (
	"context"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"runtime/debug"
	libconstant "social/pkg/lib/constant"
	liblog "social/pkg/lib/log"
	libutil "social/pkg/lib/util"
	"sync"
	"time"
)

var (
	instance *Mgr
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *Mgr {
	once.Do(func() {
		instance = NewMgr()
	})
	return instance
}

// IsEnable 是否 开启
func IsEnable() bool {
	if instance == nil {
		return false
	}
	return instance.client != nil
}

func NewMgr() *Mgr {
	return new(Mgr)
}

// Mgr 管理器
type Mgr struct {
	client                        *clientv3.Client
	kv                            clientv3.KV
	lease                         clientv3.Lease
	leaseGrantResponse            *clientv3.LeaseGrantResponse
	leaseKeepAliveResponseChannel <-chan *clientv3.LeaseKeepAliveResponse

	cancelFunc context.CancelFunc
	waitGroup  sync.WaitGroup // Stop 等待信号

	options *options
}

// Handler etcd 处理数据
func (p *Mgr) Handler(key string, val string) error {
	return p.options.onFunc(key, val)
}

// Start 开始
func (p *Mgr) Start(ctx context.Context, opts ...*options) error {
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
	// 先删除,再添加
	for _, v := range p.options.kvSlice {
		_, err = p.DelWithPrefix(v.Key)
		if err != nil {
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
		_, err = p.PutWithLease(v.Key, v.Value)
		if err != nil {
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
	}
	return nil
}

// KeepAlive 更新租约
func (p *Mgr) KeepAlive(ctx context.Context) error {
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
					liblog.PrintErr(libconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			p.waitGroup.Done()
			liblog.PrintInfo(libconstant.GoroutineDone)
		}()

		for {
			select {
			case <-ctx.Done():
				liblog.PrintInfo(libconstant.GoroutineDone)
				return
			case leaseKeepAliveResponse, ok := <-p.leaseKeepAliveResponseChannel:
				liblog.PrintInfo(leaseKeepAliveResponse, ok)
				if leaseKeepAliveResponse != nil {
					continue
				}
				if ok {
					continue
				}
				// abnormal
				liblog.PrintErr("etcd lease KeepAlive died, retrying")
				go func(ctx context.Context) {
					defer func() {
						if libutil.IsRelease() {
							if err := recover(); err != nil {
								liblog.PrintErr(libconstant.Retry, libconstant.GoroutinePanic, err, debug.Stack())
							}
						}
						liblog.PrintInfo(libconstant.Retry, libconstant.GoroutineDone)
					}()
					if err := p.Stop(); err != nil {
						liblog.PrintInfo(libconstant.Retry, libconstant.Failure, err)
						return
					}
					if err := p.retryKeepAlive(ctx); err != nil {
						liblog.PrintErr(libconstant.Retry, libconstant.Failure, err)
						return
					}
				}(context.TODO())
				return
			}
		}
	}(ctxWithCancel)
	return nil
}

// Stop 停止
func (p *Mgr) Stop() error {
	if p.client != nil { // 删除
		for _, v := range p.options.kvSlice {
			_, err := p.DelWithPrefix(v.Key)
			if err != nil {
				liblog.PrintErr(libutil.GetCodeLocation(1).String())
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
func (p *Mgr) Put(key string, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	putResponse, err := p.kv.Put(context.TODO(), key, value, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	return putResponse, nil
}

// PutWithLease 将一个键值对放入etcd中 WithLease 带ttl
func (p *Mgr) PutWithLease(key string, value string) (*clientv3.PutResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithLease(p.leaseGrantResponse.ID),
	}
	return p.Put(key, value, opts...)
}

// Del 删除
func (p *Mgr) Del(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	deleteResponse, err := p.kv.Delete(context.TODO(), key, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	return deleteResponse, nil
}

// DelWithPrefix 删除键值 匹配的键值
func (p *Mgr) DelWithPrefix(keyPrefix string) (*clientv3.DeleteResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
	}
	return p.Del(keyPrefix, opts...)
}

// DelRange 按选项删除范围内的键值
func (p *Mgr) DelRange(startKeyPrefix string, endKeyPrefix string) (*clientv3.DeleteResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithFromKey(),
		clientv3.WithRange(endKeyPrefix),
	}
	return p.Del(startKeyPrefix, opts...)
}

// Watch 监视key
func (p *Mgr) Watch(key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return p.client.Watch(context.TODO(), key, opts...)
}

// WatchPrefix 监视以key为前缀的所有 key value
func (p *Mgr) WatchPrefix(key string) clientv3.WatchChan {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
	}
	return p.Watch(key, opts...)
}

// WatchPrefixSendIntoChan 监听key变化,放入 chan 中
func (p *Mgr) WatchPrefixSendIntoChan(preFix string) error {
	eventChan := p.WatchPrefix(preFix)
	go func() {
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					liblog.PrintErr(libconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			liblog.PrintInfo(libconstant.GoroutineDone)
		}()
		for v := range eventChan {
			Key := string(v.Events[0].Kv.Key)
			Value := string(v.Events[0].Kv.Value)
			p.options.outgoingEventChan <- &KV{
				Key:   Key,
				Value: Value,
			}
		}
	}()
	return nil
}

// Get 检索键
func (p *Mgr) Get(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	getResponse, err := p.kv.Get(context.TODO(), key, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	return getResponse, nil
}

// GetPrefix 查找以key为前缀的所有 key value
func (p *Mgr) GetPrefix(key string) (*clientv3.GetResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
	}
	return p.Get(key, opts...)
}

// GetPrefixSendIntoChan  取得关心的前缀,放入 chan 中
func (p *Mgr) GetPrefixSendIntoChan(preFix string) error {
	getResponse, err := p.GetPrefix(preFix)
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

// 多次重试 Start 和 KeepAlive
func (p *Mgr) retryKeepAlive(ctx context.Context) error {
	liblog.PrintfErr("renewing etcd lease, reconfiguring.grantLeaseMaxRetries:%v, grantLeaseIntervalSecond:%v",
		*p.options.grantLeaseMaxRetries, grantLeaseRetryDuration/time.Second)
	var failedGrantLeaseAttempts = 0
	for {
		if err := p.Start(ctx, p.options); err != nil {
			failedGrantLeaseAttempts++
			if *p.options.grantLeaseMaxRetries <= failedGrantLeaseAttempts {
				return errors.WithMessagef(err, "%v exceeded max attempts to renew etcd lease %v %v",
					libutil.GetCodeLocation(1), *p.options.grantLeaseMaxRetries, failedGrantLeaseAttempts)
			}
			liblog.PrintErr("error granting etcd lease, will retry.", err)
			time.Sleep(grantLeaseRetryDuration)
			continue
		} else {
		retryKeepAlive:
			// 续租
			if err = p.KeepAlive(ctx); err != nil {
				failedGrantLeaseAttempts++
				if *p.options.grantLeaseMaxRetries <= failedGrantLeaseAttempts {
					return errors.WithMessagef(err, "%v exceeded max attempts to renew etcd lease %v %v",
						libutil.GetCodeLocation(1), *p.options.grantLeaseMaxRetries, failedGrantLeaseAttempts)
				}
				liblog.PrintErr("error granting etcd lease, will retry.", err)
				time.Sleep(grantLeaseRetryDuration)
				goto retryKeepAlive
			} else {
				return nil
			}
		}
	}
}
