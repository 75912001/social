package log

import (
	"github.com/pkg/errors"
	"os"
	liberror "social/lib/error"
	libutil "social/lib/util"
)

// options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type options struct {
	level          *Level     // 日志等级
	absPath        *string    // 日志绝对路径
	isReportCaller *bool      // 是否打印调用信息 default: true
	namePrefix     *string    // 日志名 前缀
	hooks          LevelHooks // 各日志级别对应的钩子
	isWriteFile    *bool      // 是否写文件 default: true
	enablePool     *bool      // 使用内存池 default: true
}

// NewOptions 新的Options
func NewOptions() *options {
	ops := new(options)
	ops.hooks = make(LevelHooks)
	return ops
}

func (p *options) WithLevel(level Level) *options {
	p.level = &level
	return p
}

func (p *options) WithAbsPath(absPath string) *options {
	p.absPath = &absPath
	return p
}

func (p *options) WithIsReportCaller(isReportCaller bool) *options {
	p.isReportCaller = &isReportCaller
	return p
}

func (p *options) WithNamePrefix(namePrefix string) *options {
	p.namePrefix = &namePrefix
	return p
}

func (p *options) WithHooks(hooks LevelHooks) *options {
	p.hooks = hooks
	return p
}

func (p *options) WithIsWriteFile(isWriteFile bool) *options {
	p.isWriteFile = &isWriteFile
	return p
}

func (p *options) WithEnablePool(enablePool bool) *options {
	p.enablePool = &enablePool
	return p
}

func (p *options) IsEnablePool() bool {
	return *p.enablePool
}

// AddHooks 添加钩子
func (p *options) AddHooks(hook Hook) *options {
	p.hooks.add(hook)
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
		if opt.level != nil {
			so.level = opt.level
		}
		if opt.absPath != nil {
			so.absPath = opt.absPath
		}
		if opt.isReportCaller != nil {
			so.isReportCaller = opt.isReportCaller
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
func configure(opts *options) error {
	if opts.level == nil {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	if opts.absPath == nil {
		return errors.WithMessage(liberror.Param, libutil.GetCodeLocation(1).String())
	}
	if opts.isReportCaller == nil {
		var reportCaller = true
		opts.isReportCaller = &reportCaller
	}
	if opts.isWriteFile == nil {
		var writeFile = true
		opts.isWriteFile = &writeFile
	}
	if opts.enablePool == nil {
		var enablePool = true
		opts.enablePool = &enablePool
	}
	if err := os.MkdirAll(*opts.absPath, os.ModePerm); err != nil {
		return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	return nil
}
