package server

// AvailableLoad 可用负载
func AvailableLoad() uint32 {
	if GetInstance().IsStopping() {
		return 0
	}
	return GetInstance().BenchMgr.Base.AvailableLoad
}
