package server

import (
	"github.com/pkg/errors"
	"path"
	"social/pkg/bench"
	xretcd "social/pkg/lib/etcd"
	xrutil "social/pkg/lib/util"
)

// Options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	Path                   *string // 路径
	BenchPath              *string // 配置文件路径
	SubBench               bench.ISubBench
	DefaultHandler         OnDefaultHandler // default 处理函数
	EtcdHandler            xretcd.OnFunc    // etcd 处理函数
	EtcdWatchServicePrefix *string          // etcd 关注 服务 前缀
	EtcdWatchCommandPrefix *string          // etcd 关注 命令 前缀
}

// NewOptions 新的Options
func NewOptions() *Options {
	ops := new(Options)
	return ops
}

func (p *Options) SetSubBench(subBench bench.ISubBench) *Options {
	p.SubBench = subBench
	return p
}

func (p *Options) SetDefaultHandler(defaultHandler OnDefaultHandler) *Options {
	p.DefaultHandler = defaultHandler
	return p
}

func (p *Options) SetEtcdHandler(etcdHandler xretcd.OnFunc) *Options {
	p.EtcdHandler = etcdHandler
	return p
}

func (p *Options) SetEtcdWatchServicePrefix(etcdWatchServicePrefix string) *Options {
	p.EtcdWatchServicePrefix = &etcdWatchServicePrefix
	return p
}

func (p *Options) SetEtcdWatchCommandPrefix(etcdWatchCommandPrefix string) *Options {
	p.EtcdWatchCommandPrefix = &etcdWatchCommandPrefix
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
		if opt.Path != nil {
			so.Path = opt.Path
		}
		if opt.SubBench != nil {
			so.SubBench = opt.SubBench
		}
		if opt.DefaultHandler != nil {
			so.DefaultHandler = opt.DefaultHandler
		}
		if opt.EtcdHandler != nil {
			so.EtcdHandler = opt.EtcdHandler
		}
		if opt.EtcdWatchServicePrefix != nil {
			so.EtcdWatchServicePrefix = opt.EtcdWatchServicePrefix
		}
		if opt.EtcdWatchCommandPrefix != nil {
			so.EtcdWatchCommandPrefix = opt.EtcdWatchCommandPrefix
		}
	}
	return so
}

// 配置
func configure(opts *Options) error {
	if opts.Path == nil { // 当前目录
		pathValue, err := xrutil.GetCurrentPath()
		if err != nil {
			return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
		}
		opts.Path = &pathValue
	}
	if opts.BenchPath == nil {
		benchPath := path.Join(*opts.Path, "bench.json")
		opts.BenchPath = &benchPath
	}
	return nil
}
