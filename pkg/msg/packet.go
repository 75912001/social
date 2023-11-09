package msg

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	xrpb "social/pkg/lib/pb"
	xrutil "social/pkg/lib/util"
)

type Packet struct{}

// Marshal 序列化 数据包
func (p *Packet) Marshal(packet *xrpb.UnserializedPacket) ([]byte, error) {
	var data []byte
	var err error
	if data, err = proto.Marshal(packet.Message); nil != err && proto.ErrNil != err {
		return nil, errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	}

	packetLength := GProtoHeadLength + uint32(len(data))
	buf := make([]byte, packetLength)

	packet.Header.Pack(buf)
	copy(buf[GProtoHeadLength:packetLength], data)
	return buf, nil
}

// Unmarshal 反序列化
//
//	data:完整包数据 包头+包体
func (p *Packet) Unmarshal(data []byte) (header xrpb.IHeader, msg proto.Message, err error) {
	header = &Header{}
	header.Unpack(data)
	if err = proto.Unmarshal(data[GProtoHeadLength:], msg); nil != err {
		return nil, nil, errors.WithMessagef(err, "%v data:%v", xrutil.GetCodeLocation(1).String(), data)
	}
	return header, msg, nil
}
