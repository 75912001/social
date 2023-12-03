package etcd

import (
	"github.com/pkg/errors"
	liberror "social/pkg/lib/error"
	libutil "social/pkg/lib/util"
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

// NewOptions 新的Options
func NewOptions() *options {
	return new(options)
}

// options contains options to configure instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type options struct {
	addrs                []string           // 地址
	ttl                  *int64             // Time To Live, etcd内部会按照 ttl/3 的时间(最小1秒),保持连接
	grantLeaseMaxRetries *int               // 授权租约 最大 重试次数 默认:600
	kvSlice              []KV               // 事件
	dialTimeout          *time.Duration     // dialTimeout is the timeout for failing to establish a connection.
	onFunc               OnFunc             // 回调 处理数据
	outgoingEventChan    chan<- interface{} // 传出 channel
}

func (p *options) SetAddrs(addrs []string) *options {
	p.addrs = p.addrs[0:0]
	p.addrs = append(p.addrs, addrs...)
	return p
}

func (p *options) SetTTL(ttl int64) *options {
	p.ttl = &ttl
	return p
}

func (p *options) SetGrantLeaseMaxRetries(retries int) *options {
	p.grantLeaseMaxRetries = &retries
	return p
}

func (p *options) SetKV(kv []KV) *options {
	p.kvSlice = p.kvSlice[0:0]
	p.kvSlice = append(p.kvSlice, kv...)
	return p
}

func (p *options) SetDialTimeout(dialTimeout time.Duration) *options {
	p.dialTimeout = &dialTimeout
	return p
}

func (p *options) SetOnFunc(onFunc OnFunc) *options {
	p.onFunc = onFunc
	return p
}

func (p *options) SetOutgoingEventChan(eventChan chan<- interface{}) *options {
	p.outgoingEventChan = eventChan
	return p
}

// mergeOptions combines the given *options into a single *options in a last one wins fashion.
// The specified options are merged with the existing options, with the specified options taking
// precedence.
func mergeOptions(opts ...*options) *options {
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
		if opt.outgoingEventChan != nil {
			no.SetOutgoingEventChan(opt.outgoingEventChan)
		}
	}
	return no
}

// 配置
func configure(opts *options) error {
	if len(opts.addrs) == 0 {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	if opts.ttl == nil {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	if opts.grantLeaseMaxRetries == nil {
		var v = 600
		opts.grantLeaseMaxRetries = &v
	}
	if len(opts.kvSlice) == 0 {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	if opts.dialTimeout == nil {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	if opts.onFunc == nil {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	if opts.outgoingEventChan == nil {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	return nil
}
