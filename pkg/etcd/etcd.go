package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"social/lib/etcd"
	libutil "social/lib/util"
	pkgbench "social/pkg/bench"
	"social/pkg/consts"
	"time"

	"github.com/pkg/errors"
)

// GenerateServiceKey 生成服务注册的key
func GenerateServiceKey(zoneID uint32, serviceName string, serviceID uint32) string {
	return fmt.Sprintf("%v/%v/%v/%v/%v",
		consts.ProjectName, libetcd.WatchMsgTypeService,
		zoneID, serviceName, serviceID)
}

// EtcdGenerateWatchServicePrefix 生成关注服务的前缀
func EtcdGenerateWatchServicePrefix() string {
	return fmt.Sprintf("%v/%v/",
		consts.ProjectName, libetcd.WatchMsgTypeService)
}

// EtcdGenerateWatchCommandPrefix 生成关注命令的前缀
func EtcdGenerateWatchCommandPrefix(zoneID uint32, serviceName string) string {
	return fmt.Sprintf("%v/%v/%v/%v/",
		consts.ProjectName, libetcd.WatchMsgTypeCommand,
		zoneID, serviceName)
}

// Parse
// e.g.:/objectName/service/${zoneID}/${serviceName}/${serviceID}
// e.g.:/objectName/command/${zoneID}/${serviceName}/${serviceID}
func Parse(key string) (msgType string, zoneID string, serviceName string, serviceID string) {
	serviceID = path.Base(key)

	key = path.Dir(key)
	serviceName = path.Base(key)

	key = path.Dir(key)
	zoneID = path.Base(key)

	key = path.Dir(key)
	msgType = path.Base(key)
	return msgType, zoneID, serviceName, serviceID
}

// Start 启动Etcd
func Start(conf *pkgbench.Etcd, BusChannel chan interface{}, onFunc etcd.OnFunc) error {
	etcdValue, err := json.Marshal(conf.Value)
	if err != nil {
		return errors.WithMessagef(err, libutil.GetCodeLocation(1).String())
	}
	var kvSlice []etcd.KV
	kvSlice = append(kvSlice, etcd.KV{
		Key:   conf.Key,
		Value: string(etcdValue),
	})
	err = etcd.GetInstance().Start(context.TODO(),
		etcd.NewOptions().
			WithAddrs(conf.Addrs).
			WithTTL(conf.TTL).
			WithDialTimeout(5*time.Second).
			WithKV(kvSlice).WithOnFunc(onFunc).
			WithOutgoingEventChan(BusChannel),
	)
	if err != nil {
		return errors.WithMessagef(err, libutil.GetCodeLocation(1).String())
	}

	// 续租
	err = etcd.GetInstance().Run(context.TODO())
	if err != nil {
		return errors.WithMessagef(err, libutil.GetCodeLocation(1).String())
	}

	return nil
}
