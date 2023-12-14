package handler

import (
	libconstant "social/lib/consts"
	liblog "social/lib/log"
	libutil "social/lib/util"
	pkgcommon "social/pkg/common"
	pkgetcd "social/pkg/etcd"
	pkgserver "social/pkg/server"
	"strconv"
)

// OnEventEtcd
// e.g.:/projectName/service/${zoneID}/${serviceName}/${serviceID}
// e.g.:/projectName/command/${zoneID}/${serviceName}/${serviceID}
func OnEventEtcd(key string, value string) error {
	liblog.GetInstance().Infof("%v key:%v, value:%v", libconstant.Etcd, key, value)

	var zoneIDU32 uint32
	var serviceIDU32 uint32
	msgType, zoneID, serviceName, serviceID := pkgetcd.Parse(key)

	if zoneIDU64, err := strconv.ParseUint(zoneID, 10, 32); err != nil {
		liblog.GetInstance().Errorf(libconstant.Etcd, key, value, err)
		return nil
	} else {
		zoneIDU32 = uint32(zoneIDU64)
	}

	switch msgType {
	case pkgcommon.EtcdWatchMsgTypeCommand:
		//todo menglingchao
	case pkgcommon.EtcdWatchMsgTypeService:
		if serviceIDU64, err := strconv.ParseUint(serviceID, 10, 64); err != nil {
			liblog.GetInstance().Fatal(libconstant.Etcd, key, value, serviceID, err)
			return nil
		} else {
			serviceIDU32 = uint32(serviceIDU64)
		}
		switch serviceName {
		// 收到其它 服务 启动、关闭 的信息
		case pkgserver.NameGate: //网关
		case pkgserver.NameFriend: //好友
			if 0 == len(value) {
				liblog.GetInstance().Warnf("%s delete service with key:%s, value empty", libconstant.Etcd, key)
				//todo menglingchao  将该服务从所在区域中移除
				return nil
			} else {
				//todo menglingchao
				//查找,添加,链接
			}
			liblog.GetInstance().Warnf("%s add service zone_id:%d service_id:%d with key:%s",
				libconstant.Etcd, zoneIDU32, serviceIDU32, key)
			//todo menglingchao 添加服务
		case pkgserver.NameInteraction: //交互
		case pkgserver.NameNotification: //通知
		case pkgserver.NameBlog: //博客
		case pkgserver.NameRecommendation: //推荐
		case pkgserver.NameCleansing: //清洗
		case pkgserver.NameRobot: //清洗
		default:
			liblog.GetInstance().Errorf("%v %v %v", key, value, libutil.GetCodeLocation(1))
		}
	default:
	}
	return nil
}
