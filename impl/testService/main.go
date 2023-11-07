package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"social/impl/protobuf/proto/social_service"
)

type socialServer struct {
	social_service.UnimplementedSocialServiceServer
}

func (s *socialServer) BidirectionalStreamingMethod(stream social_service.SocialService_BidirectionalStreamingMethodServer) error {
	for {
		request, err := stream.Recv()
		if err != nil {
			return err
		}

		// 根据请求类型选择处理逻辑
		switch req := request.GetRequest().(type) {
		case *social_service.CommonRequest_RegisterRequest:
			fmt.Printf("Received RegisterRequest: %s\n", req.RegisterRequest.GetServiceKey())
			// 处理 RequestTypeA 并生成响应
			response := &social_service.CommonResponse{
				Response: &social_service.CommonResponse_RegisterResponse{RegisterResponse: &social_service.RegisterResponse{Field1: "Response to " + req.RegisterRequest.GetServiceKey().GetServiceName()}},
			}
			if err := stream.Send(response); err != nil {
				return err
			}
		case *social_service.CommonRequest_LogoutRequest:
			fmt.Printf("Received LogoutRequest: %d\n", req.LogoutRequest.GetServiceKey())
			// 处理 RequestTypeB 并生成响应
			response := &social_service.CommonResponse{
				Response: &social_service.CommonResponse_LogoutResponse{LogoutResponse: &social_service.LogoutResponse{Field2: req.LogoutRequest.GetServiceKey().GetServiceID()}},
			}
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

func main() {
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
