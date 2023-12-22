package friend

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"runtime"
	"runtime/debug"
	"social/lib/actor"
	libconsts "social/lib/consts"
	libutil "social/lib/util"
	protofriend "social/pkg/proto/friend"
	pkgserver "social/pkg/server"
)

var (
	app *Friend
)

func NewFriend(normal *pkgserver.Normal) *Friend {
	app = &Friend{
		Normal: normal,
	}
	app.bus.Normal = normal
	normal.Options.
		WithDefaultHandler(app.bus.OnEventBus).
		WithEtcdHandler(app.bus.OnEventEtcd).
		WithTimerEachSecond(&pkgserver.NormalTimerSecond{
			OnTimerFun: app.OnTimerEachSecondFun,
			Arg:        app,
		}).
		WithTimerEachDay(&pkgserver.NormalTimerSecond{
			OnTimerFun: app.OnTimerEachDayFun,
			Arg:        app,
		})
	app.gateMgr.actorMgr = actor.NewMgr[string]()
	return app
}

type Friend struct {
	*pkgserver.Normal
	bus     Bus
	gateMgr UserMgr
}

func (p *Friend) String() string {
	return pkgserver.NameFriend
}

func (p *Friend) OnStart(ctx context.Context) (err error) {
	// 定时器-可用负载
	timerAvailableLoadExpireTimestamp = p.TimeMgr.ShadowTimeSecond()
	... 链接 mongodb
	... 链接 redis
	go func() { //启动grpc服务
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					p.LogMgr.Fatalf(libconsts.GoroutinePanic, err, debug.Stack())
				}
			}
			p.LogMgr.Fatalf(libconsts.GoroutineDone)
		}()
		addr := fmt.Sprintf("%v:%v", p.BenchMgr.Server.IP, p.BenchMgr.Server.Port)
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			p.LogMgr.Fatalf("Failed to listen: %v", err)
		}
		p.GrpcServer = grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 1024)) // todo menglingchao 设置接受大小
		protofriend.RegisterServiceServer(p.GrpcServer, &APIServer{})
		p.LogMgr.Tracef("%v is running on %v", p.String(), addr)
		if err = p.GrpcServer.Serve(listen); err != nil {
			p.LogMgr.Fatalf("Failed to serve: %v", err)
		}
	}()
	runtime.GC()
	return nil
}

func (p *Friend) OnPreStop(_ context.Context) (err error) {
	p.LogMgr.Warn("OnPreStop stop")
	{ // 关机前处理...业务逻辑
		p.LogMgr.Warn("grpc Service stop")
	}
	return nil
}
