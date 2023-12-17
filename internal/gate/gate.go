package gate

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"runtime"
	"runtime/debug"
	libconsts "social/lib/consts"
	libutil "social/lib/util"
	protogate "social/pkg/proto/gate"
	pkgserver "social/pkg/server"
)

var (
	gate *Gate
)

// GetInstance 获取
func GetInstance() *Gate {
	return gate
}

func NewGate(normal *pkgserver.Normal) *Gate {
	gate = &Gate{
		Normal: normal,
	}
	normal.Options.
		WithDefaultHandler(gate.bus.OnEventBus).
		WithEtcdHandler(gate.bus.OnEventEtcd).
		WithTimerEachSecond(&pkgserver.NormalTimerSecond{
			OnTimerFun: gate.OnTimerEachSecondFun,
			Arg:        gate,
		}).
		WithTimerEachDay(&pkgserver.NormalTimerSecond{
			OnTimerFun: gate.OnTimerEachDayFun,
			Arg:        gate,
		})
	return gate
}

type Gate struct {
	*pkgserver.Normal
	bus    Bus
	router Router
}

func (p *Gate) String() string {
	return pkgserver.NameGate
}

func (p *Gate) OnStart(_ context.Context) (err error) {
	// 定时器-可用负载
	timerAvailableLoadExpireTimestamp = p.TimeMgr.ShadowTimeSecond()
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
		p.GrpcServer = grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 1024)) //todo menglingchao 设置接受大小
		protogate.RegisterServiceServer(p.GrpcServer, &APIServer{})
		p.LogMgr.Tracef("Gate is running on %v", addr)
		if err = p.GrpcServer.Serve(listen); err != nil {
			p.LogMgr.Fatalf("Failed to serve: %v", err)
		}
	}()
	runtime.GC()
	return nil
}

func (p *Gate) OnPreStop(_ context.Context) (err error) {
	p.LogMgr.Warn("serverTimer stop")
	{ // todo menglingchao 关机前处理...
		p.LogMgr.Warn("grpc Service stop")
	}
	return nil
}
