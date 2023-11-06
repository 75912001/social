package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"social/proto/social_service"
)

func main() {
	// 连接 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := social_service.NewSocialServiceClient(conn)

	// 创建双向流
	stream, err := client.BidirectionalStreamingMethod(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// 发送多个请求
	requests := []*social_service.CommonRequest{
		{
			Request: &social_service.CommonRequest_RegisterRequest{
				RegisterRequest: &social_service.RegisterRequest{
					ServiceKey: &social_service.ServiceKey{
						ZoneID:      1,
						ServiceName: "name-client",
						ServiceID:   1,
					},
				},
			},
		},
		{
			Request: &social_service.CommonRequest_LogoutRequest{
				LogoutRequest: &social_service.LogoutRequest{
					ServiceKey: &social_service.ServiceKey{
						ZoneID:      1,
						ServiceName: "name-client",
						ServiceID:   1,
					},
				},
			},
		},
	}

	for _, req := range requests {
		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
	}

	// 接收多个响应
	for {
		response, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive response: %v", err)
		}

		// 根据响应类型处理响应
		switch resp := response.GetResponse().(type) {
		case *social_service.CommonResponse_RegisterResponse:
			fmt.Printf("Received RegisterResponse: %s\n", resp.RegisterResponse.GetField1())
		case *social_service.CommonResponse_LogoutResponse:
			fmt.Printf("Received RegisterResponse: %d\n", resp.LogoutResponse.GetField2())
		}
	}
}
