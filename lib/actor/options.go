package actor

import (
	"github.com/pkg/errors"
)

type OnDefaultHandler func(v interface{}) error

// Options contains Options to configure instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	defaultHandler OnDefaultHandler // default 处理函数
}

// NewOptions 新的Options
func NewOptions() *Options {
	return new(Options)
}

func (p *Options) WithDefaultHandler(defaultHandler OnDefaultHandler) *Options {
	p.defaultHandler = defaultHandler
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
		if opt.defaultHandler != nil {
			so.defaultHandler = opt.defaultHandler
		}
	}
	return so
}

// 配置
func configure(opt *Options) error {
	if opt.defaultHandler == nil {
		return errors.WithMessage(errors.New("defaultHandler is nil"), "defaultHandler is nil")
	}
	return nil
}