package server

import (
	"context"
)

type IServer interface {
	OnLoadBench(ctx context.Context, opts ...*Options) error //加载bench.json配置文件
	OnInit(ctx context.Context, opts ...*Options) error      //初始化服务资源
	OnStart(ctx context.Context) error                       //启动服务
	OnRun(ctx context.Context) error                         //运行服务
	OnPreStop(ctx context.Context) error                     //处理服务停止前的逻辑
	OnStop(ctx context.Context) error                        //停止服务
}
