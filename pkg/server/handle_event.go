package server

import (
	xretcd "social/pkg/lib/etcd"
	xrlog "social/pkg/lib/log"
	xrtimer "social/pkg/lib/timer"
	xrutil "social/pkg/lib/util"
	"time"
)

// HandleBus todo [重要] issue 在处理 event 时候, 向 eventChan 中插入 事件，注意超出eventChan的上限会阻塞.
func (p *mgr) HandleBus() {
	// 在消费eventChan时可能会往eventChan中写入事件,所以关闭服务时不能close eventChan(造成写入阻塞),通过定时检查eventChan大小来关闭
	for {
		select {
		case <-GBusChannelCheckChan:
			xrlog.GetInstance().Warn("receive GBusChannelCheckChan")
			if 0 == len(p.BusChannel) && p.IsStopping() {
				xrlog.GetInstance().Warn("server is stopping, stop consume GEventChan with length 0")
				return
			} else {
				xrlog.GetInstance().Warnf("server is stopping, waiting for consume GEventChan with length:%d", len(p.BusChannel))
			}
		case v := <-p.BusChannel:
			//TODO [*] 应拿尽拿...
			p.TimeMgr.Update()
			var err error
			switch t := v.(type) {
			case *xrtimer.Second:
				if t.IsValid() {
					t.Function(t.Arg)
				}
			case *xrtimer.Millisecond:
				if t.IsValid() {
					t.Function(t.Arg)
				}
			case *xretcd.KV:
				err = xretcd.GetInstance().Handler(t.Key, t.Value)
			default:
				if p.Opt.DefaultHandler == nil {
					xrlog.GetInstance().Fatalf("non-existent event:%v %v", v, t)
				} else {
					err = p.Opt.DefaultHandler(v)
				}
			}
			if err != nil {
				xrlog.PrintErr(v, err)
			}
			if xrutil.IsDebug() {
				dt := time.Now().Sub(p.TimeMgr.Time).Milliseconds()
				if dt > 50 {
					xrlog.GetInstance().Warnf("cost time50: %v Millisecond with event type:%T", dt, v)
				} else if dt > 20 {
					xrlog.GetInstance().Warnf("cost time20: %v Millisecond with event type:%T", dt, v)
				} else if dt > 10 {
					xrlog.GetInstance().Warnf("cost time10: %v Millisecond with event type:%T", dt, v)
				}
			}
		}
	}
}
