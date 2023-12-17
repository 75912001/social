package gate

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	libactor "social/lib/actor"
	liblog "social/lib/log"
)

func NewUser(id string, opt *libactor.Options) *User {
	user := &User{}
	actor := &Actor{}
	actor.Normal = NewActor(id, opt).Normal
	user.
		err := actor.OnStart(context.Background(), opt)
	if err != nil {
		liblog.GetInstance().Error(err)
		return nil
	}
	return actor
}

func (p *User) DelUser(ctx context.Context) error {
	DelActor(ctx, p.Actor)
	p.OnPreStop()
	return nil
}

type User struct {
	*Actor
	Stream grpc.ServerStream
}

func (p *User) OnDefaultHandler(v interface{}) error {
	return nil
}

func (p *User) OnPreStop(_ context.Context) error {
	return nil
}
