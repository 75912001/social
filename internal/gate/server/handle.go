package server

import (
	liberror "social/lib/error"
	liblog "social/lib/log"
	pkggrpcstream "social/pkg/grpcstream"
	pkgmsg "social/pkg/msg"
	pkgproto "social/pkg/proto"
	protogate "social/pkg/proto/gate"
)

func Handle(stream protogate.Service_BidirectionalBinaryDataServer, data []byte) error {
	packet := pkgmsg.Packet{}
	err := packet.Unmarshal(data)
	if err != nil {
		liblog.GetInstance().Fatal(err, libutil.GetCodeLocation(1))
		return err
	}
	liblog.GetInstance().Trace(packet.Header, packet.Message)

	//switch m := packet.Message.(type) {
	//case *protogate.RegisterReq:
	//	fmt.Println("", m.ServiceKey)
	//case *protogate.LogoutReq:
	//	fmt.Println("")
	//default:
	//	// 处理未知类型或其他情况
	//	fmt.Println("")
	//}
	switch packet.Header.MessageID {
	case protogate.RegisterReq_CMD:
		//var req pkgproto.ServiceKey //
		var req protogate.RegisterReq
		err := proto.Unmarshal(data[pkgmsg.GProtoHeadLength:], &req)
		liblog.GetInstance().Trace("Received RegisterReq:", req.String(), stream)
		// 处理 RegisterReq
		//err = pkggrpcstream.GetInstance().Add(req.GetServiceKey().GetServiceID(), stream)
		if err != nil {
			liblog.GetInstance().Fatal(err, libutil.GetCodeLocation(1))
			return err
		}
		//组包
		resPacket := pkgmsg.Packet{
			Header: pkgmsg.Header{
				MessageID: protogate.RegisterRes_CMD,
				ResultID:  0,
			},
			Message: &protogate.RegisterRes{},
		}
		sendData := pkgproto.BinaryData{}
		sendData.Data, err = resPacket.Marshal()
		//回包
		if err = stream.Send(&sendData); err != nil {
			liblog.GetInstance().Fatal(err, libutil.GetCodeLocation(1))
			return err
		}
	case protogate.LogoutReq_CMD:
		var req protogate.LogoutReq
		err := proto.Unmarshal(data[pkgmsg.GProtoHeadLength:], &req)
		liblog.GetInstance().Trace("Received LogoutReq:", req.String(), stream)
		// 处理 LogoutReq
		//删除client
		err = pkggrpcstream.GetInstance().Del(stream)
		if err != nil {
			liblog.GetInstance().Fatal(err, libutil.GetCodeLocation(1))
			return err
		}
		//组包
		resPacket := pkgmsg.Packet{
			Header: pkgmsg.Header{
				MessageID: protogate.LogoutRes_CMD,
				ResultID:  0,
			},
			Message: &protogate.LogoutRes{},
		}
		sendData := pkgproto.BinaryData{}
		sendData.Data, err = resPacket.Marshal()
		//回包
		if err = stream.Send(&sendData); err != nil {
			liblog.GetInstance().Fatal(err, libutil.GetCodeLocation(1))
			return err
		}
	default:
		liblog.GetInstance().Fatal(liberror.MessageIDNonExistent, libutil.GetCodeLocation(1))
		return liberror.MessageIDNonExistent
	}
	return nil
}
