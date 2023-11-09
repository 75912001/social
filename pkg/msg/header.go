package msg

import (
	"context"
	xrpb "dawn-server/impl/xr/lib/pb"
	"encoding/binary"
	"fmt"
	xrutil "social/pkg/lib/util"

	"github.com/pkg/errors"
)

// 包头

// GProtoHeadLength 包头长度
const GProtoHeadLength uint32 = 8

// Header 协议包头
type Header struct {
	MessageID uint32 //消息号
	ResultID  uint32 //结果id
}

func (p *Header) Pack(data []byte) {
	binary.LittleEndian.PutUint32(data[0:], p.MessageID)
	binary.LittleEndian.PutUint32(data[4:], p.ResultID)
}

func (p *Header) String() string {
	return fmt.Sprintf("MessageID:%#x, ResultID:%#x", p.MessageID, p.ResultID)
}

func (p *Header) Unpack(data []byte) {
	p.MessageID = binary.LittleEndian.Uint32(data[0:4])
	p.ResultID = binary.LittleEndian.Uint32(data[4:8])
}

type Packet struct{}

// Marshal 序列化 数据包 (发送给Room)
func (p *Packet) Marshal(packet *xrpb.UnserializedPacket) ([]byte, error) {
	var data []byte
	ph := packet.Header.(*SSHeader)
	if !packet.IsPassThrough() {
		var err error
		data, err = MarshalMsgSS(packet.CTX, packet.Message)
		if err != nil {
			return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
		}
	} else {
		data = packet.PassThroughBody
	}

	packetLength := GRoomProtoHeadLength + uint32(len(data))
	ph.PacketLength = packetLength
	buf := make([]byte, packetLength)

	packet.Header.Pack(buf)
	copy(buf[GRoomProtoHeadLength:packetLength], data)
	return buf, nil
}

// Unmarshal 反序列化
//
//	data:完整包数据 包头+包体
func (p *Packet) Unmarshal(data []byte) (header interface{}, ctx context.Context, pbData []byte, err error) {
	h := &SSHeader{}
	h.Unpack(data)

	ctx, pbData, err = UnmarshalMsgSS(data[GRoomProtoHeadLength:])
	if err != nil {
		return nil, nil, nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	return h, ctx, pbData, nil
}
