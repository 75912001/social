package robot

import (
	libconsts "social/lib/consts"
	liblog "social/lib/log"
)

// OnEventEtcd etcd 处理函数
func OnEventEtcd(key string, value string) error {
	liblog.GetInstance().Infof("%v key:%v, value:%v", libconsts.Etcd, key, value)
	return nil
}
