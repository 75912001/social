package gate

import (
	"context"
	libutil "social/lib/util"
)

type FriendMgr struct {
	*libutil.Mgr[string, *Friend]
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
