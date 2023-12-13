package robot

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"runtime"
	"social/internal/robot/handler"
	"social/internal/robot/subbench"
	"social/pkg"
	xrlog "social/pkg/lib/log"
	xrpb "social/pkg/lib/pb"
	xrutil "social/pkg/lib/util"
	"social/pkg/msg"
	"social/pkg/proto"
	"social/pkg/proto/gate"
	"social/pkg/server"
	"syscall"
)

type Server struct{}

func (p *Server) Start(ctx context.Context) (err error) {
	err = server.GetInstance().PreInit(context.TODO(),
		server.NewOptions().
			SetDefaultHandler(handler.OnEventDefault).SetSubBench(subbench.GetInstance()))
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	runtime.GC()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 服定时器
	serverTimer := new(handler.ServerTimer)
	serverTimer.Start()
	defer func() {
		serverTimer.Stop()
		xrlog.GetInstance().Warn("serverTimer stop")
	}()
	// 连接 gRPC 服务器
	addr := fmt.Sprintf("%v:%v", subbench.GetInstance().Gate.IP, subbench.GetInstance().Gate.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		xrlog.GetInstance().Fatalf("Failed to connect: %v %v", addr, err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			xrlog.GetInstance().Fatalf("err:%s", err)
		}
	}()

	client := gate.NewServiceClient(conn)
	// 创建双向流
	stream, err := client.BidirectionalBinaryData(ctx)
	if err != nil {
		xrlog.GetInstance().Fatalf("Failed to create stream: %v", err)
	}

	// 发送多个请求
	requests := []*gate.RegisterReq{
		{
			ServiceKey: &proto.ServiceKey{
				ZoneID:      pkg.GZoneID,
				ServiceName: pkg.GServiceName,
				ServiceID:   pkg.GServiceID,
			},
		},
		{
			ServiceKey: &proto.ServiceKey{
				ZoneID:      pkg.GZoneID,
				ServiceName: pkg.GServiceName,
				ServiceID:   pkg.GServiceID,
			},
		},
	}

	for _, req := range requests {
		data, err := msg.Marshal(
			&xrpb.UnserializedPacket{
				Header: &msg.Header{
					MessageID: gate.RegisterReq_CMD,
					ResultID:  0,
				},
				Message: req,
			})
		if err != nil {
			xrlog.GetInstance().Fatalf("Failed to create stream: %v", err)
			return err
		}
		if err := stream.Send(&proto.BinaryData{Data: data}); err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
	}

	// 接收多个响应
	//for {
	//	response, err := stream.Recv()
	//	if err != nil {
	//		log.Fatalf("Failed to receive response: %v", err)
	//	}
	//
	//	// 根据响应类型处理响应
	//	switch resp := response.GetResponse().(type) {
	//	case *social_service.CommonResponse_RegisterResponse:
	//		fmt.Printf("Received RegisterResponse: %s\n", resp.RegisterResponse.GetField1())
	//	case *social_service.CommonResponse_LogoutResponse:
	//		fmt.Printf("Received RegisterResponse: %d\n", resp.LogoutResponse.GetField2())
	//	}
	//}
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
	return nil
}

func (p *Server) Stop(ctx context.Context) (err error) {
	{ // todo menglingchao 关机前处理...
		// todo menglingchao 关闭grpc服务 拒绝新连接
		xrlog.GetInstance().Warn("grpc Service stop")
	}
	// 设置为关闭中
	server.GetInstance().SetStopping()
	_ = server.GetInstance().Stop()

	xrlog.PrintErr("server Log stop")
	_ = xrlog.GetInstance().Stop()
	return nil
}
