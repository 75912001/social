package actor

import (
	"context"
	liblog "social/lib/log"
	"sync"
)

// NewMgr 创建 Mgr 实例
func NewMgr[TKey comparable]() *Mgr[TKey] {
	return &Mgr[TKey]{
		actorMap: make(map[TKey]*Normal[TKey]),
	}
}

// Mgr 管理 IActor 的简单系统
type Mgr[TKey comparable] struct {
	actorMap map[TKey]*Normal[TKey]
	lock     sync.RWMutex
}

// SpawnActor 创建一个新的 IActor 并添加到系统中
func (p *Mgr[TKey]) SpawnActor(ctx context.Context, key TKey, opt *Options) {
	p.lock.Lock()
	defer p.lock.Unlock()

	normal := &Normal[TKey]{
		key:      key,
		options:  opt,
		state:    &State{},
		behavior: &Behavior{},
	}
	p.actorMap[key] = normal

	err := normal.start(ctx, opt)
	if err != nil {
		liblog.PrintfErr("IActor %v start failed err:%v", key, err)
	}
}

// DeleteActor 从系统中删除指定的 IActor
func (p *Mgr[TKey]) DeleteActor(ctx context.Context, key TKey) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if actor, ok := p.actorMap[key]; ok {
		if actor.cancelFunc != nil {
			actor.cancelFunc()
			actor.cancelFunc = nil
		}
		_ = actor.stop(ctx)
		delete(p.actorMap, key)
	}
}

// SendMessage 向指定 IActor 发送消息
func (p *Mgr[TKey]) SendMessage(key TKey, msg IMsg) {
	p.lock.Lock()
	defer p.lock.Unlock()

	actor, ok := p.actorMap[key]
	if ok {
		actor.mailBox <- msg
	} else {
		liblog.PrintfErr("IActor %v not found", key)
	}
}
