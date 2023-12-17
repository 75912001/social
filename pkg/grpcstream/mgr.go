package grpcstream

import (
	"github.com/pkg/errors"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	"sync"
)

// grpc的双向流,管理器-管理客户端

// mgr 客户端,数据流管理器
type mgr struct {
	clientMap sync.Map //k:client uuid, v:stream
	streamMap sync.Map //k:&stream, v:client uuid
}

var (
	instance *mgr
	once     sync.Once
)

// GetInstance 获取实例
func GetInstance() *mgr {
	once.Do(func() {
		instance = new(mgr)
	})
	return instance
}

// Add 添加
func (p *mgr) Add(clientUUID any, stream any) error {
	if clientUUID == nil || stream == nil {
		return errors.WithMessage(liberror.Param.WithExtraMessage("clientUUID or stream is nil"), libruntime.Location())
	}
	if p.IsClientUUIDExists(clientUUID) {
		return errors.WithMessage(liberror.Exists.WithExtraMessage("clientUUID is already exists"), libruntime.Location())
	}
	if p.IsStreamExists(stream) {
		return errors.WithMessage(liberror.Exists.WithExtraMessage("stream is already exists"), libruntime.Location())
	}
	//存储
	p.clientMap.Store(clientUUID, stream)
	p.streamMap.Store(stream, clientUUID)
	return nil
}

// Del 删除
func (p *mgr) Del(stream any) error {
	if stream == nil {
		return errors.WithMessage(liberror.Param.WithExtraMessage("stream is nil"), libruntime.Location())
	}
	if !p.IsStreamExists(stream) {
		return errors.WithMessage(liberror.NonExistent.WithExtraMessage("stream is not exists"), libruntime.Location())
	}
	c, _ := p.streamMap.Load(stream)
	p.clientMap.Delete(c)
	p.streamMap.Delete(stream)

	return nil
}

// IsClientUUIDExists 是否存在
func (p *mgr) IsClientUUIDExists(clientUUID any) bool {
	_, ok := p.clientMap.Load(clientUUID)
	return ok
}

// IsStreamExists 是否存在
func (p *mgr) IsStreamExists(stream any) bool {
	_, ok := p.streamMap.Load(stream)
	return ok
}
