package friend

import "social/pkg/server"

// Bus 总线
type Bus struct {
	*server.Normal
}

// OnEventBus 处理事件-总线
func (p *Bus) OnEventBus(v interface{}) error {
	switch t := v.(type) {
	default:
		p.Normal.LogMgr.Errorf("non-existent event:%v %v", v, t)
	}
	return nil
}
