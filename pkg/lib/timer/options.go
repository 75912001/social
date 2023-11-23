package timer

import (
	xrerror "social/pkg/lib/error"
	xrutil "social/pkg/lib/util"
	"time"

	"github.com/pkg/errors"
)

// Options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	scanSecondDuration      *time.Duration     // 扫描秒级定时器,纳秒间隔(如 100000000,则每100毫秒扫描一次秒定时器)
	scanMillisecondDuration *time.Duration     // 扫描毫秒级定时器,纳秒间隔(如 100000000,则每100毫秒扫描一次毫秒定时器)
	timeoutChan             chan<- interface{} // 是超时事件放置的channel,由外部传入.超时的*Second/*Millisecond都会放入其中
}

// NewOptions 新的Options
func NewOptions() *Options {
	ops := new(Options)
	return ops
}

func (p *Options) SetScanSecondDuration(scanSecondDuration *time.Duration) *Options {
	p.scanSecondDuration = scanSecondDuration
	return p
}

func (p *Options) SetScanMillisecondDuration(scanMillisecondDuration *time.Duration) *Options {
	p.scanMillisecondDuration = scanMillisecondDuration
	return p
}

func (p *Options) SetTimerOutChan(timeoutChan chan<- interface{}) *Options {
	p.timeoutChan = timeoutChan
	return p
}

// mergeOptions combines the given *Options into a single *Options in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func mergeOptions(opts ...*Options) *Options {
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
		if opt.timeoutChan != nil {
			so.timeoutChan = opt.timeoutChan
		}
	}
	return so
}

// 配置
func (p *mgr) configure(opts *Options) error {
	if opts.timeoutChan == nil {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	if opts.scanSecondDuration == nil && opts.scanMillisecondDuration == nil { // 秒 && 毫秒 都未启用
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	return nil
}
