package handler

import (
	libconstant "social/lib/consts"
	liblog "social/lib/log"
)

// OnEventEtcd
// e.g.:/projectName/service/${zoneID}/${serviceName}/${serviceID}
// e.g.:/projectName/command/${zoneID}/${serviceName}/${serviceID}
func OnEventEtcd(key string, value string) error {
	liblog.GetInstance().Infof("%v key:%v, value:%v", libconstant.Etcd, key, value)
	return nil
}
