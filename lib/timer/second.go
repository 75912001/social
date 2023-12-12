package timer

// Second 秒级定时器
type Second struct {
	Millisecond
}

// DelSecond 删除秒级定时器
// 同 DelMillisecond
func DelSecond(t *Second) {
	t.Millisecond.inValid()
}
