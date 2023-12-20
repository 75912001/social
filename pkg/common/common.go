package common

import "fmt"

// GenerateServiceKey 生成服务的key
func GenerateServiceKey(zoneID uint32, serviceName string, serviceID uint32) string {
	return fmt.Sprintf("%v.%v.%v", zoneID, serviceName, serviceID)
}
