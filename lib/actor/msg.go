package actor

import (
	libpb "social/lib/pb"
)

// IMsg 是消息接口
type IMsg interface {
}

// Msg 是 IMsg 的实现
type Msg struct {
	unserializedPacket libpb.UnserializedPacket
	obj                interface{}
}
