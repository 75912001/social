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

type Normal struct {
	ID string

	options *Options
	//https://zhuanlan.zhihu.com/p/427806717
	//状态 state
	// todo menglingchao
	//指actor本身的属性信息，state只能被actor自己操作，不能被其他actor共享和操作，有效的避免加锁和数据竞争

	//行为 behavior
	// todo menglingchao
	//指actor处理逻辑，如果通过行为来操作自身state
	mailBox    chan IMsg
	waitGroup  sync.WaitGroup // Stop 等待信号
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
}

func (p *Normal) SendToMailBox(msg IMsg) error {
	select {
	case <-p.cancelCtx.Done():
		return context.Canceled
	case p.mailBox <- msg:
		return nil
	}
}

func (p *Normal) OnStart(ctx context.Context, opts ...*Options) error {
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
	go func(ctx context.Context) {
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					liblog.PrintErr(libconsts.GoroutinePanic, err, debug.Stack())
				}
			}
			p.waitGroup.Done()
			liblog.PrintInfo(libconsts.GoroutineDone, p)
		}()
		for v := range p.mailBox {
			nowTime := libtime.NowTime()
			liblog.GetInstance().Tracef("Actor %v received message: %v", p.ID, v)
			switch t := v.(type) {
			case *Msg:
				// 使用 go xxx() 的方式避免阻塞... 这里需要标记消息是需要顺序处理的,还是可以多协程处理,来做不同的处理策略.
				err = p.options.defaultHandler(t.unserializedPacket)
			default:
				liblog.GetInstance().Errorf("Actor %v received message: %v", p.ID, v)
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
		// goroutine 退出,再设置chan为nil, (如果没有退出就设置为nil, 读chan == nil  会 block)
		p.mailBox = nil
	}(ctxWithCancel)
	return nil
}

func (p *Normal) OnRun(_ context.Context) error {
	return nil
}

// Exit 退出服务
func (p *Normal) Exit() {
	liblog.GetInstance().Warnf("actor exit... %v", p.ID)
	if p.cancelFunc != nil {
		p.cancelFunc()
		p.cancelFunc = nil
	}
	liblog.GetInstance().Warnf("actor exit done %v", p.ID)
}

func (p *Normal) OnPreStop(_ context.Context) error {
	//todo menglingchao 停止前的处理
	return liberror.NotImplemented
}

func (p *Normal) OnStop(_ context.Context) error {
	liblog.GetInstance().Warnf("actor OnStop... %v", p.ID)
	close(p.mailBox)
	// 等待 goroutine退出.
	p.waitGroup.Wait()
	liblog.GetInstance().Warnf("actor OnStop done %v", p.ID)
	return nil
}
