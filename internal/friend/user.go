package friend

import (
	"google.golang.org/grpc"
	libactor "social/lib/actor"
)

// User 用户
type User struct {
	key    string
	Stream grpc.ServerStream
}

// OnHandler 处理收到的消息
func (p *User) OnHandler(msg *libactor.Msg) error {
	return nil
}
