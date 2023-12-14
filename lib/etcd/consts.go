package etcd

import (
	"fmt"
	"time"
)

var (
	grantLeaseRetryDuration     = time.Second * 3 // 授权租约 重试 间隔时长
	grantLeaseMaxRetriesDefault = 600             // 授权租约 最大 重试次数 默认值
	dialTimeoutDefault          = time.Second * 5 //dialTimeout is the timeout for failing to establish a connection. 默认值
)

const (
	TtlSecondDefault    int64  = 33 //默认TTL时间 秒
	WatchMsgTypeService string = "service"
	WatchMsgTypeCommand string = "command"
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
