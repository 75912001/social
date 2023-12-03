package common

import "fmt"

const EtcdWatchMsgTypeService string = "service"
const EtcdWatchMsgTypeCommand string = "command"
const EtcdTtlSecondDefault int64 = 33 //默认TTL时间 秒

// EtcdGenerateServiceKey 生成服务注册的key
func EtcdGenerateServiceKey(zoneID uint32, serviceName string, serviceID uint32) string {
	return fmt.Sprintf("%v/%v/%v/%v/%v",
		ProjectName, EtcdWatchMsgTypeService,
		zoneID, serviceName, serviceID)
}

// EtcdGenerateWatchServicePrefix 生成关注服务的前缀
func EtcdGenerateWatchServicePrefix() string {
	return fmt.Sprintf("%v/%v/",
		ProjectName, EtcdWatchMsgTypeService)
}

// EtcdGenerateWatchCommandPrefix 生成关注命令的前缀
func EtcdGenerateWatchCommandPrefix(zoneID uint32, serviceName string) string {
	return fmt.Sprintf("%v/%v/%v/%v/",
		ProjectName, EtcdWatchMsgTypeCommand,
		zoneID, serviceName)
}
