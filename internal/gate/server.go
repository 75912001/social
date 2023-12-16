package gate

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"runtime"
	"runtime/debug"
	apigate "social/api/gate"
	libconsts "social/lib/consts"
	libutil "social/lib/util"
	protogate "social/pkg/proto/gate"
	pkgserver "social/pkg/server"
)

var (
	server *Server
)

// GetInstance 获取
func GetInstance() *Server {
	return server
}

func NewServer(normal *pkgserver.Normal) *Server {
	server = &Server{
		Normal: normal,
	}
	normal.Options.WithDefaultHandler(server.bus.OnEventBus).WithEtcdHandler(server.bus.OnEventEtcd)
	return server
}

type Server struct {
	*pkgserver.Normal
	bus         Bus
	serverTimer ServerTimer
	router      Router
}

func (p *Server) OnStart(_ context.Context) (err error) {
	// 服定时器
	p.serverTimer.Start()
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
		protogate.RegisterServiceServer(p.GrpcServer, &apigate.Server{})
		p.LogMgr.Tracef("Server is running on %v", addr)
		if err = p.GrpcServer.Serve(listen); err != nil {
			p.LogMgr.Fatalf("Failed to serve: %v", err)
		}
	}()
	runtime.GC()
	return nil
}

func (p *Server) OnPreStop(_ context.Context) (err error) {
	p.serverTimer.Stop()
	p.LogMgr.Warn("serverTimer stop")
	{ // todo menglingchao 关机前处理...
		p.LogMgr.Warn("grpc Service stop")
	}
	return nil
}
