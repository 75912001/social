package gate

import libutil "social/lib/util"

type FriendMgr struct {
	*libutil.Mgr[string, *Friend]
}
