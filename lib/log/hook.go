package log

import (
	"github.com/pkg/errors"
	libutil "social/lib/util"
)

// Hook 钩子
type Hook interface {
	Levels() []Level         //需要hook的等级列表
	Fire(entry *entry) error //执行的方法
}

type LevelHooks map[Level][]Hook

// add 添加钩子
func (hooks LevelHooks) add(hook Hook) {
	for _, level := range hook.Levels() {
		hooks[level] = append(hooks[level], hook)
	}
}

// fire 处理钩子
func (hooks LevelHooks) fire(level Level, entry *entry) error {
	for _, hook := range hooks[level] {
		if err := hook.Fire(entry); err != nil {
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
	}
	return nil
}
