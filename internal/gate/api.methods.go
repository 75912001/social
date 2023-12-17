package gate

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	pkggrpcstream "social/pkg/grpcstream"
	pkgmsg "social/pkg/msg"
	protogate "social/pkg/proto/gate"
)

func (s *APIServer) BidirectionalBinaryData(stream protogate.Service_BidirectionalBinaryDataServer) error {
	for {
		request, err := stream.Recv()
		if err != nil { // 错误处理 todo menglingchao 此处处理可做成一个函数
			gate.LogMgr.Fatal(err, liberror.Link, stream, libruntime.Location())
			// 使用 status.FromError 函数获取 gRPC 状态
			st, ok := status.FromError(err)
			if ok {
				// 获取错误代码
				code := st.Code()
				// 获取错误消息
				message := st.Message()
				gate.LogMgr.Fatal(code, message, libruntime.Location())
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
				gate.LogMgr.Fatal(st, ok, stream, libruntime.Location())
			}
			err2 := pkggrpcstream.GetInstance().Del(stream)
			if err2 != nil {
				gate.LogMgr.Fatal(err, err2, libruntime.Location())
				return errors.WithMessage(err, err2.Error())
			}
			return err
		}
		data := request.GetData()
		//获取数据-二进制
		if uint32(len(data)) < pkgmsg.GProtoHeadLength {
			gate.LogMgr.Warn(liberror.PacketHeaderLength, libruntime.Location())
			return errors.WithMessage(liberror.PacketHeaderLength, libruntime.Location())
		}
		header := &pkgmsg.Header{}
		header.Unmarshal(data[:pkgmsg.GProtoHeadLength])
		gate.LogMgr.Trace(header.String())
		if gate.CanHandle(header.MessageID) {
			err = gate.Handle(stream, header, data[pkgmsg.GProtoHeadLength:])
			if err != nil {
				gate.LogMgr.Warn(err, libruntime.Location())
			}
		} else { //非gate的消息,交给router处理
			err = gate.router.Handle(stream, header, data)
			if err != nil {
				gate.LogMgr.Warn(err, libruntime.Location())
			}
		}
	}
	//return nil
}
