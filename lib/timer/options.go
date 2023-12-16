package timer

import (
	"github.com/pkg/errors"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	"time"
)

// NewOptions 新的Options
func NewOptions() *Options {
	ops := new(Options)
	return ops
}

// Options contains Options to configure instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	scanSecondDuration      *time.Duration     // 扫描秒级定时器,纳秒间隔(如 100000000,则每100毫秒扫描一次秒定时器)
	scanMillisecondDuration *time.Duration     // 扫描毫秒级定时器,纳秒间隔(如 100000000,则每100毫秒扫描一次毫秒定时器)
	outgoingTimeoutChan     chan<- interface{} // 是超时事件放置的channel,由外部传入.超时的*Second/*Millisecond都会放入其中
}

func (p *Options) WithScanSecondDuration(scanSecondDuration *time.Duration) *Options {
	p.scanSecondDuration = scanSecondDuration
	return p
}

func (p *Options) WithScanMillisecondDuration(scanMillisecondDuration *time.Duration) *Options {
	p.scanMillisecondDuration = scanMillisecondDuration
	return p
}

func (p *Options) WithOutgoingTimerOutChan(timeoutChan chan<- interface{}) *Options {
	p.outgoingTimeoutChan = timeoutChan
	return p
}

// merge combines the given *Options into a single *Options in a last one wins fashion.
// The specified Options are merged with the existing Options on the server, with the specified Options taking
// precedence.
func merge(opts ...*Options) *Options {
	so := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.scanSecondDuration != nil {
			so.scanSecondDuration = opt.scanSecondDuration
		}
		if opt.scanMillisecondDuration != nil {
			so.scanMillisecondDuration = opt.scanMillisecondDuration
		}
		if opt.outgoingTimeoutChan != nil {
			so.outgoingTimeoutChan = opt.outgoingTimeoutChan
		}
	}
	return so
}

// 配置
func (p *Mgr) configure(opts *Options) error {
	if opts.outgoingTimeoutChan == nil {
		return errors.WithMessage(liberror.Param, libruntime.GetCodeLocation(1).String())
	}
	if opts.scanSecondDuration == nil && opts.scanMillisecondDuration == nil { // 秒 && 毫秒 都未启用
		return errors.WithMessage(liberror.Param, libruntime.GetCodeLocation(1).String())
	}
	return nil
}
