package gate

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	pkgcommon "social/pkg/common"
	pkgmsg "social/pkg/msg"
	pkgproto "social/pkg/proto"
	protogate "social/pkg/proto/gate"
)

// CanHandle 是否可以处理的命令
// 判断cmd是否符合自己可以处理的消息类型
// 如果符合，返回true；否则，返回false
func (p *Gate) CanHandle(messageID uint32) bool {
	if protogate.MessageMinCMD < messageID && messageID < protogate.MessageMaxCMD { //gate的消息
		return true
	}
	return false
}

func (p *Gate) Handle(stream protogate.Service_BidirectionalBinaryDataServer, header *pkgmsg.Header, body []byte) error {
	p.LogMgr.Trace(pkgproto.CMDMap[header.MessageID], stream, header, body)
	switch header.MessageID {
	case protogate.RegisterReqCMD:
		var req protogate.RegisterReq
		err := proto.Unmarshal(body, &req)
		if err != nil {
			p.LogMgr.Error(err, libruntime.Location())
			return err
		}
		p.LogMgr.Trace(req.String())
		//todo 从redis中验证token...

		// 处理 RegisterReq
		key := pkgcommon.GenerateServiceKey(req.ServiceKey.ZoneID, req.ServiceKey.ServiceName, req.ServiceKey.ServiceID)
		user := p.userMgr.Find(key)
		if user != nil { //注册过,返回错误码
			err = send2User(stream, protogate.RegisterResCMD, liberror.Duplicate.Code, &protogate.RegisterRes{})
			return errors.WithMessagef(err, "user already registered %v", libruntime.Location())
		}
		user = p.userMgr.FindByStream(stream)
		if user != nil { //注册过,返回错误码
			err = send2User(stream, protogate.RegisterResCMD, liberror.Duplicate.Code, &protogate.RegisterRes{})
			return errors.WithMessagef(err, "user already registered %v", libruntime.Location())
		}
		//注册
		user = p.userMgr.SpawnUser(key, stream)
		//注册-成功
		err = send2User(stream, protogate.RegisterResCMD, 0, &protogate.RegisterRes{})
		if err != nil {
			return errors.WithMessagef(err, "send2User %v", libruntime.Location())
		}
	case protogate.LogoutReqCMD:
		var req protogate.LogoutReq
		err := proto.Unmarshal(body, &req)
		if err != nil {
			p.LogMgr.Error(err, libruntime.Location())
			return err
		}
		p.LogMgr.Trace(req.String())
		// 处理 LogoutReq
		user := p.userMgr.FindByStream(stream)
		if user == nil { //未能找到
			err = send2User(stream, protogate.LogoutResCMD, liberror.NonExistent.Code, &protogate.LogoutRes{})
			if err != nil {
				return errors.WithMessagef(err, "send2User %v", libruntime.Location())
			}
		}
		//删除用户
		p.userMgr.DeleteUser(user.key)
		err = send2User(stream, protogate.LogoutResCMD, 0, &protogate.LogoutRes{})
		if err != nil {
			return errors.WithMessagef(err, "send2User %v", libruntime.Location())
		}
	default:
		p.LogMgr.Error(liberror.MessageIDNonExistent, libruntime.Location())
		return liberror.MessageIDNonExistent
	}
	return nil
}
