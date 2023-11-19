package util

import (
	"strconv"
	"time"
)

const OneDaySecond int64 = 86400
const OneMinuteSecond int64 = 60

// GenYMD 获取 e.g.:20210819
//
//	返回YMD
func GenYMD(sec int64) int {
	strYMD := time.Unix(sec, 0).Format("20060102")
	ymd, _ := strconv.Atoi(strYMD)
	return ymd
}

// TimeMgr 时间管理器
type TimeMgr struct {
	Second       int64     //近似时间（秒），上一次调用Update更新的时间
	Millisecond  int64     //近似时间（毫秒），上一次调用Update更新的时间
	Time         time.Time //上一次调用Update更新的时间
	SecondOffset int64     //时间偏移量-秒
}

// Update 更新时间管理器中的,当前时间
func (p *TimeMgr) Update() {
	p.Time = time.Now()
	p.Second = p.Time.Unix()
	p.Millisecond = p.Time.UnixMilli() // UnixNano() / int64(time.Millisecond)
}

// TimeSecond 秒
func (p *TimeMgr) TimeSecond() int64 {
	return p.Second
}

// ShadowTimeSecond 叠加偏移量的秒
func (p *TimeMgr) ShadowTimeSecond() int64 {
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
