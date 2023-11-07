package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"social/impl/common"
	"social/proto/social_service"
)

type socialServer struct {
	social_service.UnimplementedSocialServiceServer
}

func init() {
}

func (s *socialServer) BidirectionalStreamingMethod(stream social_service.SocialService_BidirectionalStreamingMethodServer) error {
	for {
		request, err := stream.Recv()
		fmt.Println(&stream)
		if err != nil {
			e := common.ClientStreamMgrInstance.Del(stream)
			if e != nil {
				//TODO
			}
			return err
		}

		// 根据请求类型选择处理逻辑
		switch req := request.GetRequest().(type) {
		case *social_service.CommonRequest_RegisterRequest:
			fmt.Printf("Received RegisterRequest: %s\n", req.RegisterRequest.GetServiceKey())
			// 处理 RegisterRequest 并生成响应
			response := &social_service.CommonResponse{
				Response: &social_service.CommonResponse_RegisterResponse{
					RegisterResponse: &social_service.RegisterResponse{
						Field1: "Response to " + req.RegisterRequest.GetServiceKey().GetServiceName(),
					},
				},
			}
			e := common.ClientStreamMgrInstance.Add(req.RegisterRequest.GetServiceKey().GetServiceID(), stream)
			if e != nil {
				//TODO
			}
			if err := stream.Send(response); err != nil {
				return err
			}
		case *social_service.CommonRequest_LogoutRequest:
			fmt.Printf("Received LogoutRequest: %d\n", req.LogoutRequest.GetServiceKey())
			// 处理 RequestTypeB 并生成响应
			response := &social_service.CommonResponse{
				Response: &social_service.CommonResponse_LogoutResponse{
					LogoutResponse: &social_service.LogoutResponse{
						Field2: req.LogoutRequest.GetServiceKey().GetServiceID(),
					},
				},
			}
			//删除client
			e := common.ClientStreamMgrInstance.Del(stream)
			if e != nil {
				//TODO
			}
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

func main() {
	// 向客户端发送通知消息
	go func() {
		//time.Sleep(time.Second * 60)
		//// 模拟异步发送通知给特定客户端
		//notification := &social_service.CommonResponse{
		//	Response: &social_service.CommonResponse_LogoutResponse{
		//		LogoutResponse: &social_service.LogoutResponse{
		//			Field2: 1,
		//		},
		//	},
		//}
		//clients.Range(func(key, value interface{}) bool {
		//	fmt.Printf("Key: %v, Value: %v\n", key, value)
		//	if err := value.(social_service.SocialService_BidirectionalStreamingMethodServer).Send(notification); err != nil {
		//		log.Printf("Error sending no")
		//	}
		//	return true // 返回 true 继续遍历，返回 false 停止遍历
		//})
	}()
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

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	social_service.RegisterSocialServiceServer(server, &socialServer{})

	fmt.Println("Server is running on :50051")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
