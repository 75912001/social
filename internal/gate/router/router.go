package router

import (
	"social/pkg/msg"
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

func (p *mgr) Handle(header *msg.Header, data []byte) error {
	//TODO
	return nil
}
