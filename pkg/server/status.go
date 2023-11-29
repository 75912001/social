package server

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync/atomic"
)

type Status uint32 //服务状态

const (
	StatusRunning  Status = 0 // 运行中
	StatusStopping Status = 1 // 关闭中
)

// IsStopping 服务是否关闭中
func (p *mgr) IsStopping() bool {
	return Status(atomic.LoadUint32((*uint32)(&p.status))) == StatusStopping
}

// IsRunning 服务是否运行中
func (p *mgr) IsRunning() bool {
	return Status(atomic.LoadUint32((*uint32)(&p.status))) == StatusRunning
}

// SetStopping 设置为关闭中
func (p *mgr) SetStopping() {
	atomic.StoreUint32((*uint32)(&p.status), uint32(StatusStopping))
}

// ServiceInfo 服务信息
//
//	NOTE 有性能影响.
//	建议 只在调试/测试时使用.
func (p *mgr) ServiceInfo() string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	return fmt.Sprintf("[goroutines-number:%v gc:%d last GC at:%v PauseTotal:%v MemStats:%+v]",
		runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal, memStats)
}
