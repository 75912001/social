package gate

import (
	"context"
	"encoding/json"
	pkgserver "social/pkg/server"
)

var timerAvailableLoadExpireTimestamp int64 //到期-时间戳
// 定时器-可用负载
func onTimerAvailableLoad(_ interface{}) {
	if timerAvailableLoadExpireTimestamp <= app.TimeMgr.ShadowTimeSecond() { //更新load
		timerAvailableLoadExpireTimestamp += 60
		app.BenchMgr.Etcd.Value.AvailableLoad = pkgserver.AvailableLoad()
		v, err := json.Marshal(app.BenchMgr.Etcd.Value)
		if err != nil {
			app.LogMgr.Warnf("OnEventEtcd value json Marshal err:%v", err)
			return
		}
		putResponse, err := app.EtcdMgr.PutWithLease(context.Background(), app.BenchMgr.Etcd.Key, string(v))
		if err != nil {
			app.LogMgr.Warn(putResponse, err)
		}
	}
}

// OnTimerEachSecondFun 新的一秒
func (p *Gate) OnTimerEachSecondFun(arg interface{}) {
	//do...
	onTimerAvailableLoad(arg)
}

// OnTimerEachDayFun 新的一天
func (p *Gate) OnTimerEachDayFun(_ interface{}) {
	//do...
}
