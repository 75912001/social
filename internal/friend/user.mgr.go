package friend

import (
	"context"
	"google.golang.org/grpc"
	libactor "social/lib/actor"
	"sync"
)

func (p *UserMgr) SpawnUser(key string, stream grpc.ServerStream) *User {
	p.lock.Lock()
	defer p.lock.Unlock()
	user := &User{
		key:    key,
		Stream: stream,
	}
	p.actorMgr.SpawnActor(context.Background(), key, libactor.NewOptions().WithDefaultHandler(user.OnHandler))
	p.userMap[key] = user
	return user
}

func (p *UserMgr) DeleteUser(key string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.actorMgr.DeleteActor(context.Background(), key)
	delete(p.userMap, key)
}

func (p *UserMgr) Find(key string) *User {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.userMap[key]
}

func (p *UserMgr) FindByStream(stream grpc.ServerStream) *User {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, user := range p.userMap {
		if user.Stream == stream {
			return user
		}
	}
	return nil
}

type UserMgr struct {
	actorMgr *libactor.Mgr[string] // e.g.:1.lp.1
	userMap  map[string]*User      // key:id, val:user
	lock     sync.RWMutex
}
