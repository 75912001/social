package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"runtime"
	libbench "social/lib/bench"
	libconsts "social/lib/consts"
	liberror "social/lib/error"
	libetcd "social/lib/etcd"
	liblog "social/lib/log"
	libpprof "social/lib/pprof"
	libruntime "social/lib/runtime"
	libtime "social/lib/time"
	libtimer "social/lib/timer"
	libutil "social/lib/util"
	pkgconsts "social/pkg/consts"
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

	normal.Options = &Options{}

	normal.BenchMgr = libbench.GetInstance()
	normal.TimeMgr = libtime.GetInstance()
	normal.TimerMgr = libtimer.GetInstance()
	normal.LogMgr = liblog.GetInstance()
	normal.EtcdMgr = &libetcd.Mgr{}

	return normal
}

type Normal struct {
	Options     *Options
	ProgramPath string // 程序路径
	ProgramName string // 程序名称
	ZoneID      uint32 // 区域ID
	ServiceName string // 服务
	ServiceID   uint32 // 服务ID

	BenchMgr *libbench.Mgr
	TimeMgr  *libtime.Mgr
	TimerMgr *libtimer.Mgr
	LogMgr   *liblog.Mgr
	EtcdMgr  *libetcd.Mgr

	busChannel          chan interface{} //总线 channel
	busChannelWaitGroup sync.WaitGroup
	busCheckChan        chan struct{} // 检查总线channel,触发检查总线中的数据是否为0,且服务status == StatusStopping
	status              status        //服务状态
	exitChan            chan struct{}
	GrpcServer          *grpc.Server

	informationPrintingTimerSecond *libtimer.Second
}

func (p *Normal) OnLoadBench(_ context.Context, opts ...*Options) error {
	p.Options = mergeOptions(opts...)
	err := configure(p.Options)
	if err != nil {
		return errors.WithMessage(err, libruntime.GetCodeLocation(1).String())
	}
	// 加载配置文件 bench.json 公共部分
	benchPath := path.Join(p.ProgramPath, "bench.json")
	err = p.BenchMgr.Parse(benchPath, pkgconsts.ProjectName, p.ZoneID, p.ServiceName, p.ServiceID)
	if err != nil {
		return errors.Errorf("Bench Load err:%v %v", err, libruntime.GetCodeLocation(1).String())
	}
	if p.Options.subBench != nil {
		// 加载配置文件 bench.json 私有部分
		subbenchPath := path.Join(p.ProgramPath, "bench.json")
		err = p.Options.subBench.Parse(subbenchPath)
		if err != nil {
			return errors.Errorf("SubBench Load err:%v %v", err, libruntime.GetCodeLocation(1).String())
		}
	}
	return nil
}

