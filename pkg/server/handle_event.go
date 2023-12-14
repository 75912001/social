package server

import (
	libetcd "social/lib/etcd"
	liblog "social/lib/log"
	libtime "social/lib/time"
	"social/lib/timer"
	libutil "social/lib/util"
)

// HandleBus todo [重要] issue 在处理 event 时候, 向 eventChan 中插入 事件，注意超出eventChan的上限会阻塞.
func (p *Normal) HandleBus() {
	// 在消费eventChan时可能会往eventChan中写入事件,所以关闭服务时不能close eventChan(造成写入阻塞),通过定时检查eventChan大小来关闭
	for {
		select {
		case <-p.busCheckChan:
			p.LogMgr.Warn("receive busCheckChan")
			if 0 == len(p.busChannel) && p.IsStopping() {
				p.LogMgr.Warn("server is stopping, stop consume EventChan with length 0")
				return
			} else {
				p.LogMgr.Warnf("server is stopping, waiting for consume EventChan with length:%d", len(p.busChannel))
			}
		case v := <-p.busChannel:
			//TODO [*] 应拿尽拿...
			p.TimeMgr.Update()
			var err error
			switch t := v.(type) {
			case *timer.Second:
				if t.IsValid() {
					t.Function(t.Arg)
				}
			case *timer.Millisecond:
				if t.IsValid() {
					t.Function(t.Arg)
				}
			case *libetcd.KV:
				err = p.EtcdMgr.Handler(t.Key, t.Value)
			default:
				if p.Options.defaultHandler == nil {
					p.LogMgr.Fatalf("non-existent event:%v %v", v, t)
				} else {
					err = p.Options.defaultHandler(v)
				}
			}
			if err != nil {
				liblog.PrintErr(v, err)
			}
			if libutil.IsDebug() {
				dt := libtime.NowTime().Sub(p.TimeMgr.Time).Milliseconds()
				if dt > 50 {
					p.LogMgr.Warnf("cost time50: %v Millisecond with event type:%T", dt, v)
				} else if dt > 20 {
					p.LogMgr.Warnf("cost time20: %v Millisecond with event type:%T", dt, v)
				} else if dt > 10 {
					p.LogMgr.Warnf("cost time10: %v Millisecond with event type:%T", dt, v)
				}
			}
		}
	}
}
