package load

import (
	pkgbench "social/lib/bench"
	pkgserver "social/pkg/server"
)

// AvailableLoad 可用负载
func AvailableLoad() uint32 {
	if pkgserver.GetInstance().IsStopping() {
		return 0
	}
	return pkgbench.GetInstance().Base.AvailableLoad
}
