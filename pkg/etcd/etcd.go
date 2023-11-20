package etcd

import (
	"context"
	"encoding/json"
	"path"
	"social/pkg/bench"
	xretcd "social/pkg/lib/etcd"
	xrutil "social/pkg/lib/util"
	"time"

	"github.com/pkg/errors"
)

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
func Start(conf *bench.Etcd, BusChannel chan interface{}, onFunc xretcd.OnFunc) error {
	if len(conf.Addrs) == 0 {
		return nil
	}
	etcdValue, err := json.Marshal(conf.Value)
	if err != nil {
		return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	}
	var kvSlice []xretcd.KV
	kvSlice = append(kvSlice, xretcd.KV{
		Key:   conf.Key,
		Value: string(etcdValue),
	})
	err = xretcd.GetInstance().Start(context.TODO(),
		xretcd.NewOptions().
			SetAddrs(conf.Addrs).
			SetTTL(conf.TTL).
			SetDialTimeout(5*time.Second).
			SetKV(kvSlice).SetOnFunc(onFunc).
			SetEventChan(BusChannel),
	)
	if err != nil {
		return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	}

	// 续租
	err = xretcd.GetInstance().KeepAlive(context.TODO())
	if err != nil {
		return errors.WithMessagef(err, xrutil.GetCodeLocation(1).String())
	}

	return nil
}
