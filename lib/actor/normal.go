package actor

import (
	"context"
	"github.com/pkg/errors"
	"runtime/debug"
	libbench "social/lib/bench"
	libconsts "social/lib/consts"
	liberror "social/lib/error"
	liblog "social/lib/log"
	libruntime "social/lib/runtime"
	libtime "social/lib/time"
	libutil "social/lib/util"
	"sync"
)

type State struct {
}

type Behavior struct {
}

type Normal[TKey comparable] struct {
	key      TKey
	options  *Options
	state    *State
	behavior *Behavior

	mailBox    chan IMsg
	waitGroup  sync.WaitGroup // Stop 等待信号
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
}

func (p *Normal[TKey]) GetKey() TKey {
	return p.key
}

func (p *Normal[TKey]) handler() error {
	var err error
	for {
		select {
		case <-p.cancelCtx.Done():
			return errors.WithMessagef(context.Canceled, "%v", p.GetKey())
		case v, ok := <-p.mailBox:
			if !ok {
				return errors.WithMessagef(liberror.ChannelClosed, "%v", p.GetKey())
			}
			nowTime := libtime.NowTime()
			liblog.GetInstance().Tracef("Actor %v received message: %v", p.GetKey(), v)
			switch t := v.(type) {
			case *Msg:
				err = p.options.onHandler(t.unserializedPacket)
			default:
				liblog.PrintfErr("Actor %v received message: %v", p.GetKey(), v)
			}
			if err != nil {
				liblog.PrintErr(v, err)
			}
			if libutil.IsDebug() {
				dt := libtime.NowTime().Sub(nowTime).Milliseconds()
				if dt > 50 {
					liblog.GetInstance().Warnf("cost time50: %v Millisecond with event type:%T", dt, v)
				} else if dt > 20 {
					liblog.GetInstance().Warnf("cost time20: %v Millisecond with event type:%T", dt, v)
				} else if dt > 10 {
					liblog.GetInstance().Warnf("cost time10: %v Millisecond with event type:%T", dt, v)
				}
			}
		}
	}
}

func (p *Normal[TKey]) start(ctx context.Context, opts ...*Options) error {
	p.options = merge(opts...)
	err := configure(p.options)
	if err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelCtx = ctxWithCancel
	p.cancelFunc = cancelFunc
	p.mailBox = make(chan IMsg, libbench.GetInstance().Base.ActorChannelNumber)
	p.waitGroup.Add(1)
	go func() {
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					liblog.PrintErr(libconsts.GoroutinePanic, err, debug.Stack())
				}
			}
			p.waitGroup.Done()
			liblog.PrintInfo(libconsts.GoroutineDone, p)
		}()
		err = p.handler()
		if err != nil {
			liblog.PrintErr(err)
		}
		// goroutine 退出,再设置chan为nil, (如果没有退出就设置为nil, 读chan == nil  会 block)
		p.mailBox = nil
	}()
	return nil
}

func (p *Normal[TKey]) stop(_ context.Context) error {
	liblog.GetInstance().Warnf("actor stop... %v", p.GetKey())
	close(p.mailBox)
	// 等待 goroutine退出. 阻塞等待mailBox的消息处理退出
	p.waitGroup.Wait()
	liblog.GetInstance().Warnf("actor stop done %v", p.GetKey())
	return nil
}
