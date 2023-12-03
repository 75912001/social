package etcd

import (
	"context"
	"encoding/json"
	"path"
	pkgbench "social/pkg/bench"
	libetcd "social/pkg/lib/etcd"
	libutil "social/pkg/lib/util"
	"time"

	"github.com/pkg/errors"
)

const WatchMsgTypeService string = "service"
const WatchMsgTypeCommand string = "command"

const TtlSecondDefault int64 = 33 //默认TTL时间 秒

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
func Start(conf *pkgbench.Etcd, BusChannel chan interface{}, onFunc libetcd.OnFunc) error {
	etcdValue, err := json.Marshal(conf.Value)
	if err != nil {
		return errors.WithMessagef(err, libutil.GetCodeLocation(1).String())
	}
	var kvSlice []libetcd.KV
	kvSlice = append(kvSlice, libetcd.KV{
		Key:   conf.Key,
		Value: string(etcdValue),
	})
	err = libetcd.GetInstance().Start(context.TODO(),
		libetcd.NewOptions().
			SetAddrs(conf.Addrs).
			SetTTL(conf.TTL).
			SetDialTimeout(5*time.Second).
			SetKV(kvSlice).SetOnFunc(onFunc).
			SetOutgoingEventChan(BusChannel),
	)
	if err != nil {
		return errors.WithMessagef(err, libutil.GetCodeLocation(1).String())
	}

	// 续租
	err = libetcd.GetInstance().KeepAlive(context.TODO())
	if err != nil {
		return errors.WithMessagef(err, libutil.GetCodeLocation(1).String())
	}

	return nil
}
