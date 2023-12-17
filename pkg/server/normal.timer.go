package server

import (
	libruntime "social/lib/runtime"
	libtime "social/lib/time"
	libtimer "social/lib/timer"
	libutil "social/lib/util"
)

type NormalTimerSecond struct {
	OnTimerFun       libtimer.OnFun
	onTimerFunHandle *libtimer.Second
	lastExpireSecond int64 // 上一个到期时间
	Arg              libutil.IObject
}

func (p *Normal) onTimerEachSecond(arg interface{}) {
	p.TimeMgr.Update()
	//p.LogMgr.Trace(p.TimeMgr.ShadowTimeSecond(), libruntime.Location())
	if p.Options.timerEachSecond != nil {
		p.Options.timerEachSecond.OnTimerFun(arg)
	}
	//使用上次的时间来确定下次的时间,这样不会堆积
	p.Options.timerEachSecond.lastExpireSecond += 1
	p.Options.timerEachSecond.onTimerFunHandle = p.TimerMgr.AddSecond(p.onTimerEachSecond, arg, p.Options.timerEachSecond.lastExpireSecond)
}

func (p *Normal) onTimerEachDay(arg interface{}) {
	p.TimeMgr.Update()
	p.LogMgr.Trace(p.TimeMgr.ShadowTimeSecond(), libruntime.Location())
	if p.Options.timerEachDay != nil {
		p.Options.timerEachDay.OnTimerFun(arg)
	}
	//使用上次的时间来确定下次的时间,这样不会堆积
	p.Options.timerEachDay.lastExpireSecond += libtime.OneDaySecond
	p.Options.timerEachDay.onTimerFunHandle = p.TimerMgr.AddSecond(p.onTimerEachDay, arg, p.Options.timerEachDay.lastExpireSecond)
}
