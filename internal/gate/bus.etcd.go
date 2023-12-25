package gate

import (
	"encoding/json"
	"github.com/pkg/errors"
	libbench "social/lib/bench"
	libconsts "social/lib/consts"
	libetcd "social/lib/etcd"
	libruntime "social/lib/runtime"
	pkgcommon "social/pkg/common"
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
		// 处理etcd命令事件
	case libetcd.WatchMsgTypeService:
		if serviceIDU64, err := strconv.ParseUint(serviceID, 10, 64); err != nil {
			gate.LogMgr.Fatal(libconsts.Etcd, key, value, serviceID, err)
			return nil
		} else {
			serviceIDU32 = uint32(serviceIDU64)
		}
		switch serviceName {
		// 收到其它 服务 启动、关闭 的信息
		case pkgserver.NameLogin: //登录服务
		case pkgserver.NameGate: //网关
		case pkgserver.NameFriend: //好友
			serverKey := pkgcommon.GenerateServiceKey(zoneIDU32, serviceName, serviceIDU32)
			if 0 == len(value) { //将该服务从所在区域中移除
				gate.LogMgr.Warnf("%s delete service with key:%s, value empty", libconsts.Etcd, key)
				gate.friendMgr.Del(serverKey)
				return nil
			}
			var etcdValueJson libbench.EtcdValueJson
			if err := json.Unmarshal([]byte(value), &etcdValueJson); err != nil {
				return errors.WithMessagef(err, "%v bench EtcdValueJson json Unmarshal %v", libconsts.Etcd, value)
			}
			//查找
			friend, ok := gate.friendMgr.Find(serverKey)
			if ok { //有,更新
				friend.EtcdValueJson = etcdValueJson
				return nil
			}
			// todo menglingchao 没有,链接,添加
			//... 链接 ...

			gate.friendMgr.Add(serverKey, &Friend{
				key:           serverKey,
				Stream:        nil,
				EtcdValueJson: libbench.EtcdValueJson{},
			})
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
