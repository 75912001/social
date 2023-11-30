package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/debug"
	"social/pkg/bench"
	"social/pkg/common"
	"social/pkg/ec"
	"social/pkg/etcd"
	xrconstant "social/pkg/lib/constant"
	xrerror "social/pkg/lib/error"
	xretcd "social/pkg/lib/etcd"
	xrlog "social/pkg/lib/log"
	xrpprof "social/pkg/lib/pprof"
	xrtimer "social/pkg/lib/timer"
	xrutil "social/pkg/lib/util"
	"sync"
	"syscall"
	"time"
)

type Normal struct {
	Options     *options
	ZoneID      uint32 // 区域ID
	ServiceName string // 服务
	ServiceID   uint32 // 服务ID
	TimeMgr     xrutil.TimeMgr
	TimerMgr    *xrtimer.Mgr
	LogMgr      *xrlog.Mgr
	EtcdMgr     *xretcd.Mgr

	busChannel          chan interface{} //总线 channel
	busChannelWaitGroup sync.WaitGroup
	busCheckChan        chan struct{} // 检查总线channel,触发检查总线中的数据是否为0,且服务status == StatusStopping
	status              status        //服务状态
	exitChan            chan struct{}
}

func (p *Normal) Start(ctx context.Context) error {
	return xrerror.NotImplemented
}

func (p *Normal) Run(ctx context.Context) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case <-p.exitChan:
		p.LogMgr.Warn("Server will stop in a few seconds")
	case s := <-sigChan:
		p.LogMgr.Warnf("Server got signal: %s, shutting down...", s)
	}
	return nil
}

func (p *Normal) PreStop(ctx context.Context) error {
	return xrerror.NotImplemented
}

func (p *Normal) Stop(ctx context.Context) error {
	// 设置为关闭中
	p.SetStopping()

	// 定时检查事件总线是否消费完成
	go func() {
		p.LogMgr.Warn("start busCheckChan timer")
		idleDuration := 500 * time.Millisecond
		idleDelay := time.NewTimer(idleDuration)
		defer func() {
			idleDelay.Stop()
		}()
		for {
			select {
			case <-idleDelay.C:
				idleDelay.Reset(idleDuration)
				p.busCheckChan <- struct{}{}
				p.LogMgr.Warn("send to GCheckBusChan")
			}
		}
	}()

	// 等待GEventChan处理结束
	p.busChannelWaitGroup.Wait()

	xrtimer.GetInstance().Stop()
	p.LogMgr.Warn("server Timer stop")

	if xretcd.IsEnable() {
		_ = xretcd.GetInstance().Stop()
		p.LogMgr.Warn("server Etcd stop")
	}

	xrlog.PrintErr("server Log stop")
	_ = p.LogMgr.Stop()
	return nil
}

// Exit 退出服务
func (p *Normal) Exit() {
	p.LogMgr.Warn("server Exit")
	p.exitChan <- struct{}{}
}

