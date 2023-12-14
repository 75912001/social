package router

import (
	pkgmsg "social/pkg/msg"
	protogate "social/pkg/proto/gate"
	"sync"
)

var (
	instance *mgr
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *mgr {
	once.Do(func() {
		instance = new(mgr)
	})
	return instance
}

// Mgr 消息路由管理器
type mgr struct {
}

func (p *mgr) Handle(stream protogate.Service_BidirectionalBinaryDataServer, header *pkgmsg.Header, body []byte) error {
	//TODO
	return nil
}
