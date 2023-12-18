package gate

import (
	libconsts "social/lib/consts"
	libetcd "social/lib/etcd"
	libruntime "social/lib/runtime"
	pkgserver "social/pkg/server"
	"strconv"
)

// OnEventEtcd 处理事件-etcd
func (p *Bus) OnEventEtcd(key string, value string) error {
	gate.LogMgr.Infof("%v key:%v, value:%v", libconsts.Etcd, key, value)

	var zoneIDU32 uint32
	var serviceIDU32 uint32
	msgType, zoneID, serviceName, serviceID := libetcd.Parse(key)

	if zoneIDU64, err := strconv.ParseUint(zoneID, 10, 32); err != nil {
		gate.LogMgr.Errorf(libconsts.Etcd, key, value, err)
		return nil
	} else {
		zoneIDU32 = uint32(zoneIDU64)
	}

	switch msgType {
	case libetcd.WatchMsgTypeCommand:
		//todo menglingchao 处理etcd命令事件
	case libetcd.WatchMsgTypeService:
		if serviceIDU64, err := strconv.ParseUint(serviceID, 10, 64); err != nil {
			gate.LogMgr.Fatal(libconsts.Etcd, key, value, serviceID, err)
			return nil
		} else {
			serviceIDU32 = uint32(serviceIDU64)
		}
		switch serviceName {
		// 收到其它 服务 启动、关闭 的信息
		case pkgserver.NameGate: //网关
		case pkgserver.NameFriend: //好友
			if 0 == len(value) {
				gate.LogMgr.Warnf("%s delete service with key:%s, value empty", libconsts.Etcd, key)
				//todo menglingchao  将该服务从所在区域中移除
				return nil
			}
			gate.LogMgr.Warnf("%s service zone_id:%d service_id:%d with key:%s",
				libconsts.Etcd, zoneIDU32, serviceIDU32, key)
			//todo menglingchao 查找,添加,链接
			//有,更新
			//没有,添加,链接
		case pkgserver.NameInteraction: //交互
		case pkgserver.NameNotification: //通知
		case pkgserver.NameBlog: //博客
		case pkgserver.NameRecommendation: //推荐
		case pkgserver.NameCleansing: //清洗
		case pkgserver.NameRobot: //机器人
		default:
			gate.LogMgr.Errorf("%v %v %v", key, value, libruntime.Location())
		}
	default:
	}
	return nil
}
