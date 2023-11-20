package load

import "social/pkg/common"

// AvailableLoad 可用负载
func AvailableLoad() uint32 {
	return common.BenchJsonAvailableLoadMaxDefault // todo menglingchao
}
