package actor

import (
	"github.com/pkg/errors"
)

type OnDefaultHandler func(msg *Msg) error

// Options contains Options to configure instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	onHandler        OnDefaultHandler // default 处理函数
	actorChannelSize *uint32          // actor channel 大小
}

// NewOptions 新的Options
func NewOptions() *Options {
	return new(Options)
}

func (p *Options) WithDefaultHandler(defaultHandler OnDefaultHandler) *Options {
	p.onHandler = defaultHandler
	return p
}

func (p *Options) WithActorChannelSize(actorChannelSize uint32) *Options {
	p.actorChannelSize = &actorChannelSize
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
		if opt.onHandler != nil {
			so.onHandler = opt.onHandler
		}
		if opt.actorChannelSize != nil {
			so.actorChannelSize = opt.actorChannelSize
		}
	}
	return so
}

// 配置
func configure(opt *Options) error {
	if opt.onHandler == nil {
		return errors.WithMessage(errors.New("onHandler is nil"), "onHandler is nil")
	}
	if opt.actorChannelSize == nil {
		return errors.WithMessage(errors.New("actorChannelSize is nil"), "actorChannelSize is nil")
	}
	return nil
}
