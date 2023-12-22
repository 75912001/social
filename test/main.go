package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pkgconsts "social/pkg/consts"
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

// GetShardKey 获取分片键.
func GetShardKey(serverStream grpc.ServerStream) (string, error) {
	md, ok := metadata.FromIncomingContext(serverStream.Context())
	if !ok {
		return "", errors.WithStack(errors.New("metadata is not found"))
	}
	values := md.Get(pkgconsts.ShardKey)
	if len(values) == 0 {
		return "", errors.WithStack(errors.New("shard key is not found"))
	}
	return values[0], nil
}

// SetShardKey 设置分片键.
// 保留原有serverStream中context的元数据,并添加一个元数据.
func SetShardKey(serverStream grpc.ServerStream, shardKey string) error {
	//md, ok := metadata.FromOutgoingContext(serverStream.Context())
	//if !ok {
	//	return errors.WithMessagef(nil, "failed to get metadata %v", libruntime.Location())
	//}
	//md.Set(pkgconsts.ShardKey, shardKey)
	//
	//// 创建一个新的元数据并设置分片键
	//newMD := metadata.Pairs(pkgconsts.ShardKey, shardKey)
	//return errors.WithStack(serverStream.SetHeader(newMD))

	// 创建一个新的元数据并设置分片键
	newMD := metadata.Pairs(pkgconsts.ShardKey, shardKey)

	// 使用新的元数据设置头部
	return errors.WithStack(serverStream.SetHeader(newMD))
}