func (p *Normal) OnInit(ctx context.Context, opts ...*Options) error {
	rand.Seed(libtime.NowTime().UnixNano())
	p.TimeMgr.Update()
	// 小端
	if !libutil.IsLittleEndian() {
		return errors.Errorf("system is bigEndian! %v", libruntime.GetCodeLocation(1).String())
	}
	// 开启UUID随机
	uuid.EnableRandPool()

	p.Options = mergeOptions(opts...)
	err := configure(p.Options)
	if err != nil {
		return errors.WithMessage(err, libruntime.GetCodeLocation(1).String())
	}
	//GoMaxProcess
	previous := runtime.GOMAXPROCS(p.BenchMgr.Base.GoMaxProcess)
	liblog.PrintfInfo("go max process new:%v, previous setting:%v", p.BenchMgr.Base.GoMaxProcess, previous)
	// log
	err = p.LogMgr.Start(ctx,
		liblog.NewOptions().
			WithLevel(liblog.Level(p.BenchMgr.Base.LogLevel)).
			WithAbsPath(p.BenchMgr.Base.LogAbsPath).
			WithNamePrefix(fmt.Sprintf("%v-%v-%v", p.ZoneID, p.ServiceName, p.ServiceID)),
	)
	if err != nil {
		return errors.Errorf("log OnStart err:%v %v ", err, libruntime.GetCodeLocation(1).String())
	}

	p.busCheckChan = make(chan struct{}, 1)
	p.exitChan = make(chan struct{}, 1)
	p.busChannel = make(chan interface{}, p.BenchMgr.Base.BusChannelNumber)
	go func() {
		defer func() {
			// 主事件channel报错 不recover
			p.LogMgr.Fatalf(libconsts.GoroutineDone)
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
		libtimer.NewOptions().
			WithScanSecondDuration(p.BenchMgr.Timer.ScanSecondDuration).
			WithScanMillisecondDuration(p.BenchMgr.Timer.ScanMillisecondDuration).
			WithOutgoingTimerOutChan(p.busChannel),
	)
	if err != nil {
		return errors.Errorf("timer OnStart err:%v %v ", err, libruntime.GetCodeLocation(1).String())
	}
	// 启动Etcd
	etcdValue, err := json.Marshal(p.BenchMgr.Etcd.Value)
	if err != nil {
		return errors.WithMessagef(err, libruntime.Location())
	}
	kvSlice := []libetcd.KV{
		{
			Key:   p.BenchMgr.Etcd.Key,
			Value: string(etcdValue),
		},
	}
	err = p.EtcdMgr.Start(context.TODO(),
		libetcd.NewOptions().
			WithAddrs(p.BenchMgr.Etcd.Addrs).
			WithTTL(p.BenchMgr.Etcd.TTL).
			WithKV(kvSlice).
			WithOnFunc(p.Options.etcdHandler).
			WithOutgoingEventChan(p.busChannel).
			WithWatchServicePrefix(libetcd.GenerateWatchServicePrefix(pkgconsts.ProjectName)).
			WithWatchCommandPrefix(libetcd.GenerateWatchCommandPrefix(pkgconsts.ProjectName, p.ZoneID, p.ServiceName)),
	)
	if err != nil {
		return errors.WithMessagef(err, libruntime.Location())
	}
	// 续租
	err = p.EtcdMgr.Run(context.TODO())
	if err != nil {
		return errors.WithMessagef(err, libruntime.Location())
	}
	p.serviceInformationPrintingStart()
	if p.Options.timerEachSecond != nil {
		p.Options.timerEachSecond.lastExpireSecond = p.TimeMgr.TimeSecond() + 1
		p.Options.timerEachSecond.onTimerFunHandle = p.TimerMgr.AddSecond(p.onTimerEachSecond, p.Options.timerEachSecond.Arg, p.Options.timerEachSecond.lastExpireSecond)
	}
	if p.Options.timerEachDay != nil {
		p.Options.timerEachDay.lastExpireSecond = libtime.DayBeginSec(p.TimeMgr.TimeSecond()) + libtime.OneDaySecond
		p.Options.timerEachDay.onTimerFunHandle = p.TimerMgr.AddSecond(p.onTimerEachDay, p.Options.timerEachDay, p.Options.timerEachDay.lastExpireSecond)
	}

	runtime.GC()
	return nil
}

func (p *Normal) OnStart(_ context.Context) error {
	return liberror.NotImplemented
}

func (p *Normal) OnRun(_ context.Context) error {
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

// Exit 退出服务
func (p *Normal) Exit() {
	if p.GrpcServer != nil {
		p.GrpcServer.GracefulStop()
	}
	p.LogMgr.Warn("server Exit")
	p.exitChan <- struct{}{}
}

func (p *Normal) OnPreStop(_ context.Context) error {
	return liberror.NotImplemented
}

func (p *Normal) OnStop(_ context.Context) error {
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

	if p.Options.timerEachSecond != nil {
		libtimer.DelSecond(p.Options.timerEachSecond.onTimerFunHandle)
	}
	if p.Options.timerEachDay != nil {
		libtimer.DelSecond(p.Options.timerEachDay.onTimerFunHandle)
	}
	if p.informationPrintingTimerSecond != nil {
		libtimer.DelSecond(p.informationPrintingTimerSecond)
	}
	p.TimerMgr.Stop()
	p.LogMgr.Warn("server Timer stop")

	err := p.EtcdMgr.Stop()
	p.LogMgr.Warn(err, "server Etcd stop")

	err = p.LogMgr.Stop()
	liblog.PrintInfo(err, "server Log stop")
	return nil
}
