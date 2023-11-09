package log

import (
	"os"
	xrerror "social/pkg/lib/error"
	xrutil "social/pkg/lib/util"

	"github.com/pkg/errors"
)

// Options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	level          *Level     // 日志等级
	absPath        *string    // 日志绝对路径
	isReportCaller *bool      // 是否打印调用信息 default: true
	namePrefix     *string    // 日志名 前缀
	hooks          LevelHooks // 各日志级别对应的钩子
	isWriteFile    *bool      // 是否写文件 default: true
	enablePool     *bool      // 使用内存池 default: true
}

// NewOptions 新的Options
func NewOptions() *Options {
	ops := new(Options)
	ops.hooks = make(LevelHooks)
	return ops
}

func (p *Options) SetLevel(level Level) *Options {
	p.level = &level
	return p
}

func (p *Options) SetIsReportCaller(isReportCaller bool) *Options {
	p.isReportCaller = &isReportCaller
	return p
}

func (p *Options) SetIsWriteFile(isWriteFile bool) *Options {
	p.isWriteFile = &isWriteFile
	return p
}

func (p *Options) SetAbsPath(absPath string) *Options {
	p.absPath = &absPath
	return p
}

func (p *Options) SetNamePrefix(namePrefix string) *Options {
	p.namePrefix = &namePrefix
	return p
}

func (p *Options) SetHooks(hooks LevelHooks) *Options {
	p.hooks = hooks
	return p
}

func (p *Options) SetEnablePool(enablePool bool) *Options {
	p.enablePool = &enablePool
	return p
}

func (p *Options) IsEnablePool() bool {
	return *p.enablePool
}

// AddHooks 添加钩子
func (p *Options) AddHooks(hook Hook) *Options {
	p.hooks.add(hook)
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

		if opt.level != nil {
			so.level = opt.level
		}
		if opt.isReportCaller != nil {
			so.isReportCaller = opt.isReportCaller
		}
		if opt.absPath != nil {
			so.absPath = opt.absPath
		}
		if opt.namePrefix != nil {
			so.namePrefix = opt.namePrefix
		}
		if opt.hooks != nil {
			so.hooks = opt.hooks
		}
		if opt.isWriteFile != nil {
			so.isWriteFile = opt.isWriteFile
		}
		if opt.enablePool != nil {
			so.enablePool = opt.enablePool
		}
	}
	return so
}

// 配置
func configure(opts *Options) error {
	if opts.level == nil {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	if opts.absPath == nil {
		return errors.WithMessage(xrerror.Param, xrutil.GetCodeLocation(1).String())
	}
	if opts.isWriteFile == nil {
		var writeFile = true
		opts.isWriteFile = &writeFile
	}
	if opts.enablePool == nil {
		var enablePool = true
		opts.enablePool = &enablePool
	}
	if opts.isReportCaller == nil {
		var reportCaller = true
		opts.isReportCaller = &reportCaller
	}
	if err := os.MkdirAll(*opts.absPath, os.ModePerm); err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}

	return nil
}
