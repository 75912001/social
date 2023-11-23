package handler

import (
	"encoding/json"
	"social/internal/gate/load"
	"social/pkg/bench"
	xretcd "social/pkg/lib/etcd"
	xrlog "social/pkg/lib/log"
	xrtimer "social/pkg/lib/timer"
	xrutil "social/pkg/lib/util"
	"social/pkg/server"
)

type ServerTimer struct {
	newDay   *xrtimer.Second
	each5min *xrtimer.Second
}

func (p *ServerTimer) Start() {
	newDayBeginSec := xrutil.DayBeginSec(&server.GetInstance().TimeMgr.Time) + xrutil.OneDaySecond
	p.newDay = xrtimer.GetInstance().AddSecond(newDayTimeOut, p, newDayBeginSec)
	p.each5min = xrtimer.GetInstance().AddSecond(each5minTimeOut, p, server.GetInstance().TimeMgr.Second+xrutil.OneMinuteSecond*5)
}

func (p *ServerTimer) Stop() {
	xrtimer.DelSecond(p.newDay)
	xrtimer.DelSecond(p.each5min)
}

// 新的一天
func newDayTimeOut(arg interface{}) {
	p := arg.(*ServerTimer)
	newDayBeginSec := xrutil.DayBeginSec(&server.GetInstance().TimeMgr.Time) + xrutil.OneDaySecond
	p.newDay = xrtimer.GetInstance().AddSecond(newDayTimeOut, p, newDayBeginSec)
	//do...
}

// 每5分钟
func each5minTimeOut(arg interface{}) {
	p := arg.(*ServerTimer)
	p.each5min = xrtimer.GetInstance().AddSecond(each5minTimeOut, p, server.GetInstance().TimeMgr.Second+xrutil.OneMinuteSecond*5)
	{ //更新load
		bench.GetInstance().Etcd.Value.AvailableLoad = load.AvailableLoad()
		v, err := json.Marshal(bench.GetInstance().Etcd.Value)
		if err != nil {
			xrlog.GetInstance().Warnf("OnEventEtcd value json Marshal err:%v", err)
			return
		}
		_, _ = xretcd.GetInstance().PutWithLease(bench.GetInstance().Etcd.Key, string(v))
	}
}
