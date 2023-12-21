package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

func main() {
	{
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("router-key", "client"))

		// 从上下文中获取路由信息
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			fmt.Println("")
		}
		routerKey := md.Get("router-key")
		fmt.Printf("Router Key: %v\n", routerKey)

		// 在响应中添加路由信息
		//routerKey = append(routerKey, "server") // 这里只是一个简单的演示，实际应用中需要更复杂的逻辑
		//stream.Send(&HelloReply{Message: fmt.Sprintf("Hello, %s!", msg.GetBody()), RouterKey: routerKey})
	}
}
