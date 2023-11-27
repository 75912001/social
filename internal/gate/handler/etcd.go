package handler

import (
	"social/pkg"
	"social/pkg/etcd"
	xrconstant "social/pkg/lib/constant"
	xrlog "social/pkg/lib/log"
	xrutil "social/pkg/lib/util"
	"social/pkg/server"
	"strconv"
)

// OnEventEtcd
// e.g.:/projectName/service/${zoneID}/${serviceName}/${serviceID}
// e.g.:/projectName/command/${zoneID}/${serviceName}/${serviceID}
func OnEventEtcd(key string, value string) error {
	xrlog.GetInstance().Infof("%v key:%v, value:%v", xrconstant.Etcd, key, value)

	var zoneIDU32 uint32
	var serviceIDU32 uint32
	msgType, zoneID, serviceName, serviceID := etcd.Parse(key)

	if zoneIDU64, err := strconv.ParseUint(zoneID, 10, 32); err != nil {
		xrlog.GetInstance().Fatal(xrconstant.Etcd, key, value, err)
		return nil
	} else {
		zoneIDU32 = uint32(zoneIDU64)
	}

	switch msgType {
	case pkg.EtcdWatchMsgTypeCommand:
		//todo menglingchao
	case pkg.EtcdWatchMsgTypeService:
		if serviceIDU64, err := strconv.ParseUint(serviceID, 10, 64); err != nil {
			xrlog.GetInstance().Fatal(xrconstant.Etcd, key, value, serviceID, err)
			return nil
		} else {
			serviceIDU32 = uint32(serviceIDU64)
		}
		switch serviceName {
		// 收到其它 服务 启动、关闭 的信息
		case server.NameGate: //网关
		case server.NameFriend: //好友
			if 0 == len(value) {
				xrlog.GetInstance().Warnf("%s delete service with key:%s, value empty", xrconstant.Etcd, key)
				//todo menglingchao  将该服务从所在区域中移除
				return nil
			}
			xrlog.GetInstance().Warnf("%s add service zone_id:%d service_id:%d with key:%s",
				xrconstant.Etcd, zoneIDU32, serviceIDU32, key)
			//todo menglingchao 添加服务
		case server.NameInteraction: //交互
		case server.NameNotification: //通知
		case server.NameBlog: //博客
		case server.NameRecommendation: //推荐
		case server.NameCleansing: //清洗
		case server.NameRobot: //清洗
		default:
			xrlog.GetInstance().Errorf("%v %v %v", key, value, xrutil.GetCodeLocation(1))
		}
	default:
	}
	return nil
}
