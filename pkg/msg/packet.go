package msg

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	libpb "social/pkg/lib/pb"
	libutil "social/pkg/lib/util"
)

type Packet struct{}

// Marshal 序列化 数据包
func Marshal(unserializedPacket *libpb.UnserializedPacket) ([]byte, error) {
	headerBuf := unserializedPacket.Header.Pack()

	var messageBuf []byte
	var err error
	if messageBuf, err = proto.Marshal(unserializedPacket.Message); nil != err && proto.ErrNil != err {
		return nil, errors.WithMessagef(err, libutil.GetCodeLocation(1).String())
	}

	packetLength := GProtoHeadLength + uint32(len(messageBuf))
	packetBuf := make([]byte, packetLength)

	copy(packetBuf[0:GProtoHeadLength], headerBuf)
	copy(packetBuf[GProtoHeadLength:], messageBuf)
	return packetBuf, nil
}

// Unmarshal 反序列化
//
//	data:完整包数据 包头+包体
func Unmarshal(data []byte) (unserializedPacket *libpb.UnserializedPacket, err error) {
	var msg proto.Message
	if err = proto.Unmarshal(data[GProtoHeadLength:], msg); nil != err {
		return nil, errors.WithMessagef(err, "%v data:%v", libutil.GetCodeLocation(1).String(), data)
	}
	unserializedPacket = &libpb.UnserializedPacket{
		Header:  &Header{},
		Message: msg,
	}
	unserializedPacket.Header.Unpack(data)
	return unserializedPacket, nil
}
