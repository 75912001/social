package gate

import (
	"context"
	"github.com/pkg/errors"
	libactor "social/lib/actor"
	liblog "social/lib/log"
)

func NewActor(id string, opt *libactor.Options) *Actor {
	actor := &Actor{
		Normal: &libactor.Normal{
			ID: id,
		},
	}
	err := actor.OnStart(context.Background(), opt)
	if err != nil {
		liblog.GetInstance().Error(err)
		return nil
	}
	return actor
}

func DelActor(ctx context.Context, actor *Actor) error {
	actor.Exit()
	err := actor.OnPreStop(ctx)
	if err != nil {
		liblog.GetInstance().Error(err)
		//不能返回,需要继续执行OnStop,清理actor资源
		//return errors.WithMessage(err, "OnPreStop")
	}
	err = actor.OnStop(ctx)
	if err != nil {
		return errors.WithMessage(err, "OnStop")
	}
	return nil
}

type Actor struct {
	*libactor.Normal
}

func (p *Actor) OnDefaultHandler(v interface{}) error {
	return nil
}

func (p *Actor) OnPreStop(_ context.Context) error {
	return nil
}
