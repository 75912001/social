package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	libbench "social/lib/bench"
	libconsts "social/lib/consts"
	libetcd "social/lib/etcd"
	libruntime "social/lib/runtime"
	pkgcommon "social/pkg/common"
	protofriend "social/pkg/proto/friend"
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
			serverKey := pkgcommon.GenerateServiceKey(zoneID, serviceName, serviceID)
			if 0 == len(value) { //将该服务从所在区域中移除
				app.LogMgr.Warnf("%s delete service with key:%s, value empty", libconsts.Etcd, key)
				app.friendMgr.Del(serverKey)
				return nil
			}
			var etcdValueJson libbench.EtcdValueJson
			if err := json.Unmarshal([]byte(value), &etcdValueJson); err != nil {
				return errors.WithMessagef(err, "%v bench EtcdValueJson json Unmarshal %v", libconsts.Etcd, value)
			}
			//查找
			friend, ok := app.friendMgr.Find(serverKey)
			if ok { //有,更新
				friend.EtcdValueJson = etcdValueJson
				return nil
			}
			//没有,链接,添加
			stream, err := func(ctx context.Context, addr string) (protofriend.FriendService_BidirectionalBinaryDataClient, error) {
				// 连接 gRPC 服务端
				conn, err := grpc.Dial(addr, grpc.WithInsecure())
				if err != nil {
					app.LogMgr.Errorf("%v bench EtcdValueJson grpc.Dial %v", libconsts.Etcd, err)
					return nil, err
				}
				// 创建 gRPC 客户端
				client := protofriend.NewFriendServiceClient(conn)

				// 调用双向流 RPC 方法
				stream, err := client.BidirectionalBinaryData(ctx)
				if err != nil {
					app.LogMgr.Errorf("%v bench EtcdValueJson client.BidirectionalBinaryData %v", libconsts.Etcd, err)
					return nil, err
				}
				for {
					response, err := stream.Recv()
					if err != nil {

					}
				}
				return stream, nil
			}(context.Background(), fmt.Sprintf("%v:%v", etcdValueJson.ServiceNetTCP.IP, etcdValueJson.ServiceNetTCP.Port))
			if err != nil {
				app.LogMgr.Error(err, libruntime.Location())
				return errors.WithStack(err)
			}
			app.friendMgr.Add(serverKey, &Friend{
				key:           serverKey,
				Stream:        stream,
				EtcdValueJson: libbench.EtcdValueJson{},
			})
		case pkgserver.NameInteraction: //交互
		case pkgserver.NameNotification: //通知
		case pkgserver.NameBlog: //博客
		case pkgserver.NameRecommendation: //推荐
		case pkgserver.NameCleansing: //清洗
		case pkgserver.NameRobot: //机器人
		default:
			app.LogMgr.Errorf("%v %v %v", key, value, libruntime.Location())
		}
	default:
	}
	return nil
}
