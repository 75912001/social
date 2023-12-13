package server

import (
	"context"
)

type IServer interface {
	LoadBench(ctx context.Context, opts ...*options) error
	Init(ctx context.Context, opts ...*options) error
	Start(ctx context.Context) error
	Run(ctx context.Context) error
	PreStop(ctx context.Context) error //处理服务停止前的逻辑
	Stop(ctx context.Context) error
}
