package gate

import (
	"context"
	"google.golang.org/grpc"
	libactor "social/lib/actor"
	libutil "social/lib/util"
	"sync"
)

func (p *FriendMgr) SpawnFriend(key string, stream grpc.ServerStream) *Friend {
	friend := &Friend{}
	go friend.Service.OnBidirectionalRecv()

	return nil
}

type FriendMgr struct {
	*libutil.Mgr[string, *Friend]
	actorMgr *libactor.Mgr[string] // e.g.:1.lp.1
	lock     sync.RWMutex
}

// 获取一个可用的服务
func (p *FriendMgr) getAvailable() *Friend {
	//todo
	return nil
}

// Send 发送
func (p *FriendMgr) Send(ctx context.Context) error {
	//todo
	//var friend * Friend
	//friend.Stream.
	return nil
}
