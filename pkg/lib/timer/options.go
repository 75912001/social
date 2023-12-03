package timer

import (
	"github.com/pkg/errors"
	liberror "social/pkg/lib/error"
	libutil "social/pkg/lib/util"
	"time"
)

// NewOptions 新的Options
func NewOptions() *options {
	ops := new(options)
	return ops
}

// options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type options struct {
	scanSecondDuration      *time.Duration     // 扫描秒级定时器,纳秒间隔(如 100000000,则每100毫秒扫描一次秒定时器)
	scanMillisecondDuration *time.Duration     // 扫描毫秒级定时器,纳秒间隔(如 100000000,则每100毫秒扫描一次毫秒定时器)
	outgoingTimeoutChan     chan<- interface{} // 是超时事件放置的channel,由外部传入.超时的*Second/*Millisecond都会放入其中
}

func (p *options) SetScanSecondDuration(scanSecondDuration *time.Duration) *options {
	p.scanSecondDuration = scanSecondDuration
	return p
}

func (p *options) SetScanMillisecondDuration(scanMillisecondDuration *time.Duration) *options {
	p.scanMillisecondDuration = scanMillisecondDuration
	return p
}

func (p *options) SetOutgoingTimerOutChan(timeoutChan chan<- interface{}) *options {
	p.outgoingTimeoutChan = timeoutChan
	return p
}

// mergeOptions combines the given *options into a single *options in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func mergeOptions(opts ...*options) *options {
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
func (p *Mgr) configure(opts *options) error {
	if opts.outgoingTimeoutChan == nil {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	if opts.scanSecondDuration == nil && opts.scanMillisecondDuration == nil { // 秒 && 毫秒 都未启用
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	return nil
}
