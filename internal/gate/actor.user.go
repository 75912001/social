package gate

import (
	"context"
	"google.golang.org/grpc"
	libactor "social/lib/actor"
)

func NewUser(id string) *User {
	user := &User{}
	user.Normal = libactor.NewNormal(id, libactor.NewOptions().WithDefaultHandler(user.OnDefaultHandler))
	return user
}

// User 用户
type User struct {
	*libactor.Normal
	Stream grpc.ServerStream
}

func (p *User) Exit(ctx context.Context) error {
	//todo menglingchao 退出之前的操作
	return p.Normal.Exit(ctx)
}

func (p *User) OnDefaultHandler(v interface{}) error {
	return nil
}
