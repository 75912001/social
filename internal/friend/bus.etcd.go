package friend

import (
	"github.com/pkg/errors"
	libconsts "social/lib/consts"
	libetcd "social/lib/etcd"
	libruntime "social/lib/runtime"
	pkgcommon "social/pkg/common"
	pkgserver "social/pkg/server"
)

// OnEventEtcd 处理事件-etcd
func (p *Bus) OnEventEtcd(key string, value string) error {
	app.LogMgr.Infof("%v key:%v, value:%v", libconsts.Etcd, key, value)
	msgType, zoneID, serviceName, serviceID, err := pkgcommon.ParseEtcdKey(key)
	if err != nil {
		return errors.WithStack(err)
	}
	app.LogMgr.Info("msgType:%v, zoneID:%v, serviceName:%v, serviceID:%v", msgType, zoneID, serviceName, serviceID)
	switch msgType {
	case libetcd.WatchMsgTypeCommand:
		// 处理etcd命令事件
	case libetcd.WatchMsgTypeService:
		switch serviceName {
		// 收到其它 服务 启动、关闭 的信息
		case pkgserver.NameLogin: //登录服务
		case pkgserver.NameGate: //网关
		case pkgserver.NameFriend: //好友
		case pkgserver.NameInteraction: //交互
		case pkgserver.NameBlog: //博客
		case pkgserver.NameRecommendation: //推荐
		case pkgserver.NameRobot: //机器人
		case pkgserver.NameNotification: //通知
		// todo menglingchao friend 链接 notification ...
		//	... 链接
		//	... 添加
		//	... 移除
		case pkgserver.NameCleansing: //清洗
		// todo menglingchao friend 链接 cleansing ...
		//	... 链接
		//	... 添加
		//	... 移除
		default:
			app.LogMgr.Errorf("%v %v %v", key, value, libruntime.Location())
		}
	default:
	}
	return nil
}
