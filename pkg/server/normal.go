package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	libconstant "social/lib/consts"
	liberror "social/lib/error"
	libetcd "social/lib/etcd"
	"social/lib/log"
	libpprof "social/lib/pprof"
	libtime "social/lib/time"
	"social/lib/timer"
	libutil "social/lib/util"
	pkgbench "social/pkg/bench"
	pkgec "social/pkg/ec"
	pkgetcd "social/pkg/etcd"
	"sync"
	"syscall"
	"time"
)

var (
	instance *Normal
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *Normal {
	once.Do(func() {
		instance = NewNormal()
	})
	return instance
}

func NewNormal() *Normal {
	normal := new(Normal)

	normal.Options = &options{}

	normal.BenchMgr = pkgbench.GetInstance()
	normal.TimeMgr = libtime.GetInstance()
	normal.TimerMgr = timer.GetInstance()
	normal.LogMgr = log.GetInstance()
	normal.EtcdMgr = libetcd.GetInstance()

	return normal
}

type Normal struct {
	Options     *options
	ProgramPath string // 程序路径
	ProgramName string // 程序名称
	ZoneID      uint32 // 区域ID
	ServiceName string // 服务
	ServiceID   uint32 // 服务ID

	BenchMgr *pkgbench.Mgr
	TimeMgr  *libtime.Mgr
	TimerMgr *timer.Mgr
	LogMgr   *log.Mgr
	EtcdMgr  *libetcd.Mgr

	busChannel          chan interface{} //总线 channel
	busChannelWaitGroup sync.WaitGroup
	busCheckChan        chan struct{} // 检查总线channel,触发检查总线中的数据是否为0,且服务status == StatusStopping
	status              status        //服务状态
	exitChan            chan struct{}
}

func (p *Normal) LoadBench(ctx context.Context, opts ...*options) error {
	p.Options = mergeOptions(opts...)
	err := configure(p.Options)
	if err != nil {
		return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	// 加载配置文件 bench.json 公共部分
	err = p.BenchMgr.Parse(*p.Options.benchPath, p.ZoneID, p.ServiceName, p.ServiceID)
	if err != nil {
		return errors.Errorf("Bench Load err:%v %v", err, libutil.GetCodeLocation(1).String())
	}
	// 加载配置文件 bench.json 私有部分
	if p.Options.subBench != nil {
		err = p.Options.subBench.Load(*p.Options.benchPath)
		if err != nil {
			return errors.Errorf("SubBench Load err:%v %v", err, libutil.GetCodeLocation(1).String())
		}
	}
	return nil
}

func (p *Normal) Init(ctx context.Context, opts ...*options) error {
	p.busCheckChan = make(chan struct{}, 1)
	p.exitChan = make(chan struct{}, 1)

	rand.Seed(libtime.NowTime().UnixNano())
	p.TimeMgr.Update()
	// 小端
	if !libutil.IsLittleEndian() {
		return errors.Errorf("system is bigEndian! %v", libutil.GetCodeLocation(1).String())
	}
	// 开启UUID随机
	uuid.EnableRandPool()
	// 初始化 错误码
	if err := pkgec.Init(); err != nil {
		return errors.Errorf("ec Start err:%v %v", err, libutil.GetCodeLocation(1).String())
	}
	p.Options = mergeOptions(opts...)
	err := configure(p.Options)
	if err != nil {
		return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
	}
	//GoMaxProcess
	previous := runtime.GOMAXPROCS(p.BenchMgr.Base.GoMaxProcess)
	log.PrintfInfo("go max process new:%v, previous setting:%v",
		p.BenchMgr.Base.GoMaxProcess, previous)
	// log
	err = p.LogMgr.Start(ctx,
		log.NewOptions().
			WithLevel(log.Level(p.BenchMgr.Base.LogLevel)).
			WithAbsPath(p.BenchMgr.Base.LogAbsPath).
			WithNamePrefix(fmt.Sprintf("%v-%v-%v", p.ZoneID, p.ServiceName, p.ServiceID)),
	)
	if err != nil {
		return errors.Errorf("log Start err:%v %v ", err, libutil.GetCodeLocation(1).String())
	}

	// eventChan
	p.busChannel = make(chan interface{}, p.BenchMgr.Base.BusChannelNumber)
	go func() {
		defer func() {
			// 主事件channel报错 不recover
			p.LogMgr.Fatalf(libconstant.GoroutineDone)
		}()
		p.busChannelWaitGroup.Add(1)
		defer p.busChannelWaitGroup.Done()

		p.HandleBus()
	}()
	// 是否开启http采集分析
	if 0 < p.BenchMgr.Base.PprofHttpPort {
		libpprof.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", p.BenchMgr.Base.PprofHttpPort))
	}
	// 全局定时器
	err = p.TimerMgr.Start(ctx,
		timer.NewOptions().
			SetScanSecondDuration(p.BenchMgr.Timer.ScanSecondDuration).
			SetScanMillisecondDuration(p.BenchMgr.Timer.ScanMillisecondDuration).
			SetOutgoingTimerOutChan(p.busChannel),
	)
	if err != nil {
		return errors.Errorf("timer Start err:%v %v ", err, libutil.GetCodeLocation(1).String())
	}
	// 启动Etcd
	err = pkgetcd.Start(&p.BenchMgr.Etcd, p.busChannel, p.Options.etcdHandler)
	if err != nil {
		return errors.Errorf("Etcd start err:%v %v", err, libutil.GetCodeLocation(1).String())
	}

	p.serviceInformationPrintingStart()
	runtime.GC()

	return nil
}

func (p *Normal) Start(ctx context.Context) error {
	return liberror.NotImplemented
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
	return liberror.NotImplemented
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

	p.TimerMgr.Stop()
	p.LogMgr.Warn("server Timer stop")

	if libetcd.IsEnable() {
		_ = p.EtcdMgr.Stop()
		p.LogMgr.Warn("server Etcd stop")
	}

	log.PrintErr("server Log stop")
	_ = p.LogMgr.Stop()
	return nil
}

// Exit 退出服务
func (p *Normal) Exit() {
	p.LogMgr.Warn("server Exit")
	p.exitChan <- struct{}{}
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
