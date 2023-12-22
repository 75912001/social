package gate

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	liberror "social/lib/error"
	pkgconsts "social/pkg/consts"
	pkgmsg "social/pkg/msg"
	protofriend "social/pkg/proto/friend"
	protogate "social/pkg/proto/gate"
	"strconv"
)

// Router 消息路由管理器
type Router struct {
}

func (p *Router) Handle(routerKey string, stream protogate.Service_BidirectionalBinaryDataServer, header *pkgmsg.Header, data []byte) error {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(pkgconsts.ShardKey, routerKey))

	if protofriend.MessageMinCMD < header.MessageID && header.MessageID < protofriend.MessageMaxCMD {
		//friend
		...找到一个负载较少的friend服务...
		...将请求投递过去...
		stream.SetHeader()
		send2Friend()
		return nil
	}

	return liberror.MessageIDNonExistent
}


func SetTeamShardKeyToSession(ctx context.Context, key uint64) (context.Context, error) {
	k := consts.SessionPrefix + "team_id"
	v := fmt.Sprintf("%v", key)
	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	md.Set(k, v)
	return metadata.NewOutgoingContext(ctx, md), errors.WithStack(grpc.SetHeader(ctx, metadata.Pairs(k, v)))
}

func GetTeamShardKeyFromSession(ctx context.Context) (key uint64, ok bool) {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		if ks := md.Get(consts.SessionPrefix + "team_id"); len(ks) != 0 {
			i, err := strconv.ParseUint(ks[0], 10, 64)
			if err != nil {
				return 0, false
			}
			return i, true
		}
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	ks := md.Get(consts.SessionPrefix + "team_id")
	if len(ks) == 0 {
		return 0, false
	}
	i, err := strconv.ParseUint(ks[0], 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

/////////////////////////////////////

package main

import (
"context"
"log"
"net"

"google.golang.org/grpc"
"google.golang.org/grpc/metadata"
"google.golang.org/grpc/peer"
)

const (
	shardKey = "shard-key" // 分片键在元数据中的键名
)

type server struct{}

func (s *server) YourBiDiStreamMethod(stream YourService_YourBiDiStreamMethodServer) error {
	// 从上下文中提取元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("no metadata in context")
	}

	// 从元数据中获取分片键
	shardValues := md.Get(shardKey)
	if len(shardValues) == 0 {
		return errors.New("shard key not found in metadata")
	}
	shardKeyValue := shardValues[0]

	// 使用分片键处理业务逻辑
	// ...

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// 注册服务
	// ...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}



//////////////////////////////////////


package main

import (
"context"
"log"

"google.golang.org/grpc"
"google.golang.org/grpc/metadata"
)

const (
	shardKey     = "shard-key" // 分片键在元数据中的键名
	shardKeyValue = "your-shard-value" // 分片键的值
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewYourServiceClient(conn)

	// 创建带有分片键的上下文
	md := metadata.Pairs(shardKey, shardKeyValue)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 创建双向流
	stream, err := client.YourBiDiStreamMethod(ctx)
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	// 使用该流发送和接收消息
	// ...
}



注释说明
服务端：

YourBiDiStreamMethod：这是一个双向流方法的示例，你需要将它替换为你自己的方法。
从上下文中提取元数据，并从中获取分片键。
客户端：

创建一个带有分片键的元数据。
使用这个元数据创建一个新的上下文。
使用该上下文创建双向流。
这个示例展示了如何在gRPC的双向流中使用上下文来传递分片键。你需要根据你的具体需求和gRPC服务定义来调整这些代码。