package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math/rand"
	"path"
	"runtime"
	"runtime/debug"
	"social/pkg/bench"
	"social/pkg/common"
	"social/pkg/error_code"
	"social/pkg/etcd"
	xrconstant "social/pkg/lib/constant"
	xretcd "social/pkg/lib/etcd"
	xrlog "social/pkg/lib/log"
	xrpprof "social/pkg/lib/pprof"
	xrtimer "social/pkg/lib/timer"
	xrutil "social/pkg/lib/util"
	"sync"
	"time"
)

var (
	instance *mgr
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *mgr {
	once.Do(func() {
		instance = new(mgr)
	})
	return instance
}

type mgr struct {
	ZoneID      uint32 // 区域ID
	ServiceName string // 服务
	ServiceID   uint32 // 服务ID

	Opt *Options

	TimeMgr xrutil.TimeMgr

	BusChannel          chan interface{} //  总线 channel
	BusChannelWaitGroup sync.WaitGroup

	//

	status uint32

	QuitChan chan bool
	// 检查总线channel
	checkBusChan chan struct{}
}

// PreInit 初始化之前的操作
func (p *mgr) PreInit(ctx context.Context, opts ...*Options) error {
	p.checkBusChan = make(chan struct{}, 1)
	p.QuitChan = make(chan bool)

	rand.Seed(time.Now().UnixNano())
	p.TimeMgr.Update()
	// 小端
	if !xrutil.IsLittleEndian() {
		return errors.Errorf("system is bigEndian! %v", xrutil.GetCodeLocation(1).String())
	}
	// 开启UUID随机
	uuid.EnableRandPool()
	// 初始化 错误码
	if err := error_code.Init(); err != nil {
		return errors.Errorf("error_code Start err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	p.Opt = mergeOptions(opts...)
	err := configure(p.Opt)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	// 加载配置文件 bench.json 公共部分
	// 当前目录
	pathValue, err := xrutil.GetCurrentPath()
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	benchPath := path.Join(pathValue, *p.Opt.BenchPath)
	err = bench.GetInstance().Parse(benchPath)
	if err != nil {
		return errors.Errorf("Bench Load err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	//GoMaxProcess
	previous := runtime.GOMAXPROCS(bench.GetInstance().Base.GoMaxProcess)
	xrlog.PrintfInfo("go max process new:%v, previous setting:%v",
		bench.GetInstance().Base.GoMaxProcess, previous)
	// log
	err = xrlog.GetInstance().Start(context.TODO(),
		xrlog.NewOptions().
			SetLevel(xrlog.Level(bench.GetInstance().Base.LogLevel)).
			SetAbsPath(bench.GetInstance().Base.LogAbsPath).
			SetNamePrefix(fmt.Sprintf("%v-%v-%v", p.ZoneID, p.ServiceName, p.ServiceID)),
	)
	if err != nil {
		return errors.Errorf("log Start err:%v %v ", err, xrutil.GetCodeLocation(1).String())
	}
	// 加载配置文件 bench.json 私有部分
	if p.Opt.SubBench != nil {
		err = p.Opt.SubBench.Load(benchPath)
		if err != nil {
			return errors.Errorf("GSubBench Load err:%v %v", err, xrutil.GetCodeLocation(1).String())
		}
	}
	// eventChan
	p.BusChannel = make(chan interface{}, bench.GetInstance().Base.BusChannelNumber)
	go func() {
		defer func() {
			// 主事件channel报错 不recover
			xrlog.GetInstance().Fatalf(xrconstant.GoroutineDone)
		}()
		p.BusChannelWaitGroup.Add(1)
		defer p.BusChannelWaitGroup.Done()

		p.HandleBus()
	}()

	// 是否开启http采集分析
	if 0 < bench.GetInstance().Base.PprofHttpPort {
		xrpprof.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", bench.GetInstance().Base.PprofHttpPort))
	}

	// 全局定时器
	if bench.GetInstance().Timer.ScanSecondDuration != nil || bench.GetInstance().Timer.ScanMillisecondDuration != nil {
		err = xrtimer.GetInstance().Start(context.TODO(),
			xrtimer.NewOptions().
				SetScanSecondDuration(bench.GetInstance().Timer.ScanSecondDuration).
				SetScanMillisecondDuration(bench.GetInstance().Timer.ScanMillisecondDuration).
				SetTimerOutChan(p.BusChannel),
		)
		if err != nil {
			return errors.Errorf("timer Start err:%v %v ", err, xrutil.GetCodeLocation(1).String())
		}
	}

	runtime.GC()
	return nil
}

func (p *mgr) PostInit(ctx context.Context, opts ...*Options) error {
	p.Opt = mergeOptions(opts...)
	err := configure(p.Opt)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}
	// 启动Etcd
	bench.GetInstance().Etcd.Key = fmt.Sprintf("/%v/%v/%v/%v/%v",
		common.ProjectName, etcd.WatchMsgTypeService, p.ZoneID, p.ServiceName, p.ServiceID)
	err = etcd.Start(&bench.GetInstance().Etcd, p.BusChannel, p.Opt.EtcdHandler)
	if err != nil {
		return errors.Errorf("Etcd start err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	// etcd 关注 服务 首次启动服务需要拉取一次
	if p.Opt.EtcdWatchServicePrefix != nil {
		if err = xretcd.GetInstance().WatchPrefixIntoChan(*p.Opt.EtcdWatchServicePrefix); err != nil {
			return errors.Errorf("EtcdWatchPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
		}
		if err = xretcd.GetInstance().GetPrefixIntoChan(*p.Opt.EtcdWatchServicePrefix); err != nil {
			return errors.Errorf("EtcdGetPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
		}
	}
	// etcd 关注 命令
	if p.Opt.EtcdWatchCommandPrefix != nil {
		if err = xretcd.GetInstance().WatchPrefixIntoChan(*p.Opt.EtcdWatchCommandPrefix); err != nil {
			return errors.Errorf("EtcdWatchPrefix err:%v %v", err, xrutil.GetCodeLocation(1).String())
		}
	}
	serviceInformationPrintingStart()
	runtime.GC()
	return nil
}

func serviceInformationPrintingStart() {
	xrtimer.GetInstance().AddSecond(serviceInformationPrinting, nil, GetInstance().TimeMgr.ShadowTimeSecond()+ServiceInfoTimeOutSec)
}

// 服务信息 打印
func serviceInformationPrinting(_ interface{}) {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	xrlog.GetInstance().Infof("goroutineCnt:%d, BusChannel:%d, numGC:%d, lastGC:%v, GCPauseTotal:%v",
		runtime.NumGoroutine(), len(GetInstance().BusChannel), s.NumGC, s.LastGC, s.PauseTotal)
	serviceInformationPrintingStart()
}

func (p *mgr) Stop() error {
	// 定时检查事件总线是否消费完成
	go func() {
		xrlog.GetInstance().Warn("start checkGBusChannel timer")

		idleDuration := 500 * time.Millisecond
		idleDelay := time.NewTimer(idleDuration)
		defer func() {
			idleDelay.Stop()
		}()

		for {
			select {
			case <-idleDelay.C:
				idleDelay.Reset(idleDuration)
				p.checkBusChan <- struct{}{}
				xrlog.GetInstance().Warn("send to GCheckBusChan")
			}
		}
	}()

	// 等待GEventChan处理结束
	p.BusChannelWaitGroup.Wait()

	if bench.GetInstance().Timer.ScanSecondDuration != nil || bench.GetInstance().Timer.ScanMillisecondDuration != nil {
		xrtimer.GetInstance().Stop()
		xrlog.GetInstance().Warn("GTimer stop")
	}
	if xretcd.IsEnable() {
		_ = xretcd.GetInstance().Stop()
		xrlog.GetInstance().Warn("GEtcd stop")
	}
	return nil
}
