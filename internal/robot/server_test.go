package robot

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	protogate "social/pkg/proto/gate"
	"testing"
	"time"
)

func TestName(t *testing.T) {

	// 连接 gRPC 服务端
	conn, err := grpc.Dial("127.0.0.1:5101", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到服务器: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := protogate.NewServiceClient(conn)

	// 调用双向流 RPC 方法
	stream, err := client.BidirectionalBinaryData(context.Background())
	if err != nil {
		log.Fatalf("调用 gRPC 服务端时发生错误: %v", err)
	}

	// 发送一些请求到服务端
	for i := 0; i < 5; i++ {
		request := &your_proto.YourRequest{
			// 填充请求参数
		}

		err := stream.Send(request)
		if err != nil {
			log.Fatalf("无法发送请求到服务端: %v", err)
		}

		time.Sleep(time.Second)
	}

	// 接收服务端的响应
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break // 服务端流结束
		}
		if err != nil {
			log.Fatalf("接收服务端响应时发生错误: %v", err)
		}

		// 处理服务端的响应
		fmt.Printf("收到服务端响应: %v\n", response)
	}

}
