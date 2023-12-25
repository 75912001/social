package gate

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	libruntime "social/lib/runtime"
	pkgmsg "social/pkg/msg"
	pkgproto "social/pkg/proto"
	protofriend "social/pkg/proto/friend"
)

func send2Friend(ctx context.Context, stream protofriend.FriendService_BidirectionalBinaryDataServer, cmd uint32, resultID uint32, pb proto.Message) error {
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
