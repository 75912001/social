package gate

import (
	pkgmsg "social/pkg/msg"
	protogate "social/pkg/proto/gate"
)

// Router 消息路由管理器
type Router struct {
}

func (p *Router) Handle(stream protogate.Service_BidirectionalBinaryDataServer, header *pkgmsg.Header, body []byte) error {
	//TODO
	return nil
}
