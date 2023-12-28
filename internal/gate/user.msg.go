package gate

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	libruntime "social/lib/runtime"
	pkgmsg "social/pkg/msg"
	pkgproto "social/pkg/proto"
	protogate "social/pkg/proto/gate"
)

func send2User(stream protogate.GateService_BidirectionalBinaryDataServer, cmd uint32, resultID uint32, pb proto.Message) error {
	var err error
	res := pkgmsg.NewPacket(cmd, resultID, pb)
	sendData := pkgproto.BinaryData{}
	sendData.Data, err = res.Marshal()
	if err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	if err = stream.Send(&sendData); err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	return nil
}

// 总线-用户管理器-添加用户
type busUserMgrAddUser struct {
	key    string
	stream grpc.ServerStream
}

type busUserMgrDelUser struct {
	key    string
	stream grpc.ServerStream
}
