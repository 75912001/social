package common

import (
	"github.com/pkg/errors"
	"sync"
)

// grpc的双向流,客户端管理器

// ClientStreamMgr 客户端,数据流管理器
type ClientStreamMgr struct {
	clientMap sync.Map //k:client uuid, v:stream
	streamMap sync.Map //k:&stream, v:client uuid
}

var ClientStreamMgrInstance *ClientStreamMgr

func init() {
	ClientStreamMgrInstance = new(ClientStreamMgr)
}

// Add 添加
func (p *ClientStreamMgr) Add(clientUUID any, stream any) error {
	if clientUUID == nil || stream == nil {
		return errors.New("clientUUID or stream is nil")
	}
	if p.IsClientUUIDExists(clientUUID) {
		return errors.New("clientUUID is already exists")
	}
	if p.IsStreamExists(stream) {
		return errors.New("stream is already exists")
	}
	//存储
	p.clientMap.Store(clientUUID, stream)
	p.streamMap.Store(stream, clientUUID)
	return nil
}

// Del 删除
func (p *ClientStreamMgr) Del(stream any) error {
	if stream == nil {
		return errors.New("stream is nil")
	}

	if !p.IsStreamExists(stream) {
		return errors.New("stream is not exists")
	}
	c, _ := p.streamMap.Load(stream)
	p.clientMap.Delete(c)
	p.streamMap.Delete(stream)

	return nil
}

// TODO 通过clientUUID获取stream

// IsClientUUIDExists 是否存在
func (p *ClientStreamMgr) IsClientUUIDExists(clientUUID any) bool {
	_, ok := p.clientMap.Load(clientUUID)
	return ok
}

// IsStreamExists 是否存在
func (p *ClientStreamMgr) IsStreamExists(stream any) bool {
	_, ok := p.streamMap.Load(stream)
	return ok
}
