package gate

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gaterouter "social/internal/gate/router"
	liberror "social/lib/error"
	liblog "social/lib/log"
	libutil "social/lib/util"
	pkggrpcstream "social/pkg/grpcstream"
	pkgmsg "social/pkg/msg"
	pkgproto "social/pkg/proto"
	protogate "social/pkg/proto/gate"
)

func (s *Server) BidirectionalBinaryData(stream protogate.Service_BidirectionalBinaryDataServer) error {
	defer func() {

	}()
	for {
		request, err := stream.Recv()
		if err != nil {
			liblog.GetInstance().Fatal(err, liberror.Link, stream, libutil.GetCodeLocation(1))
			// 使用 status.FromError 函数获取 gRPC 状态
			st, ok := status.FromError(err)
			if ok {
				// 获取错误代码
				code := st.Code()
				// 获取错误消息
				message := st.Message()
				liblog.GetInstance().Fatal(code, message, libutil.GetCodeLocation(1))

				// 根据错误代码采取不同的处理方式
				switch code {
				case codes.Unavailable:
					// 服务不可用，可能是网络中断
					// 处理网络问题
				case codes.Canceled:
					// 请求被取消
					// 处理取消请求
				case codes.Unknown:
					// 未知错误
					// 处理未知错误
				default:
					// 其他错误
					// 处理其他错误
				}
				// 在处理不同类型的错误后，可以根据需要进行其他操作
			} else {
				liblog.GetInstance().Fatal(st, ok, stream, libutil.GetCodeLocation(1))
			}
			err2 := pkggrpcstream.GetInstance().Del(stream)
			if err2 != nil {
				liblog.GetInstance().Fatal(err, err2, libutil.GetCodeLocation(1))
				return errors.WithMessage(err, err2.Error())
			}
			return err
		}
		//获取数据-二进制
		if uint32(len(request.GetData())) < pkgmsg.GProtoHeadLength {
			//todo menglingchao 消息长度不足,断开链接 grpc
			return errors.WithMessage(liberror.Packet, libutil.GetCodeLocation(1).String())
		}
		header := &pkgmsg.Header{}
		header.Unpack(request.GetData())
		liblog.GetInstance().Trace(header.String())
		//todo menglingchao 按照CMD来 处理/分发 数据包...
		//return xxx
		err = gaterouter.GetInstance().Handle(header, request.GetData())
		if err != nil {
			liblog.GetInstance().Fatal(err, libutil.GetCodeLocation(1))
			return err
		}

		err = handle(stream, request.GetData())
		if err != nil {
			liblog.GetInstance().Fatal(err, libutil.GetCodeLocation(1))
			return err
		}
	}
	return nil
}

func handle(stream protogate.Service_BidirectionalBinaryDataServer, data []byte) error {
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
