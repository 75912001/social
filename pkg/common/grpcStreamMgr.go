package common

import (
	"github.com/pkg/errors"
	xrerror "social/pkg/lib/error"
	xrutil "social/pkg/lib/util"
	"sync"
)

// grpc的双向流,客户端管理器

// GrpcStreamMgr 客户端,数据流管理器
type GrpcStreamMgr struct {
	clientMap sync.Map //k:client uuid, v:stream
	streamMap sync.Map //k:&stream, v:client uuid
}

var (
	instance *GrpcStreamMgr
	once     sync.Once
)

// GetGrpcStreamMgrInstance 获取实例
func GetGrpcStreamMgrInstance() *GrpcStreamMgr {
	once.Do(func() {
		instance = new(GrpcStreamMgr)
	})
	return instance
}

// Add 添加
func (p *GrpcStreamMgr) Add(clientUUID any, stream any) error {
	if clientUUID == nil || stream == nil {
		return errors.WithMessage(xrerror.Param.WithExtraMessage("clientUUID or stream is nil"), xrutil.GetCodeLocation(1).String())
	}
	if p.IsClientUUIDExists(clientUUID) {
		return errors.WithMessage(xrerror.Exists.WithExtraMessage("clientUUID is already exists"), xrutil.GetCodeLocation(1).String())
	}
	if p.IsStreamExists(stream) {
		return errors.WithMessage(xrerror.Exists.WithExtraMessage("stream is already exists"), xrutil.GetCodeLocation(1).String())
	}
	//存储
	p.clientMap.Store(clientUUID, stream)
	p.streamMap.Store(stream, clientUUID)
	return nil
}

// Del 删除
func (p *GrpcStreamMgr) Del(stream any) error {
	if stream == nil {
		return errors.WithMessage(xrerror.Param.WithExtraMessage("stream is nil"), xrutil.GetCodeLocation(1).String())
	}
	if !p.IsStreamExists(stream) {
		return errors.WithMessage(xrerror.NonExistent.WithExtraMessage("stream is not exists"), xrutil.GetCodeLocation(1).String())
	}
	c, _ := p.streamMap.Load(stream)
	p.clientMap.Delete(c)
	p.streamMap.Delete(stream)

	return nil
}

// IsClientUUIDExists 是否存在
func (p *GrpcStreamMgr) IsClientUUIDExists(clientUUID any) bool {
	_, ok := p.clientMap.Load(clientUUID)
	return ok
}

// IsStreamExists 是否存在
func (p *GrpcStreamMgr) IsStreamExists(stream any) bool {
	_, ok := p.streamMap.Load(stream)
	return ok
}
