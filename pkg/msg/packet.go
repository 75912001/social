package msg

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	libruntime "social/lib/runtime"
)

// Packet 协议包
type Packet struct {
	Header  Header
	Message proto.Message // Marshal 时候使用
}

// Marshal 序列化 数据包
func (p *Packet) Marshal() ([]byte, error) {
	headerBuf := p.Header.Marshal()

	var messageBuf []byte
	var err error
	if messageBuf, err = proto.Marshal(p.Message); nil != err && proto.ErrNil != err {
		return nil, errors.WithMessagef(err, libruntime.GetCodeLocation(1).String())
	}

	packetLength := GProtoHeadLength + uint32(len(messageBuf))
	packetBuf := make([]byte, packetLength)

	copy(packetBuf[:GProtoHeadLength], headerBuf)
	copy(packetBuf[GProtoHeadLength:], messageBuf)
	return packetBuf, nil
}

// Unmarshal 反序列化
// data:完整包数据 包头+包体
func (p *Packet) Unmarshal(data []byte) (err error) {
	p.Header.Unmarshal(data)
	return nil
}
