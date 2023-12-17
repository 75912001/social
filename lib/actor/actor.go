package actor

import (
	"context"
)

// IActor 是 Actor 接口
type IActor interface {
	OnStart(ctx context.Context, opts ...*Options) error //启动
	OnPreStop(ctx context.Context) error                 //停止前的处理
	OnStop(ctx context.Context) error                    //停止
	SendToMailBox(msg IMsg) error
}
