package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	api_gate "social/api/gate"
	"social/pkg/proto/gate"
)

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
