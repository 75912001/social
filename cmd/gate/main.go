package gate

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	api_gate "social/api/gate"
	"social/internal/gate/handler"
	"social/internal/gate/stop"
	"social/pkg/common"
	"social/pkg/etcd"
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
			SetEtcdWatchServicePrefix(fmt.Sprintf("/%v/%v/", common.ProjectName, etcd.WatchMsgTypeService)).
			SetEtcdWatchCommandPrefix(fmt.Sprintf("/%v/%v/%v/%v/",
				common.ProjectName, etcd.WatchMsgTypeCommand,
				server.GetInstance().ZoneID,
				server.GetInstance().ServiceName)),
	)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}

	runtime.GC()

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

func main() {
	// 向客户端发送通知消息
	//go func() {
	//	//time.Sleep(time.Second * 60)
	//	//// 模拟异步发送通知给特定客户端
	//	//notification := &social_service.CommonResponse{
	//	//	Response: &social_service.CommonResponse_LogoutResponse{
	//	//		LogoutResponse: &social_service.LogoutResponse{
	//	//			Field2: 1,
	//	//		},
	//	//	},
	//	//}
	//	//clients.Range(func(key, value interface{}) bool {
	//	//	fmt.Printf("Key: %v, Value: %v\n", key, value)
	//	//	if err := value.(social_service.SocialService_BidirectionalStreamingMethodServer).Send(notification); err != nil {
	//	//		log.Printf("Error sending no")
	//	//	}
	//	//	return true // 返回 true 继续遍历，返回 false 停止遍历
	//	//})
	//}()
	////		LogoutResponse: &pb.ResponseTypeA{
	////			Field1: "Notification to client",
	////		},
	////	}
	////	clientStream, ok := s.clients.Load(clientID)
	////	if ok {
	////		// 发送通知消息
	////		if err := clientStream.(pb.MyService_MyBidirectionalStreamingMethodServer).Send(notification); err != nil {
	////			log.Printf("Error sending notification: %v", err)
	////		}
	////	}
	////}()

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	gate.RegisterServiceServer(server, &api_gate.Server{})

	fmt.Println("Server is running on :50051")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
