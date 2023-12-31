package robot

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	liblog "social/lib/log"
	pkgmsg "social/pkg/msg"
	pkgproto "social/pkg/proto"
	protogate "social/pkg/proto/gate"
	pkgserver "social/pkg/server"
	"time"
)

var (
	robot *Robot
)

func NewRobot(normal *pkgserver.Normal) *Robot {
	robot = &Robot{
		Normal: normal,
	}
	normal.Options.
		WithDefaultHandler(robot.bus.OnEventBus).
		WithEtcdHandler(robot.bus.OnEventEtcd).WithSubBench(&robot.subBenchMgr)
	return robot
}

type Robot struct {
	*pkgserver.Normal
	bus         Bus
	subBenchMgr SubBenchMgr
}

func (p *Robot) OnStart(ctx context.Context) (err error) {
	p.Options.WithDefaultHandler(OnEventDefault)

	// 连接 gRPC 服务器
	addr := fmt.Sprintf("%v:%v", p.subBenchMgr.Gate.IP, p.subBenchMgr.Gate.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		liblog.GetInstance().Fatalf("Failed to connect: %v %v", addr, err)
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			liblog.GetInstance().Fatalf("err:%s", err)
		}
	}()

	client := protogate.NewGateServiceClient(conn)
	// 创建双向流
	stream, err := client.BidirectionalBinaryData(ctx)
	if err != nil {
		liblog.GetInstance().Fatalf("Failed to create stream: %v", err)
		return err
	}
	// 发送请求
	req := &protogate.GateRegisterReq{
		ServiceKey: &pkgproto.ServiceKey{
			ZoneID:      p.ZoneID,
			ServiceName: p.ServiceName,
			ServiceID:   p.ServiceID,
		},
	}

	packet := pkgmsg.Packet{
		Header: pkgmsg.Header{
			MessageID: protogate.GateRegisterReqCMD,
			ResultID:  0,
		},
		Message: req,
	}
	data, err := packet.Marshal()
	if err != nil {
		liblog.GetInstance().Fatalf("Failed to create stream: %v", err)
		return err
	}
	bd := &pkgproto.BinaryData{Data: data}
	//_ = data
	//bd := &pkgproto.BinaryData{}
	if err := stream.Send(bd); err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	time.Sleep(time.Hour)
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
	return nil
}

func (p *Robot) OnPreStop(ctx context.Context) (err error) {
	{ // todo menglingchao 关机前处理...
		// todo menglingchao 关闭grpc服务 拒绝新连接
		liblog.GetInstance().Warn("grpc Service stop")
	}
	return nil
}
