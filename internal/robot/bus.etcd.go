package robot

import (
	libconsts "social/lib/consts"
	libetcd "social/lib/etcd"
	pkgserver "social/pkg/server"
	"strconv"
)

// OnEventEtcd 处理事件-etcd
func (p *Bus) OnEventEtcd(key string, value string) error {
	robot.LogMgr.Infof("%v key:%v, value:%v", libconsts.Etcd, key, value)

	var zoneIDU32 uint32
	var serviceIDU32 uint32
	msgType, zoneID, serviceName, serviceID := libetcd.Parse(key)

	if zoneIDU64, err := strconv.ParseUint(zoneID, 10, 32); err != nil {
		robot.LogMgr.Errorf(libconsts.Etcd, key, value, err)
		return nil
	} else {
		zoneIDU32 = uint32(zoneIDU64)
	}

	switch msgType {
	case libetcd.WatchMsgTypeCommand:
		//todo menglingchao 处理etcd命令事件
	case libetcd.WatchMsgTypeService:
		if serviceIDU64, err := strconv.ParseUint(serviceID, 10, 64); err != nil {
			robot.LogMgr.Fatal(libconsts.Etcd, key, value, serviceID, err)
			return nil
		} else {
			serviceIDU32 = uint32(serviceIDU64)
		}
		switch serviceName {
		// 收到其它 服务 启动、关闭 的信息
		case pkgserver.NameGate: //网关
			// todo menglingchao 筛选,链接
			robot.LogMgr.Info(key, value, zoneIDU32, serviceName, serviceIDU32)
		default:
		}
	default:
	}
	return nil
}
