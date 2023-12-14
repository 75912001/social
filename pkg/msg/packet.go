package msg

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	libutil "social/lib/util"
)

type Packet struct {
	Header  Header
	Message proto.Message // Marshal 时候使用
}

// Marshal 序列化 数据包
func (p *Packet) Marshal() ([]byte, error) {
	headerBuf := p.Header.Pack()

	var messageBuf []byte
	var err error
	if messageBuf, err = proto.Marshal(p.Message); nil != err && proto.ErrNil != err {
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
func (p *Packet) Unmarshal(data []byte) (err error) {
	p.Header.Unpack(data)
	//err = proto.Unmarshal(data[GProtoHeadLength:], p.Message)

	return nil
}
