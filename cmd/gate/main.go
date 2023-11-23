package gate

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	apigate "social/api/gate"
	"social/internal/gate/handler"
	"social/internal/gate/stop"
	"social/pkg"
	"social/pkg/bench"
	"social/pkg/common"
	xrconstant "social/pkg/lib/constant"
	xrlog "social/pkg/lib/log"
	xrutil "social/pkg/lib/util"
	"social/pkg/proto/gate"
	"social/pkg/server"
	"syscall"
)

type Server struct {
}

func (p *Server) Stop() (err error) {
	xrlog.PrintErr("server Log stop")
	_ = xrlog.GetInstance().Stop()
	return nil
}

func (p *Server) Start() (err error) {
	err = server.GetInstance().PreInit(context.TODO(),
		server.NewOptions().
			SetDefaultHandler(handler.OnEventDefault))
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}

	// world服定时器
	serverTimer := new(handler.ServerTimer)
	serverTimer.Start()
	defer func() {
		serverTimer.Stop()
		xrlog.GetInstance().Warn("serverTimer stop")
	}()

	err = server.GetInstance().PostInit(context.TODO(),
		server.NewOptions().
			SetEtcdHandler(handler.OnEventEtcd).
			SetEtcdWatchServicePrefix(fmt.Sprintf("/%v/%v/", common.ProjectName, pkg.EtcdWatchMsgTypeService)).
			SetEtcdWatchCommandPrefix(fmt.Sprintf("/%v/%v/%v/%v/",
				common.ProjectName, pkg.EtcdWatchMsgTypeCommand,
				pkg.GZoneID,
				pkg.GServiceName)),
	)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}

	runtime.GC()

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

	go func() { //启动grpc服务
		defer func() {
			if xrutil.IsRelease() {
				if err := recover(); err != nil {
					xrlog.GetInstance().Fatalf(xrconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			//todo menglingchao 关闭grpc服务...
			// p.waitGroupOutPut.Done()
			xrlog.GetInstance().Fatalf(xrconstant.GoroutineDone)
		}()
		addr := fmt.Sprintf("%v:%v", bench.GetInstance().Server.IP, bench.GetInstance().Server.Port)
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			xrlog.GetInstance().Fatalf("Failed to listen: %v", err)
		}

		server := grpc.NewServer()
		gate.RegisterServiceServer(server, &apigate.Server{})

		xrlog.GetInstance().Tracef("Server is running on %v", addr)
		if err = server.Serve(listen); err != nil {
			xrlog.GetInstance().Fatalf("Failed to serve: %v", err)
		}
	}()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// 退出服务
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case <-server.GetInstance().QuitChan:
		xrlog.GetInstance().Warn("GServer will stop in a few seconds")
	case s := <-sigChan:
		xrlog.GetInstance().Warnf("GServer got signal: %s, shutting down...", s)
	}

	stop.PreStop()
	_ = server.GetInstance().Stop()

	return nil
}
