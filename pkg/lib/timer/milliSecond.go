package timer

// Millisecond 毫秒级定时器
type Millisecond struct {
	Arg      interface{} // 参数
	Function OnFun       // 超时调用的函数
	expire   int64       // 过期时间戳
	valid    bool        // 有效(false:不执行,扫描时自动删除)
}

// IsValid 判断是否有效
func (p *Millisecond) IsValid() bool {
	return p.valid
}

func (p *Millisecond) reset() {
	p.Arg = nil
	p.Function = nil
	p.expire = 0
	p.valid = false
}

// DelMillisecond 删除毫秒级定时器
//	NOTE 必须与该timerOutChan线性处理.如:在同一个goroutine select中处理数据
//	参数:
//		毫秒定时器
func DelMillisecond(t *Millisecond) {
	t.inValid()
}

// 设为无效
func (p *Millisecond) inValid() {
	p.Arg = nil
	p.Function = nil
	p.expire = 0
	p.valid = false
}
