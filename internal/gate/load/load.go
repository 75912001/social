package load

import (
	"social/pkg/bench"
	"social/pkg/server"
)

// AvailableLoad 可用负载
func AvailableLoad() uint32 {
	if server.GetInstance().IsStopping() {
		return 0
	}
	return bench.GetInstance().Base.AvailableLoad // todo menglingchao
}
