package gate

import (
	"google.golang.org/grpc"
)

// User 用户
type User struct {
	key    string
	Stream grpc.ServerStream
}

// OnHandler 处理收到的消息
func (p *User) OnHandler(v interface{}) error {
	return nil
}
