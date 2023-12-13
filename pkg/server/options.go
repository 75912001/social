package server

import (
	"github.com/pkg/errors"
	"path"
	"path/filepath"
	pkgbench "social/pkg/bench"
	libetcd "social/pkg/lib/etcd"
	libutil "social/pkg/lib/util"
)

type OnDefaultHandler func(v interface{}) error

// options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type options struct {
	path                   *string // 路径
	benchPath              *string // 配置文件路径
	subBench               pkgbench.ISubBench
	defaultHandler         OnDefaultHandler // default 处理函数
	etcdHandler            libetcd.OnFunc   // etcd 处理函数
	etcdWatchServicePrefix *string          // etcd 关注 服务 前缀
	etcdWatchCommandPrefix *string          // etcd 关注 命令 前缀
}

// NewOptions 新的Options
func NewOptions() *options {
	ops := new(options)
	return ops
}

func (p *options) SetPath(path string) *options {
	p.path = &path
	return p
}

func (p *options) SetBenchPath(benchPath string) *options {
	p.benchPath = &benchPath
	return p
}

func (p *options) SetSubBench(subBench pkgbench.ISubBench) *options {
	p.subBench = subBench
	return p
}

func (p *options) SetDefaultHandler(defaultHandler OnDefaultHandler) *options {
	p.defaultHandler = defaultHandler
	return p
}

func (p *options) SetEtcdHandler(etcdHandler libetcd.OnFunc) *options {
	p.etcdHandler = etcdHandler
	return p
}

func (p *options) SetEtcdWatchServicePrefix(etcdWatchServicePrefix string) *options {
	p.etcdWatchServicePrefix = &etcdWatchServicePrefix
	return p
}

func (p *options) SetEtcdWatchCommandPrefix(etcdWatchCommandPrefix string) *options {
	p.etcdWatchCommandPrefix = &etcdWatchCommandPrefix
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
		if opt.path != nil {
			so.path = opt.path
		}
		if opt.benchPath != nil {
			so.benchPath = opt.benchPath
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
		if opt.etcdWatchServicePrefix != nil {
			so.etcdWatchServicePrefix = opt.etcdWatchServicePrefix
		}
		if opt.etcdWatchCommandPrefix != nil {
			so.etcdWatchCommandPrefix = opt.etcdWatchCommandPrefix
		}
	}
	return so
}

// 配置
func configure(opts *options) error {
	if opts.path == nil { // 当前目录
		pathValue, err := libutil.GetCurrentPath()
		if err != nil {
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
		opts.path = &pathValue
	}
	if opts.benchPath == nil {
		benchPath := path.Join(*opts.path, "bench.json")
		benchPath = filepath.ToSlash(benchPath)
		opts.benchPath = &benchPath
	}
	return nil
}
