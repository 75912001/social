package common

import (
	"fmt"
	"github.com/pkg/errors"
	libetcd "social/lib/etcd"
	"strconv"
)

// GenerateServiceKey 生成服务的key
func GenerateServiceKey(zoneID uint32, serviceName string, serviceID uint32) string {
	return fmt.Sprintf("%v.%v.%v", zoneID, serviceName, serviceID)
}

// ParseEtcdKey 解析etcd同步过来的key
func ParseEtcdKey(key string) (msgType string, zoneID uint32, serviceName string, serviceID uint32, err error) {
	msgType, strZoneID, serviceName, strServiceID := libetcd.Parse(key)
	if zoneIDU64, err := strconv.ParseUint(strZoneID, 10, 32); err != nil {
		return "", 0, "", 0, errors.WithStack(err)
	} else {
		zoneID = uint32(zoneIDU64)
	}
	if serviceIDU64, err := strconv.ParseUint(strServiceID, 10, 32); err != nil {
		return "", 0, "", 0, errors.WithStack(err)
	} else {
		serviceID = uint32(serviceIDU64)
	}
	return msgType, zoneID, serviceName, serviceID, nil
}
