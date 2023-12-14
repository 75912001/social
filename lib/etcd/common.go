package etcd

import (
	"fmt"
	"path"
)

// GenerateServiceKey 生成服务注册的key
func GenerateServiceKey(projectName string, zoneID uint32, serviceName string, serviceID uint32) string {
	return fmt.Sprintf("%v/%v/%v/%v/%v",
		projectName, WatchMsgTypeService,
		zoneID, serviceName, serviceID)
}

// GenerateWatchServicePrefix 生成关注服务的前缀
func GenerateWatchServicePrefix(projectName string) string {
	return fmt.Sprintf("%v/%v/",
		projectName, WatchMsgTypeService)
}

// GenerateWatchCommandPrefix 生成关注命令的前缀
func GenerateWatchCommandPrefix(projectName string, zoneID uint32, serviceName string) string {
	return fmt.Sprintf("%v/%v/%v/%v/",
		projectName, WatchMsgTypeCommand,
		zoneID, serviceName)
}

// Parse
// e.g.:${objectName}/service/${zoneID}/${serviceName}/${serviceID}
// e.g.:${objectName}/command/${zoneID}/${serviceName}/${serviceID}
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