func (p *Normal) Init(ctx context.Context, opts ...*options) error {
	p.busCheckChan = make(chan struct{}, 1)
	p.exitChan = make(chan struct{}, 1)

	rand.Seed(time.Now().UnixNano())
	p.TimeMgr.Update()
	// 小端
	if !xrutil.IsLittleEndian() {
		return errors.Errorf("system is bigEndian! %v", xrutil.GetCodeLocation(1).String())
	}
	// 开启UUID随机
	uuid.EnableRandPool()
	// 初始化 错误码
	if err := ec.Init(); err != nil {
		return errors.Errorf("ec Start err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	p.Options = mergeOptions(opts...)
	err := configure(p.Options)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	// 加载配置文件 bench.json 公共部分
	// 当前目录
	pathValue, err := xrutil.GetCurrentPath()
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	benchPath := path.Join(pathValue, *p.Options.benchPath)
	err = bench.GetInstance().Parse(benchPath)
	if err != nil {
		return errors.Errorf("Bench Load err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	if len(bench.GetInstance().Etcd.Key) == 0 {
		bench.GetInstance().Etcd.Key = fmt.Sprintf("%v/%v/%v/%v/%v",
			common.ProjectName, etcd.WatchMsgTypeService,
			p.ZoneID, p.ServiceName, p.ServiceID)
	}
	if bench.GetInstance().Etcd.TTL == 0 {
		bench.GetInstance().Etcd.TTL = etcd.TtlSecondDefault
	}
	if len(bench.GetInstance().Base.LogAbsPath) == 0 {
		bench.GetInstance().Base.LogAbsPath = common.LogAbsPath
	}
	//GoMaxProcess
	previous := runtime.GOMAXPROCS(bench.GetInstance().Base.GoMaxProcess)
	xrlog.PrintfInfo("go max process new:%v, previous setting:%v",
		bench.GetInstance().Base.GoMaxProcess, previous)
	// log
	p.LogMgr = xrlog.GetInstance()
	err = p.LogMgr.Start(ctx,
		xrlog.NewOptions().
			SetLevel(xrlog.Level(bench.GetInstance().Base.LogLevel)).
			SetAbsPath(bench.GetInstance().Base.LogAbsPath).
			SetNamePrefix(fmt.Sprintf("%v-%v-%v", p.ZoneID, p.ServiceName, p.ServiceID)),
	)
	if err != nil {
		return errors.Errorf("log Start err:%v %v ", err, xrutil.GetCodeLocation(1).String())
	}
	// 加载配置文件 bench.json 私有部分
	if p.Options.subBench != nil {
		err = p.Options.subBench.Load(benchPath)
		if err != nil {
			return errors.Errorf("GSubBench Load err:%v %v", err, xrutil.GetCodeLocation(1).String())
		}
	}
	// eventChan
	p.busChannel = make(chan interface{}, bench.GetInstance().Base.BusChannelNumber)
	go func() {
		defer func() {
			// 主事件channel报错 不recover
			p.LogMgr.Fatalf(xrconstant.GoroutineDone)
		}()
		p.busChannelWaitGroup.Add(1)
		defer p.busChannelWaitGroup.Done()

		p.HandleBus()
	}()
	// 是否开启http采集分析
	if 0 < bench.GetInstance().Base.PprofHttpPort {
		xrpprof.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", bench.GetInstance().Base.PprofHttpPort))
	}
	// 全局定时器
	p.TimerMgr = xrtimer.GetInstance()
	err = p.TimerMgr.Start(ctx,
		xrtimer.NewOptions().
			SetScanSecondDuration(bench.GetInstance().Timer.ScanSecondDuration).
			SetScanMillisecondDuration(bench.GetInstance().Timer.ScanMillisecondDuration).
			SetTimerOutChan(p.busChannel),
	)
	if err != nil {
		return errors.Errorf("timer Start err:%v %v ", err, xrutil.GetCodeLocation(1).String())
	}
	// 启动Etcd
	p.EtcdMgr = xretcd.GetInstance()
	err = etcd.Start(&bench.GetInstance().Etcd, p.busChannel, p.Options.etcdHandler)
	if err != nil {
		return errors.Errorf("Etcd start err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	// etcd 关注 服务 首次启动服务需要拉取一次
	if err = p.EtcdMgr.WatchPrefixIntoChan(*p.Options.etcdWatchServicePrefix); err != nil {
		return errors.Errorf("EtcdWatchPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	if err = p.EtcdMgr.GetPrefixIntoChan(*p.Options.etcdWatchServicePrefix); err != nil {
		return errors.Errorf("EtcdGetPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	// etcd 关注 命令
	if err = p.EtcdMgr.WatchPrefixIntoChan(*p.Options.etcdWatchCommandPrefix); err != nil {
		return errors.Errorf("EtcdWatchPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	p.serviceInformationPrintingStart()
	runtime.GC()

	return nil
}

func (p *Normal) serviceInformationPrintingStart() {
	p.TimerMgr.AddSecond(p.serviceInformationPrinting, nil, p.TimeMgr.ShadowTimeSecond()+ServiceInfoTimeOutSec)
}

// 服务信息 打印
func (p *Normal) serviceInformationPrinting(_ interface{}) {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	p.LogMgr.Infof("goroutineCnt:%d, busChannel:%d, numGC:%d, lastGC:%v, GCPauseTotal:%v",
		runtime.NumGoroutine(), len(p.busChannel), s.NumGC, s.LastGC, s.PauseTotal)
	p.serviceInformationPrintingStart()
}
