package gate

import (
	"github.com/golang/protobuf/proto"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	pkgmsg "social/pkg/msg"
	pkgproto "social/pkg/proto"
	protogate "social/pkg/proto/gate"
)

// CanHandle 是否可以处理的命令 todo menglingchao
func (p *Server) CanHandle(messageID uint32) bool {
	// 在这里编写判断逻辑
	// 判断cmd是否符合自己可以处理的消息类型
	// 如果符合，返回true；否则，返回false
	//if protogate.User2GateMessageMinCMD < header.MessageID && header.MessageID < protogate.User2GateMessageMaxCMD { //gate的消息

	return false
}

func (p *Server) Handle(stream protogate.Service_BidirectionalBinaryDataServer, header *pkgmsg.Header, data []byte) error {
	p.LogMgr.Trace(pkgproto.CMDMap[header.MessageID], stream, header, data)
	switch header.MessageID {
	case protogate.RegisterReqCMD:
		var req protogate.RegisterReq
		err := proto.Unmarshal(data, &req)
		if err != nil {
			p.LogMgr.Error(err, libruntime.Location())
			return err
		}
		p.LogMgr.Trace(req.String())
		// 处理 RegisterReq todo menglingchao
		//err = pkggrpcstream.GetInstance().Add(req.GetServiceKey().GetServiceID(), stream)
		//组包
		resPacket := pkgmsg.NewPacket(protogate.RegisterResCMD, 0, &protogate.RegisterRes{})
		sendData := pkgproto.BinaryData{}
		sendData.Data, err = resPacket.Marshal()
		if err != nil {
			p.LogMgr.Warn(err, libruntime.Location())
			return err
		}
		//回包
		if err = stream.Send(&sendData); err != nil {
			p.LogMgr.Error(err, libruntime.Location())
			return err
		}
	case protogate.LogoutReqCMD:
		var req protogate.LogoutReq
		err := proto.Unmarshal(data, &req)
		if err != nil {
			p.LogMgr.Error(err, libruntime.Location())
			return err
		}
		p.LogMgr.Trace(req.String())
		// 处理 LogoutReq todo menglingchao
		//删除client
		//err = pkggrpcstream.GetInstance().Del(stream)
		//if err != nil {
		//	liblog.GetInstance().Fatal(err, libruntime.Location())
		//	return err
		//}
		//组包
		resPacket := pkgmsg.NewPacket(protogate.LogoutResCMD, 0, &protogate.LogoutRes{})
		sendData := pkgproto.BinaryData{}
		sendData.Data, err = resPacket.Marshal()
		if err != nil {
			p.LogMgr.Warn(err, libruntime.Location())
			return err
		}
		//回包
		if err = stream.Send(&sendData); err != nil {
			p.LogMgr.Error(err, libruntime.Location())
			return err
		}
	default:
		p.LogMgr.Error(liberror.MessageIDNonExistent, libruntime.Location())
		return liberror.MessageIDNonExistent
	}
	return nil
}
