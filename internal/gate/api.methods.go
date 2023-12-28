package gate

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	pkggrpcstream "social/pkg/grpcstream"
	pkgmsg "social/pkg/msg"
	protogate "social/pkg/proto/gate"
)

func RecvError(err error, fun func()) {
	// 使用 status.FromError 函数获取 gRPC 状态
	st, ok := status.FromError(err)
	if ok {
		// 获取错误代码
		code := st.Code()
		// 获取错误消息
		message := st.Message()
		app.LogMgr.Error(code, message, libruntime.Location())
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
		app.LogMgr.Error(st, ok, libruntime.Location())
	}
	fun()
}

func (s *APIServer) BidirectionalBinaryData(stream protogate.GateService_BidirectionalBinaryDataServer) error {
	for {
		request, err := stream.Recv()
		if err != nil { // 错误处理
			app.LogMgr.Error(err, liberror.Link, stream, libruntime.Location())
			RecvError(err, func() {
				err2 := pkggrpcstream.GetInstance().Del(stream)
				if err2 != nil {
					app.LogMgr.Error(err, err2, libruntime.Location())
				}
			})
			return err
		}
		data := request.GetData()
		//获取数据-二进制
		if uint32(len(data)) < pkgmsg.GProtoHeadLength {
			app.LogMgr.Error(liberror.PacketHeaderLength, libruntime.Location())
			continue
		}
		rawHeader := &pkgmsg.Header{}
		rawHeader.Unmarshal(data[:pkgmsg.GProtoHeadLength])
		app.LogMgr.Trace(rawHeader.String())
		message := app.userPBFunMgr.Find(rawHeader.MessageID)
		if message != nil {
			rawMessage, err := message.Unmarshal(data[pkgmsg.GProtoHeadLength:])
			if err != nil {
				app.LogMgr.Error(err, libruntime.Location())
				continue
			}
			app.LogMgr.Trace(rawMessage)
			err = message.Handler(rawHeader, rawMessage, stream)
			if err != nil {
				app.LogMgr.Warn(err, libruntime.Location())
			}
			continue
		}
		if app.CanHandle(rawHeader.MessageID) {
			err = app.Handle(stream, rawHeader, data[pkgmsg.GProtoHeadLength:])
			if err != nil {
				app.LogMgr.Warn(err, libruntime.Location())
			}
			continue
		}
		//非gate的消息,交给router处理
		user := app.userMgr.FindByStream(stream)
		if user == nil { //未注册
			app.LogMgr.Error(liberror.Unregistered, libruntime.Location())
			continue
			//return errors.WithMessage(liberror.Unregistered, "user not registered")
		}
		err = app.router.Handle(user.key, stream, rawHeader, data)
		if err != nil {
			app.LogMgr.Warn(err, libruntime.Location())
		}
	}
}
