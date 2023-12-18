package robot

// Bus 总线
type Bus struct {
}

// OnEventBus 处理事件-总线
func (p *Bus) OnEventBus(v interface{}) error {
	switch t := v.(type) {
	default:
		robot.LogMgr.Errorf("non-existent event:%v %v", v, t)
	}
	return nil
}
