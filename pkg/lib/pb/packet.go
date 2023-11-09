package pb

import (
	"github.com/golang/protobuf/proto"
)

// UnserializedPacket 未序列化的数据包
type UnserializedPacket struct {
	Header  IHeader       // 包头
	Message proto.Message // 数据
}

// IPacket 接口-数据包
type IPacket interface {
	// Marshal 序列化
	//	返回:
	//		数据
	Marshal(unserializedPacket *UnserializedPacket) (data []byte, err error)
	// Unmarshal 反序列化
	//	参数:
	//		data:数据(包头+包体)
	//	返回:
	//		header:包头
	//		message:消息
	Unmarshal(data []byte) (header IHeader, message proto.Message, err error)
}
