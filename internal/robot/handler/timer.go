package handler

import (
	"encoding/json"
	"social/internal/gate/load"
	libetcd "social/lib/etcd"
	liblog "social/lib/log"
	libutil "social/lib/time"
	"social/lib/timer"
	"social/pkg/bench"
	"social/pkg/server"
)

type ServerTimer struct {
	newDay   *timer.Second
	each5min *timer.Second
}

func (p *ServerTimer) Start() {
	newDayBeginSec := libutil.DayBeginSec(&server.GetInstance().TimeMgr.Time) + libutil.OneDaySecond
	p.newDay = timer.GetInstance().AddSecond(newDayTimeOut, p, newDayBeginSec)
	p.each5min = timer.GetInstance().AddSecond(each5minTimeOut, p, server.GetInstance().TimeMgr.Second+libutil.OneMinuteSecond*5)
}

func (p *ServerTimer) Stop() {
	timer.DelSecond(p.newDay)
	timer.DelSecond(p.each5min)
}

// 新的一天
func newDayTimeOut(arg interface{}) {
	p := arg.(*ServerTimer)
	newDayBeginSec := libutil.DayBeginSec(&server.GetInstance().TimeMgr.Time) + libutil.OneDaySecond
	p.newDay = timer.GetInstance().AddSecond(newDayTimeOut, p, newDayBeginSec)
	//do...
}

// 每5分钟
func each5minTimeOut(arg interface{}) {
	p := arg.(*ServerTimer)
	p.each5min = timer.GetInstance().AddSecond(each5minTimeOut, p, server.GetInstance().TimeMgr.Second+libutil.OneMinuteSecond*5)
	{ //更新load
		bench.GetInstance().Etcd.Value.AvailableLoad = load.AvailableLoad()
		v, err := json.Marshal(bench.GetInstance().Etcd.Value)
		if err != nil {
			liblog.GetInstance().Warnf("OnEventEtcd value json Marshal err:%v", err)
			return
		}
		_, _ = libetcd.GetInstance().PutWithLease(bench.GetInstance().Etcd.Key, string(v))
	}
}
