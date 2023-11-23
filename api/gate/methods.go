package gate

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"social/internal/gate/router"
	"social/pkg/grpcstream"
	xrerror "social/pkg/lib/error"
	xrlog "social/pkg/lib/log"
	xrpb "social/pkg/lib/pb"
	xrutil "social/pkg/lib/util"
	"social/pkg/msg"
	pkgproto "social/pkg/proto"
	"social/pkg/proto/gate"
)

func (s *Server) BidirectionalStreamingMethod(stream gate.Service_BidirectionalBinaryDataServer) error {
	defer func() {

	}()
	for {
		request, err := stream.Recv()
		if err != nil {
			xrlog.GetInstance().Fatal(err, xrerror.Link, stream, xrutil.GetCodeLocation(1))
			// 使用 status.FromError 函数获取 gRPC 状态
			st, ok := status.FromError(err)
			if ok {
				// 获取错误代码
				code := st.Code()
				// 获取错误消息
				message := st.Message()
				xrlog.GetInstance().Fatal(code, message, xrutil.GetCodeLocation(1))

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
				xrlog.GetInstance().Fatal(st, ok, stream, xrutil.GetCodeLocation(1))
			}
			err2 := grpcstream.GetInstance().Del(stream)
			if err2 != nil {
				xrlog.GetInstance().Fatal(err, err2, xrutil.GetCodeLocation(1))
				return errors.WithMessage(err, err2.Error())
			}
			return err
		}
		{ //获取数据-二进制
			data := request.GetData()
			if uint32(len(data)) < msg.GProtoHeadLength {
				//todo menglingchao 消息长度不足,断开链接 grpc
				return errors.WithMessage(xrerror.Packet, xrutil.GetCodeLocation(1).String())
			}

			header := &msg.Header{
				MessageID: 0,
				ResultID:  0,
			}
			header.Unpack(data)
			xrlog.GetInstance().Trace(header.String())
			//todo menglingchao 按照CMD来 处理/分发 数据包...
			//return xxx
			err = router.GetInstance().Handle(header, data)
			if err != nil {
				xrlog.GetInstance().Fatal(err, xrutil.GetCodeLocation(1))
				return err
			}
		}
		err = handle(stream, request.GetData())
		if err != nil {
			xrlog.GetInstance().Fatal(err, xrutil.GetCodeLocation(1))
			return err
		}
	}
	return nil
}

func handle(stream gate.Service_BidirectionalBinaryDataServer, data []byte) error {
	unserializedPacket, err := msg.Unmarshal(data)
	if err != nil {
		xrlog.GetInstance().Fatal(err, xrutil.GetCodeLocation(1))
		return err
	}
	xrlog.GetInstance().Trace(unserializedPacket.Header, unserializedPacket.Message)

	// 根据请求类型选择处理逻辑
	switch req := unserializedPacket.Message.(type) {
	case *gate.RegisterReq:
		xrlog.GetInstance().Trace("Received RegisterReq:", req.String(), stream)
		// 处理 RegisterReq
		err = grpcstream.GetInstance().Add(req.GetServiceKey().GetServiceID(), stream)
		if err != nil {
			xrlog.GetInstance().Fatal(err, xrutil.GetCodeLocation(1))
			return err
		}
		//组包
		resPacket := &xrpb.UnserializedPacket{
			Header: &msg.Header{
				MessageID: gate.RegisterRes_CMD,
				ResultID:  0,
			},
			Message: &gate.RegisterRes{},
		}
		sendData := pkgproto.BinaryData{}
		sendData.Data, err = msg.Marshal(resPacket)
		//回包
		if err = stream.Send(&sendData); err != nil {
			xrlog.GetInstance().Fatal(err, xrutil.GetCodeLocation(1))
			return err
		}
	case *gate.LogoutReq:
		xrlog.GetInstance().Trace("Received LogoutReq:", req.String(), stream)
		// 处理 LogoutReq
		//删除client
		err = grpcstream.GetInstance().Del(stream)
		if err != nil {
			xrlog.GetInstance().Fatal(err, xrutil.GetCodeLocation(1))
			return err
		}
		//组包
		resPacket := &xrpb.UnserializedPacket{
			Header: &msg.Header{
				MessageID: gate.LogoutRes_CMD,
				ResultID:  0,
			},
			Message: &gate.LogoutRes{},
		}
		sendData := pkgproto.BinaryData{}
		sendData.Data, err = msg.Marshal(resPacket)
		//回包
		if err = stream.Send(&sendData); err != nil {
			xrlog.GetInstance().Fatal(err, xrutil.GetCodeLocation(1))
			return err
		}
	default:
		xrlog.GetInstance().Fatal(xrerror.MessageIDNonExistent, xrutil.GetCodeLocation(1))
		return xrerror.MessageIDNonExistent
	}
	return nil
}
