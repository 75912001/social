package msg

import (
	"encoding/binary"
	"fmt"
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

func (p *Header) Unpack(data []byte) {
	p.MessageID = binary.LittleEndian.Uint32(data[0:4])
	p.ResultID = binary.LittleEndian.Uint32(data[4:8])
}

func (p *Header) String() string {
	return fmt.Sprintf("MessageID:%#x, ResultID:%#x", p.MessageID, p.ResultID)
}
