package server

import "sync/atomic"

const StatusRunning uint32 = 0  // 服务状态：运行中
const StatusStopping uint32 = 1 // 服务状态：关闭中

// IsStopping 服务是否关闭中
func (p *mgr) IsStopping() bool {
	return atomic.LoadUint32(&p.status) == StatusStopping
}

// IsRunning 服务是否运行中
func (p *mgr) IsRunning() bool {
	return atomic.LoadUint32(&p.status) == StatusRunning
}

// SetStopping 设置为关闭中
func (p *mgr) SetStopping() {
	atomic.StoreUint32(&p.status, StatusStopping)
}
