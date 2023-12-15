package actor

import (
	"context"
	"fmt"
	"time"
)

// IActor 是 Actor 接口
type IActor interface {
	Send(msg IMsg) error                 //发送消息给Actor邮箱中
	OnStart(ctx context.Context) error   //启动
	OnRun(ctx context.Context) error     //运行
	OnPreStop(ctx context.Context) error //停止前的处理
	OnStop(ctx context.Context) error    //停止

}

// Actor 是 IActor 接口的实现
type Actor struct {
	ID      string
	ctx     context.Context
	cancel  context.CancelFunc
	mailbox chan IMsg
}

// NewActor 创建一个新的 Actor
func NewActor(id string) *Actor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Actor{
		ID:      id,
		ctx:     ctx,
		cancel:  cancel,
		mailbox: make(chan IMsg, 10), // 缓冲区大小为 10
	}
}

// Send 实现 IActor 接口的 Send 方法
func (a *Actor) Send(msg IMsg) error {
	select {
	case a.mailbox <- msg:
		return nil
	case <-a.ctx.Done():
		return context.Canceled
	}
}

// Start 实现 IActor 接口的 Start 方法
func (a *Actor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-a.mailbox:
				// 在这里处理接收到的消息 todo menglingchao []处理Actor收到的消息
				// 使用 go xxx() 的方式避免阻塞... 这里需要标记消息是需要顺序处理的,还是可以多协程处理,来做不同的处理策略.
				fmt.Printf("Actor %s received message: %s\n", a.ID, msg.String())
			}
		}
	}()
}

// Stop 实现 IActor 接口的 Stop 方法
func (a *Actor) Stop() {
	a.cancel()
}

func main() {
	// 创建一个 Actor
	actor := NewActor("1")

	// 启动 Actor
	actor.Start(context.Background())

	// 发送消息给 Actor
	message := &Msg{Content: "Hello, Actor!"}
	actor.Send(message)

	// 等待一段时间，以便 Actor 处理消息
	<-time.After(2 * time.Second)

	// 停止 Actor
	actor.Stop()
}
