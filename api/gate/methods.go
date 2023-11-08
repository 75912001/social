package gate

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"social/pkg/common"
	xrerror "social/pkg/lib/error"
	xrutil "social/pkg/lib/util"
	"social/pkg/proto/gate"
)

func (s *Server) BidirectionalStreamingMethod(stream gate.Service_BidirectionalStreamingMethodServer) error {
	for {
		request, err := stream.Recv()
		if err != nil {
			log.Fatalln(err, xrerror.Link, stream, xrutil.GetCodeLocation(1))
			// 使用 status.FromError 函数获取 gRPC 状态
			st, ok := status.FromError(err)
			if ok {
				// 获取错误代码
				code := st.Code()
				// 获取错误消息
				message := st.Message()
				log.Fatalln(code, message, xrutil.GetCodeLocation(1))

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
				log.Fatalln(st, ok, stream, xrutil.GetCodeLocation(1))
			}
			e := common.GetGrpcStreamMgrInstance().Del(stream)
			if e != nil {
				log.Fatalln(err, e, xrutil.GetCodeLocation(1))
				return errors.WithMessage(err, e.Error())
			}
			return err
		}

		// 根据请求类型选择处理逻辑
		switch req := request.GetRequest().(type) {
		case *gate.Request_RegisterReq:
			log.Println("Received Request_RegisterReq:", req.RegisterReq.GetServiceKey())
			// 处理 RegisterReq 并生成响应
			res := &gate.Response{
				Response: &gate.Response_RegisterRes{
					RegisterRes: &gate.RegisterRes{},
				},
			}
			err = common.GetGrpcStreamMgrInstance().Add(req.RegisterReq.GetServiceKey().GetServiceID(), stream)
			if err != nil {
				log.Fatalln(err, stream, xrutil.GetCodeLocation(1))
				return err
			}
			if err = stream.Send(res); err != nil {
				log.Fatalln(err, stream, xrutil.GetCodeLocation(1))
				return err
			}
		case *gate.Request_LogoutReq:
			log.Println("Received Request_LogoutReq:", req.LogoutReq.GetServiceKey())
			// 处理 RequestTypeB 并生成响应
			response := &gate.Response{
				Response: &gate.Response_LogoutRes{
					LogoutRes: &gate.LogoutRes{},
				},
			}
			//删除client
			err = common.GetGrpcStreamMgrInstance().Del(stream)
			if err != nil {
				log.Fatalln(err, stream, xrutil.GetCodeLocation(1))
				return err
			}
			if err = stream.Send(response); err != nil {
				log.Fatalln(err, stream, xrutil.GetCodeLocation(1))
				return err
			}
		default:

			log.Fatalln(xrerror.MessageIDNonExistent, xrutil.GetCodeLocation(1))
			return nil
		}
	}
}

func (s *Server) ForwardBinaryData(stream gate.Service_ForwardBinaryDataServer) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			// 客户端关闭流，结束循环
			return nil
		}
		if err != nil {
			log.Printf("Error receiving data: %v", err)
			return err
		}

		// 在这里可以根据需要处理接收到的二进制数据
		// 在本示例中，我们将接收到的数据直接打印出来
		log.Printf("Received binary data: %v", data.Data)

		// 在这里可以根据需要处理数据，然后将数据发送回客户端
		// 在本示例中，我们将接收到的数据原样发送回客户端
		if err := stream.Send(data); err != nil {
			log.Printf("Error sending data: %v", err)
			return err
		}
	}
}
