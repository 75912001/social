package gate

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(pkgconsts.RouterKey, routerKey))

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