package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"runtime"
	"runtime/debug"
	apigate "social/api/gate"
	gatehandler "social/internal/gate/bus"
	libconstant "social/lib/consts"
	libutil "social/lib/util"
	protogate "social/pkg/proto/gate"
	pkgserver "social/pkg/server"
)

var (
	instance *Server
)

// GetInstance 获取
func GetInstance() *Server {
	return instance
}

func NewServer(normal *pkgserver.Normal) *Server {
	instance = &Server{
		Normal: normal,
	}
	normal.Options.WithDefaultHandler(gatehandler.OnEventDefault).WithEtcdHandler(gatehandler.OnEventEtcd)
	return instance
}

type Server struct {
	*pkgserver.Normal
	timer *gatehandler.Timer
}

func (p *Server) OnStart(ctx context.Context) (err error) {
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//
	//向客户端发送通知消息
	//go func() {
	//	time.Sleep(time.Second * 60)
	//	// 模拟异步发送通知给特定客户端
	//	notification := &social_service.CommonResponse{
	//		Response: &social_service.CommonResponse_LogoutResponse{
	//			LogoutResponse: &social_service.LogoutResponse{
	//				Field2: 1,
	//			},
	//		},
	//	}
	//	clients.Range(func(key, value interface{}) bool {
	//		fmt.Printf("Key: %v, Value: %v\n", key, value)
	//		if err := value.(social_service.SocialService_BidirectionalStreamingMethodServer).Send(notification); err != nil {
	//			log.Printf("Error sending no")
	//		}
	//		return true // 返回 true 继续遍历，返回 false 停止遍历
	//
	//
	//		LogoutResponse: &pb.ResponseTypeA{
	//			Field1: "Notification to client",
	//		},
	//	}
	//	clientStream, ok := s.clients.Load(clientID)
	//	if ok {
	//		// 发送通知消息
	//		if err := clientStream.(pb.MyService_MyBidirectionalStreamingMethodServer).Send(notification); err != nil {
	//			log.Printf("Error sending notification: %v", err)
	//		}
	//	}
	//}()

	// 服定时器
	p.timer = gatehandler.NewTimer()
	p.timer.Start()

	go func() { //启动grpc服务
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					p.LogMgr.Fatalf(libconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			// todo menglingchao 关闭grpc服务...
			// p.waitGroupOutPut.Done()
			p.LogMgr.Fatalf(libconstant.GoroutineDone)
		}()
		addr := fmt.Sprintf("%v:%v", p.BenchMgr.Server.IP, p.BenchMgr.Server.Port)
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			p.LogMgr.Fatalf("Failed to listen: %v", err)
		}

		newServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 1024)) //todo menglingchao 设置接受大小
		protogate.RegisterServiceServer(newServer, &apigate.Server{})

		p.LogMgr.Tracef("Server is running on %v", addr)
		if err = newServer.Serve(listen); err != nil {
			p.LogMgr.Fatalf("Failed to serve: %v", err)
		}
	}()
	runtime.GC()
	return nil
}

func (p *Server) OnPreStop(ctx context.Context) (err error) {
	p.timer.Stop()
	p.LogMgr.Warn("serverTimer stop")
	{ // todo menglingchao 关机前处理...
		// todo menglingchao 关闭grpc服务 拒绝新连接
		p.LogMgr.Warn("grpc Service stop")
	}
	return nil
}
