package gate

import (
	liberror "social/lib/error"
	pkgmsg "social/pkg/msg"
	protogate "social/pkg/proto/gate"
)

// Router 消息路由管理器
type Router struct {
}

func (p *Router) Handle(routerKey string, stream protogate.GateService_BidirectionalBinaryDataServer, header *pkgmsg.Header, data []byte) error {
	//ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(pkgconsts.ShardKey, routerKey))
	//
	//if protofriend.MessageMinCMD < header.MessageID && header.MessageID < protofriend.MessageMaxCMD {
	//	//friend
	//	//...找到一个负载较少的friend服务...
	//	//...将请求投递过去...
	//	//stream.SetHeader()
	//	//send2Friend()
	//	return nil
	//}

	return liberror.MessageIDNonExistent
}
