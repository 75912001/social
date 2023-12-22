package friend

import (
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	pkgmsg "social/pkg/msg"
	pkgproto "social/pkg/proto"
	protofriend "social/pkg/proto/friend"
)

// CanHandle 是否可以处理的命令
// 判断cmd是否符合自己可以处理的消息类型
// 如果符合，返回true；否则，返回false
func (p *Friend) CanHandle(messageID uint32) bool {
	if protofriend.MessageMinCMD < messageID && messageID < protofriend.MessageMaxCMD { //gate的消息
		return true
	}
	return false
}

func (p *Friend) Handle(stream protofriend.Service_BidirectionalBinaryDataServer, header *pkgmsg.Header, body []byte) error {
	p.LogMgr.Trace(pkgproto.CMDMap[header.MessageID], stream, header, body)
	switch header.MessageID {
	case protofriend.UpdateFriendMaxReqCMD: // uint32 = 0x2010         // 修改好友数量最大值
	case protofriend.GetFriendListReqCMD: // uint32 = 0x2020           // 获取好友列表请求
	case protofriend.ApplyFriendReqCMD: // uint32 = 0x2030             // 申请成为好友请求
	case protofriend.AgreeApplyFriendReqCMD: // uint32 = 0x2040        // 接受申请好友
	case protofriend.RejectApplyFriendReqCMD: // uint32 = 0x2050       // 拒绝申请好友请求
	case protofriend.RemoveFriendReqCMD: // uint32 = 0x2060            // 移除好友请求
	case protofriend.UpdateFriendRemarkReqCMD: // uint32 = 0x2070      // 修改好友备注请求
	case protofriend.UpdateFriendRelationReqCMD: // uint32 = 0x2080    // 改变好友的关系值请求
	case protofriend.AddUserToBlackListReqCMD: // uint32 = 0x2090      // 将用户加入黑名单请求
	case protofriend.RemoveUserFromBlackListReqCMD: // uint32 = 0x20a0 // 将用户从黑名单中移除请求
	case protofriend.GetUserStatusReqCMD: // uint32 = 0x20b0           // 获取用户状态请求
	case protofriend.UpdateUserStatusReqCMD: // uint32 = 0x20c0        // 改变状态请求
	case protofriend.UpdateUserLocationReqCMD: // uint32 = 0x20d0      // 改变经纬值请求
	default:
		p.LogMgr.Error(liberror.MessageIDNonExistent, libruntime.Location())
		return liberror.MessageIDNonExistent
	}
	return nil
}
