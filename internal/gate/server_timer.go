package gate

import (
	"context"
	"encoding/json"
	libutil "social/lib/time"
	libtimer "social/lib/timer"
	pkgserver "social/pkg/server"
)

type ServerTimer struct {
	newDay   *libtimer.Second
	each5min *libtimer.Second
}

func (p *ServerTimer) Start() {
	newDayBeginSec := libutil.DayBeginSec(server.TimeMgr.ShadowTimeSecond()) + libutil.OneDaySecond
	p.newDay = server.TimerMgr.AddSecond(newDayTimeOut, nil, newDayBeginSec)
	p.each5min = server.TimerMgr.AddSecond(each5minTimeOut, nil, server.TimeMgr.ShadowTimeSecond()+libutil.OneMinuteSecond*5)
}

func (p *ServerTimer) Stop() {
	libtimer.DelSecond(p.newDay)
	libtimer.DelSecond(p.each5min)
}

// 新的一天
func newDayTimeOut(_ interface{}) {
	newDayBeginSec := libutil.DayBeginSec(server.TimeMgr.ShadowTimeSecond()) + libutil.OneDaySecond
	server.serverTimer.newDay = server.TimerMgr.AddSecond(newDayTimeOut, nil, newDayBeginSec)
	//do...
}

// 每5分钟
func each5minTimeOut(_ interface{}) {
	server.serverTimer.each5min = server.TimerMgr.AddSecond(each5minTimeOut, nil, server.TimeMgr.ShadowTimeSecond()+libutil.OneMinuteSecond*5)
	{ //更新load
		server.BenchMgr.Etcd.Value.AvailableLoad = pkgserver.AvailableLoad()
		v, err := json.Marshal(server.BenchMgr.Etcd.Value)
		if err != nil {
			server.LogMgr.Warnf("OnEventEtcd value json Marshal err:%v", err)
			return
		}
		putResponse, err := server.EtcdMgr.PutWithLease(context.Background(), server.BenchMgr.Etcd.Key, string(v))
		if err != nil {
			server.LogMgr.Warn(putResponse, err)
		}
	}
}
