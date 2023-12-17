package server

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync/atomic"
)

type status uint32 //服务状态

const (
	StatusRunning  status = 0 // 运行中
	StatusStopping status = 1 // 关闭中
)

// IsStopping 服务是否关闭中
func (p *Normal) IsStopping() bool {
	return status(atomic.LoadUint32((*uint32)(&p.status))) == StatusStopping
}

// IsRunning 服务是否运行中
func (p *Normal) IsRunning() bool {
	return status(atomic.LoadUint32((*uint32)(&p.status))) == StatusRunning
}

// SetStopping 设置为关闭中
func (p *Normal) SetStopping() {
	atomic.StoreUint32((*uint32)(&p.status), uint32(StatusStopping))
}

// Info 服务信息
// [NOTE] 有性能影响.
// 建议 只在调试/测试时使用.
func (p *Normal) Info() string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	return fmt.Sprintf("[goroutines-number:%v gc:%d last GC at:%v PauseTotal:%v MemStats:%+v]",
		runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal, memStats)
}

func (p *Normal) serviceInformationPrintingStart() {
	p.informationPrintingTimerSecond = p.TimerMgr.AddSecond(p.serviceInformationPrinting, nil, p.TimeMgr.TimeSecond()+InfoTimeOutSec)
}

// 服务信息 打印
func (p *Normal) serviceInformationPrinting(_ interface{}) {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	p.LogMgr.Infof("goroutineCnt:%d, busChannel:%d, numGC:%d, lastGC:%v, GCPauseTotal:%v",
		runtime.NumGoroutine(), len(p.busChannel), s.NumGC, s.LastGC, s.PauseTotal)
	p.serviceInformationPrintingStart()
}
