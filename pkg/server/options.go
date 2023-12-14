package server

import (
	"github.com/pkg/errors"
	libetcd "social/lib/etcd"
	pkgsubbench "social/pkg/subbench"
)

type OnDefaultHandler func(v interface{}) error

// Options contains Options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	subBench       pkgsubbench.ISubBench
	defaultHandler OnDefaultHandler // default 处理函数
	etcdHandler    libetcd.OnFunc   // etcd 处理函数
}

// NewOptions 新的Options
func NewOptions() *Options {
	return new(Options)
}

func (p *Options) WithSubBench(subBench pkgsubbench.ISubBench) *Options {
	p.subBench = subBench
	return p
}

func (p *Options) WithDefaultHandler(defaultHandler OnDefaultHandler) *Options {
	p.defaultHandler = defaultHandler
	return p
}

func (p *Options) WithEtcdHandler(etcdHandler libetcd.OnFunc) *Options {
	p.etcdHandler = etcdHandler
	return p
}

// mergeOptions combines the given *Options into a single *Options in a last one wins fashion.
// The specified Options are merged with the existing Options on the server, with the specified Options taking
// precedence.
func mergeOptions(opts ...*Options) *Options {
	so := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.subBench != nil {
			so.subBench = opt.subBench
		}
		if opt.defaultHandler != nil {
			so.defaultHandler = opt.defaultHandler
		}
		if opt.etcdHandler != nil {
			so.etcdHandler = opt.etcdHandler
		}
	}
	return so
}

// 配置
func configure(opt *Options) error {
	if opt.etcdHandler == nil {
		return errors.WithMessage(errors.New("etcdHandler is nil"), "etcdHandler is nil")
	}
	return nil
}
