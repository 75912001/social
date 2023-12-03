package handler

import (
	"encoding/json"
	gateload "social/internal/gate/load"
	pkgbench "social/pkg/bench"
	libetcd "social/pkg/lib/etcd"
	liblog "social/pkg/lib/log"
	libutil "social/pkg/lib/time"
	libtimer "social/pkg/lib/timer"
	pkgserver "social/pkg/server"
)

func NewTimer() *Timer {
	return new(Timer)
}

type Timer struct {
	newDay   *libtimer.Second
	each5min *libtimer.Second
}

func (p *Timer) Start() {
	newDayBeginSec := libutil.DayBeginSec(&pkgserver.GetInstance().TimeMgr.Time) + libutil.OneDaySecond
	p.newDay = libtimer.GetInstance().AddSecond(newDayTimeOut, p, newDayBeginSec)
	p.each5min = libtimer.GetInstance().AddSecond(each5minTimeOut, p, pkgserver.GetInstance().TimeMgr.Second+libutil.OneMinuteSecond*5)
}

func (p *Timer) Stop() {
	libtimer.DelSecond(p.newDay)
	libtimer.DelSecond(p.each5min)
}

// 新的一天
func newDayTimeOut(arg interface{}) {
	p := arg.(*Timer)
	newDayBeginSec := libutil.DayBeginSec(&pkgserver.GetInstance().TimeMgr.Time) + libutil.OneDaySecond
	p.newDay = libtimer.GetInstance().AddSecond(newDayTimeOut, p, newDayBeginSec)
	//do...
}

// 每5分钟
func each5minTimeOut(arg interface{}) {
	p := arg.(*Timer)
	p.each5min = libtimer.GetInstance().AddSecond(each5minTimeOut, p, pkgserver.GetInstance().TimeMgr.Second+libutil.OneMinuteSecond*5)
	{ //更新load
		pkgbench.GetInstance().Etcd.Value.AvailableLoad = gateload.AvailableLoad()
		v, err := json.Marshal(pkgbench.GetInstance().Etcd.Value)
		if err != nil {
			liblog.GetInstance().Warnf("OnEventEtcd value json Marshal err:%v", err)
			return
		}
		_, _ = libetcd.GetInstance().PutWithLease(pkgbench.GetInstance().Etcd.Key, string(v))
	}
}
