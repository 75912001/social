package etcd

import (
	"context"
	"runtime/debug"
	xrconstant "social/pkg/lib/constant"
	xrlog "social/pkg/lib/log"
	xrutil "social/pkg/lib/util"
	"sync"
	"time"

	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	instance *mgr
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *mgr {
	once.Do(func() {
		instance = new(mgr)
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

// mgr 管理器
type mgr struct {
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
func (p *mgr) Handler(key string, val string) error {
	return p.options.onFunc(key, val)
}

// Start 开始
func (p *mgr) Start(_ context.Context, opts ...*Options) error {
	p.options = mergeOptions(opts...)
	err := configure(p.options)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}

	p.client, err = clientv3.New(clientv3.Config{
		Endpoints:   p.options.addrs,
		DialTimeout: *p.options.dialTimeout,
	})
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	// 获得kv api子集
	p.kv = clientv3.NewKV(p.client)
	// 申请一个lease 租约
	p.lease = clientv3.NewLease(p.client)
	// 申请一个ttl秒的租约
	p.leaseGrantResponse, err = p.lease.Grant(context.TODO(), *p.options.ttl)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}

	// 先删除,再添加
	for _, v := range p.options.kvSlice {
		_, err = p.DelWithPrefix(v.Key)
		if err != nil {
			return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
		}
		_, err = p.PutWithLease(v.Key, v.Value)
		if err != nil {
			return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
		}
	}
	return nil
}

// Stop 停止
func (p *mgr) Stop() error {
	if p.client != nil { // 删除
		for _, v := range p.options.kvSlice {
			_, err := p.DelWithPrefix(v.Key)
			if err != nil {
				xrlog.PrintErr(xrutil.GetCodeLocation(1).String())
				//	return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
			}
		}

		err := p.client.Close()
		if err != nil {
			return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
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

func (p *mgr) getLeaseKeepAliveResponseChannel() <-chan *clientv3.LeaseKeepAliveResponse {
	return p.leaseKeepAliveResponseChannel
}

// 多次重试 Start 和 KeepAlive
func (p *mgr) retryKeepAlive(ctx context.Context) error {
	xrlog.PrintfErr("renewing etcd lease, reconfiguring.grantLeaseMaxRetries:%v, grantLeaseIntervalSecond:%v",
		*p.options.grantLeaseMaxRetries, grantLeaseRetryDuration/time.Second)
	var failedGrantLeaseAttempts = 0
	for {
		if err := p.Start(ctx, p.options); err != nil {
			failedGrantLeaseAttempts++
			if *p.options.grantLeaseMaxRetries <= failedGrantLeaseAttempts {
				return errors.WithMessagef(err, "%v exceeded max attempts to renew etcd lease %v %v",
					xrutil.GetCodeLocation(1), *p.options.grantLeaseMaxRetries, failedGrantLeaseAttempts)
			}
			xrlog.PrintErr("error granting etcd lease, will retry.", err)
			time.Sleep(grantLeaseRetryDuration)
			continue
		} else {
			// 续租
			if err = p.KeepAlive(ctx); err != nil {
				failedGrantLeaseAttempts++
				if failedGrantLeaseAttempts >= *p.options.grantLeaseMaxRetries {
					return errors.WithMessagef(err, "%v exceeded max attempts to renew etcd lease %v %v",
						xrutil.GetCodeLocation(1), *p.options.grantLeaseMaxRetries, failedGrantLeaseAttempts)
				}
				xrlog.PrintErr("error granting etcd lease, will retry.", err)
				time.Sleep(grantLeaseRetryDuration)
				continue
			} else {
				return nil
			}
		}
	}
}

// KeepAlive 更新租约
func (p *mgr) KeepAlive(ctx context.Context) error {
	var err error
	p.leaseKeepAliveResponseChannel, err = p.lease.KeepAlive(ctx, p.leaseGrantResponse.ID)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}

	p.waitGroup.Add(1)

	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc

	go func(ctx context.Context) {
		defer func() {
			if xrutil.IsRelease() {
				if err := recover(); err != nil {
					xrlog.PrintErr(xrconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			p.waitGroup.Done()
			xrlog.PrintInfo(xrconstant.GoroutineDone)
		}()

		for {
			select {
			case <-ctx.Done():
				xrlog.PrintInfo(xrconstant.GoroutineDone)
				return
			case leaseKeepAliveResponse, ok := <-p.getLeaseKeepAliveResponseChannel():
				xrlog.PrintInfo(leaseKeepAliveResponse, ok)
				if leaseKeepAliveResponse != nil {
					continue
				}
				if ok {
					continue
				}
				// abnormal
				xrlog.PrintErr("etcd lease KeepAlive died, retrying")
				go func(ctx context.Context) {
					defer func() {
						if xrutil.IsRelease() {
							if err := recover(); err != nil {
								xrlog.PrintErr(xrconstant.Retry, xrconstant.GoroutinePanic, err, debug.Stack())
							}
						}
						xrlog.PrintInfo(xrconstant.Retry, xrconstant.GoroutineDone)
					}()
					if err := p.Stop(); err != nil {
						xrlog.PrintInfo(xrconstant.Retry, xrconstant.Failure, err)
						return
					}
					if err := p.retryKeepAlive(ctx); err != nil {
						xrlog.PrintErr(xrconstant.Retry, xrconstant.Failure, err)
						return
					}
				}(context.TODO())
				return
			}
		}
	}(ctxWithCancel)
	return nil
}

// PutWithLease 将一个键值对放入etcd中 WithLease 带ttl
func (p *mgr) PutWithLease(key string, value string) (*clientv3.PutResponse, error) {
	putResponse, err := p.kv.Put(context.TODO(), key, value, clientv3.WithLease(p.leaseGrantResponse.ID))
	if err != nil {
		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	return putResponse, nil
}

// Put 将一个键值对放入etcd中
func (p *mgr) Put(key string, value string) (*clientv3.PutResponse, error) {
	putResponse, err := p.kv.Put(context.TODO(), key, value)
	if err != nil {
		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	return putResponse, nil
}

// DelWithPrefix 删除键值 匹配的键值
func (p *mgr) DelWithPrefix(keyPrefix string) (*clientv3.DeleteResponse, error) {
	deleteResponse, err := p.kv.Delete(context.TODO(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	return deleteResponse, nil
}

// Del 删除键值
func (p *mgr) Del(key string) (*clientv3.DeleteResponse, error) {
	deleteResponse, err := p.kv.Delete(context.TODO(), key)
	if err != nil {
		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	return deleteResponse, nil
}

// DelRange 按选项删除范围内的键值
func (p *mgr) DelRange(startKeyPrefix string, endKeyPrefix string) (*clientv3.DeleteResponse, error) {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithFromKey(),
		clientv3.WithRange(endKeyPrefix),
	}
	deleteResponse, err := p.kv.Delete(context.TODO(), startKeyPrefix, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	return deleteResponse, nil
}

// WatchPrefix 监视以key为前缀的所有 key value
func (p *mgr) WatchPrefix(key string) clientv3.WatchChan {
	return p.client.Watch(context.TODO(), key, clientv3.WithPrefix())
}

// Get 检索键
func (p *mgr) Get(key string) (*clientv3.GetResponse, error) {
	getResponse, err := p.kv.Get(context.TODO(), key)
	if err != nil {
		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	return getResponse, nil
}

// GetPrefix 查找以key为前缀的所有 key value
func (p *mgr) GetPrefix(key string) (*clientv3.GetResponse, error) {
	getResponse, err := p.kv.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	return getResponse, nil
}

// GetPrefixIntoChan  取得关心的前缀,放入 chan 中
func (p *mgr) GetPrefixIntoChan(preFix string) (err error) {
	getResponse, err := p.GetPrefix(preFix)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	for _, v := range getResponse.Kvs {
		p.options.eventChan <- &KV{
			Key:   string(v.Key),
			Value: string(v.Value),
		}
	}
	return
}

// WatchPrefixIntoChan 监听key变化,放入 chan 中
func (p *mgr) WatchPrefixIntoChan(preFix string) (err error) {
	eventChan := p.WatchPrefix(preFix)
	go func() {
		defer func() {
			if xrutil.IsRelease() {
				if err := recover(); err != nil {
					xrlog.PrintErr(xrconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			xrlog.PrintInfo(xrconstant.GoroutineDone)
		}()
		for v := range eventChan {
			Key := string(v.Events[0].Kv.Key)
			Value := string(v.Events[0].Kv.Value)
			p.options.eventChan <- &KV{
				Key:   Key,
				Value: Value,
			}
		}
	}()
	return
}
