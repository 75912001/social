package etcd

import (
	"github.com/pkg/errors"
	xrerror "social/pkg/lib/error"
	xrutil "social/pkg/lib/util"
	"time"
)

var (
	grantLeaseRetryDuration = time.Second * 3 // 授权租约 重试 间隔时长
)

// OnFunc 处理数据
type OnFunc func(key string, value string) error

// KV key-value pair
type KV struct {
	Key   string
	Value string
}

// Options contains options to configure instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	addrs                []string           // 地址
	ttl                  *int64             // Time To Live, etcd内部会按照 ttl/3 的时间(最小1秒),保持连接
	grantLeaseMaxRetries *int               // 授权租约 最大 重试次数 默认:600
	kvSlice              []KV               // 事件
	dialTimeout          *time.Duration     // dialTimeout is the timeout for failing to establish a connection.
	onFunc               OnFunc             // 回调 处理数据
	eventChan            chan<- interface{} // 传出 channel
}

// NewOptions 新的Options
func NewOptions() *Options {
	return new(Options)
}

func (p *Options) SetAddrs(addrs []string) *Options {
	p.addrs = p.addrs[0:0]
	p.addrs = append(p.addrs, addrs...)
	return p
}

func (p *Options) SetTTL(ttl int64) *Options {
	p.ttl = &ttl
	return p
}

func (p *Options) SetGrantLeaseMaxRetries(retries int) *Options {
	p.grantLeaseMaxRetries = &retries
	return p
}

func (p *Options) SetKV(kv []KV) *Options {
	p.kvSlice = p.kvSlice[0:0]
	p.kvSlice = append(p.kvSlice, kv...)
	return p
}

func (p *Options) SetDialTimeout(dialTimeout time.Duration) *Options {
	p.dialTimeout = &dialTimeout
	return p
}

func (p *Options) SetOnFunc(onFunc OnFunc) *Options {
	p.onFunc = onFunc
	return p
}

func (p *Options) SetEventChan(eventChan chan<- interface{}) *Options {
	p.eventChan = eventChan
	return p
}

// mergeOptions combines the given *Options into a single *Options in a last one wins fashion.
// The specified options are merged with the existing options, with the specified options taking
// precedence.
func mergeOptions(opts ...*Options) *Options {
	no := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if len(opt.addrs) != 0 {
			no.SetAddrs(opt.addrs)
		}
		if opt.ttl != nil {
			no.SetTTL(*opt.ttl)
		}
		if opt.grantLeaseMaxRetries != nil {
			no.SetGrantLeaseMaxRetries(*opt.grantLeaseMaxRetries)
		}
		if len(opt.kvSlice) != 0 {
			no.SetKV(opt.kvSlice)
		}
		if opt.dialTimeout != nil {
			no.SetDialTimeout(*opt.dialTimeout)
		}
		if opt.onFunc != nil {
			no.SetOnFunc(opt.onFunc)
		}
		if opt.eventChan != nil {
			no.SetEventChan(opt.eventChan)
		}
	}
	return no
}

// 配置
func configure(opts *Options) error {
	if len(opts.addrs) == 0 {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	if opts.ttl == nil {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	if opts.grantLeaseMaxRetries == nil {
		var v = 600
		opts.grantLeaseMaxRetries = &v
	}
	if len(opts.kvSlice) == 0 {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	if opts.dialTimeout == nil {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	if opts.onFunc == nil {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	if opts.eventChan == nil {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	return nil
}
