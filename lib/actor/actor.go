package actor

import (
	"context"
)

// IActor 是 Actor 接口
type IActor interface {
	OnStart(ctx context.Context, opts ...*Options) error //启动
	OnStop(ctx context.Context) error                    //停止
	SendToMailBox(msg IMsg) error
}
