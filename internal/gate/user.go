package gate

import (
	"fmt"
	"google.golang.org/grpc"
	libactor "social/lib/actor"
)

// User 用户
type User struct {
	key    string
	Stream grpc.ServerStream
}

func (p *User) String() string {
	return fmt.Sprintf("%v", p.key)
}

// OnMailBox 处理收到的消息
func (p *User) OnMailBox(msg *libactor.Msg) error {
	return nil
}
