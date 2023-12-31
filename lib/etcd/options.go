package etcd

import (
	"github.com/pkg/errors"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	"time"
)

// OnFunc 处理数据
type OnFunc func(key string, value string) error

// KV key-value pair
type KV struct {
	Key   string
	Value string
}

// NewOptions 新的Options
func NewOptions() *Options {
	return new(Options)
}

// Options contains Options to configure instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	addrs                []string           // 地址
	ttl                  *int64             // Time To Live, etcd内部会按照 ttl/3 的时间(最小1秒),保持连接
	grantLeaseMaxRetries *int               // 授权租约 最大 重试次数 默认:600
	kvSlice              []KV               // Put puts a key-value pair into etcd.
	dialTimeout          *time.Duration     // dialTimeout is the timeout for failing to establish a connection.
	onFunc               OnFunc             // 回调 处理数据
	outgoingEventChan    chan<- interface{} // 传出 channel
	watchServicePrefix   *string            // 关注 服务 前缀
	watchCommandPrefix   *string            // 关注 命令 前缀
}

func (p *Options) WithWatchServicePrefix(watchServicePrefix string) *Options {
	p.watchServicePrefix = &watchServicePrefix
	return p
}

func (p *Options) WithWatchCommandPrefix(watchCommandPrefix string) *Options {
	p.watchCommandPrefix = &watchCommandPrefix
	return p
}

func (p *Options) WithAddrs(addrs []string) *Options {
	p.addrs = p.addrs[:0]
	p.addrs = append(p.addrs, addrs...)
	return p
}

func (p *Options) WithTTL(ttl int64) *Options {
	p.ttl = &ttl
	return p
}

func (p *Options) WithGrantLeaseMaxRetries(retries int) *Options {
	p.grantLeaseMaxRetries = &retries
	return p
}

func (p *Options) WithKV(kv []KV) *Options {
	p.kvSlice = p.kvSlice[:0]
	p.kvSlice = append(p.kvSlice, kv...)
	return p
}

func (p *Options) WithDialTimeout(dialTimeout time.Duration) *Options {
	p.dialTimeout = &dialTimeout
	return p
}

func (p *Options) WithOnFunc(onFunc OnFunc) *Options {
	p.onFunc = onFunc
	return p
}

func (p *Options) WithOutgoingEventChan(eventChan chan<- interface{}) *Options {
	p.outgoingEventChan = eventChan
	return p
}

// merge combines the given *Options into a single *Options in a last one wins fashion.
// The specified Options are merged with the existing Options, with the specified Options taking
// precedence.
func merge(opts ...*Options) *Options {
	newOptions := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if len(opt.addrs) != 0 {
			newOptions.WithAddrs(opt.addrs)
		}
		if opt.ttl != nil {
			newOptions.WithTTL(*opt.ttl)
		}
		if opt.grantLeaseMaxRetries != nil {
			newOptions.WithGrantLeaseMaxRetries(*opt.grantLeaseMaxRetries)
		}
		if len(opt.kvSlice) != 0 {
			newOptions.WithKV(opt.kvSlice)
		}
		if opt.dialTimeout != nil {
			newOptions.WithDialTimeout(*opt.dialTimeout)
		}
		if opt.onFunc != nil {
			newOptions.WithOnFunc(opt.onFunc)
		}
		if opt.outgoingEventChan != nil {
			newOptions.WithOutgoingEventChan(opt.outgoingEventChan)
		}
		if opt.watchServicePrefix != nil {
			newOptions.watchServicePrefix = opt.watchServicePrefix
		}
		if opt.watchCommandPrefix != nil {
			newOptions.watchCommandPrefix = opt.watchCommandPrefix
		}
	}
	return newOptions
}

// 配置
func configure(opt *Options) error {
	if len(opt.addrs) == 0 {
		return errors.WithMessage(liberror.Param, libruntime.GetCodeLocation(1).String())
	}
	if opt.ttl == nil {
		ttl := TtlSecondDefault
		opt.ttl = &ttl
	}
	if opt.grantLeaseMaxRetries == nil {
		opt.grantLeaseMaxRetries = &grantLeaseMaxRetriesDefault
	}
	if len(opt.kvSlice) == 0 {
		return errors.WithMessage(liberror.Param, libruntime.GetCodeLocation(1).String())
	}
	if opt.dialTimeout == nil {
		opt.dialTimeout = &dialTimeoutDefault
	}
	if opt.onFunc == nil {
		return errors.WithMessage(liberror.Param, libruntime.GetCodeLocation(1).String())
	}
	if opt.outgoingEventChan == nil {
		return errors.WithMessage(liberror.Param, libruntime.GetCodeLocation(1).String())
	}
	return nil
}
