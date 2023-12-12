package time

import (
	"strconv"
	"sync"
	"time"
)

// GenYMD 获取 e.g.:20210819
//
//	返回YMD
func GenYMD(sec int64) int {
	strYMD := time.Unix(sec, 0).Format("20060102")
	ymd, _ := strconv.Atoi(strYMD)
	return ymd
}

var (
	instance *Mgr
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *Mgr {
	once.Do(func() {
		instance = new(Mgr)
	})
	return instance
}

// Mgr 时间管理器
type Mgr struct {
	Second       int64     //近似时间（秒），上一次调用Update更新的时间
	Millisecond  int64     //近似时间（毫秒），上一次调用Update更新的时间
	Time         time.Time //上一次调用Update更新的时间
	SecondOffset int64     //时间偏移量-秒
}

// Update 更新时间管理器中的,当前时间
func (p *Mgr) Update() {
	p.Time = time.Now().UTC()
	p.Second = p.Time.Unix()
	p.Millisecond = p.Time.UnixMilli() // UnixNano() / int64(time.Millisecond)
}

// TimeSecond 秒
func (p *Mgr) TimeSecond() int64 {
	return p.Second
}

// ShadowTimeSecond 叠加偏移量的秒
func (p *Mgr) ShadowTimeSecond() int64 {
	return p.Second + p.SecondOffset
}

// DayBeginSec 当天开始时间戳
func DayBeginSec(t *time.Time) int64 {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}

// DayEndSec 当天最后一秒
func DayEndSec(dayBeginSec int64) int64 {
	return dayBeginSec + OneDaySecond - 1
}
