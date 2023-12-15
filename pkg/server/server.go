package server

import (
	"context"
)

type IServer interface {
	OnLoadBench(ctx context.Context, opts ...*Options) error //加载bench.json配置文件
	OnInit(ctx context.Context, opts ...*Options) error      //初始化服务资源
	OnStart(ctx context.Context) error                       //启动
	OnRun(ctx context.Context) error                         //运行
	OnPreStop(ctx context.Context) error                     //停止前的处理
	OnStop(ctx context.Context) error                        //停止
}
