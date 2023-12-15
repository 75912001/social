package actor

import (
	"context"
	liberror "social/lib/error"
)

type Normal struct {
	//https://zhuanlan.zhihu.com/p/427806717
	//状态 state
	//指actor本身的属性信息，state只能被actor自己操作，不能被其他actor共享和操作，有效的避免加锁和数据竞争

	//行为 behavior
	//指actor处理逻辑，如果通过行为来操作自身state

	//Mailbox邮箱

}

func (a *Normal) Send(msg IMsg) error {
	return liberror.NotImplemented
}

func (a *Normal) OnStart(ctx context.Context) error {
	//todo menglingchao 创建一个协程,处理逻辑
	return liberror.NotImplemented
}

func (a *Normal) OnRun(ctx context.Context) error {
	return liberror.NotImplemented
}

func (a *Normal) OnPreStop(ctx context.Context) error {
	//todo menglingchao 停止前的处理
	return liberror.NotImplemented
}

func (a *Normal) OnStop(ctx context.Context) error {
	//todo menglingchao 停止处理
	return liberror.NotImplemented
}
