package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math/rand"
	"path"
	"runtime"
	"social/pkg/bench"
	"social/pkg/common"
	"social/pkg/error_code"
	"social/pkg/etcd"
	xrconstant "social/pkg/lib/constant"
	xretcd "social/pkg/lib/etcd"
	xrlog "social/pkg/lib/log"
	xrpporf "social/pkg/lib/pprof"
	xrtimer "social/pkg/lib/timer"
	xrutil "social/pkg/lib/util"
	"sync"
	"sync/atomic"
	"time"
)

const ServiceNameGate string = "gate"                     //网关
const ServiceNameFriend string = "friend"                 //好友
const ServiceNameInteraction string = "interaction"       //交互
const ServiceNameNotification string = "notification"     //通知
const ServiceNameBlog string = "blog"                     //博客
const ServiceNameRecommendation string = "recommendation" //推荐
const ServiceNameCleansing string = "cleansing"           //清洗

const StatusRunning = 0  // 服务状态：运行中
const StatusStopping = 1 // 服务状态：关闭中

var GBusChannelWaitGroup sync.WaitGroup
var GBusChannelCheckChan = make(chan struct{}, 1)

var GServerStatus uint32

var GQuitChan = make(chan bool)

// IsServerStopping 服务是否关闭中
func IsServerStopping() bool {
	return atomic.LoadUint32(&GServerStatus) == StatusStopping
}

// IsServerRunning 服务是否运行中
func IsServerRunning() bool {
	return atomic.LoadUint32(&GServerStatus) == StatusRunning
}

// SetServerStopping 设置为关闭中
func SetServerStopping() {
	atomic.StoreUint32(&GServerStatus, StatusStopping)
}

type IServer interface {
	Start() (err error)
	Stop() (err error)
}

var GMgr Mgr

type Mgr struct {
	ZoneID      uint32 // 区域ID
	ServiceName string // 服务
	ServiceID   uint32 // 服务ID

	Opt     *Options
	TimeMgr xrutil.TimeMgr
	Timer   xrtimer.Mgr
	Bench   bench.Mgr

	BusChannel          chan interface{} //  总线 channel
	BusChannelWaitGroup sync.WaitGroup
}

func (p *Mgr) PreInit(ctx context.Context, opts ...*Options) error {
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
	err = p.Bench.Parse(benchPath)
	if err != nil {
		return errors.Errorf("Bench Load err:%v %v", err, xrutil.GetCodeLocation(1).String())
	}
	//GoMaxProcess
	previous := runtime.GOMAXPROCS(p.Bench.Json.Base.GoMaxProcess)
	xrlog.PrintfInfo("go max process new:%v, previous setting:%v",
		p.Bench.Json.Base.GoMaxProcess, previous)
	// log
	err = xrlog.GetInstance().Start(context.TODO(),
		xrlog.NewOptions().
			SetLevel(xrlog.Level(p.Bench.Json.Base.LogLevel)).
			SetAbsPath(p.Bench.Json.Base.LogAbsPath).
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
	p.BusChannel = make(chan interface{}, p.Bench.Json.Base.BusChannelNumber)

	go func() {
		defer func() {
			// 主事件channel报错 不recover
			xrlog.GetInstance().Infof(xrconstant.GoroutineDone)
		}()
		p.BusChannelWaitGroup.Add(1)
		defer p.BusChannelWaitGroup.Done()

		HandleEvent(p.BusChannel, p.Opt.DefaultHandler)
	}()

	// 是否开启http采集分析
	if 0 < p.Bench.Json.Base.PprofHttpPort {
		xrpporf.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", p.Bench.Json.Base.PprofHttpPort))
	}

	// 全局定时器
	if p.Bench.Json.Timer.ScanSecondDuration != nil || p.Bench.Json.Timer.ScanMillisecondDuration != nil {
		err = p.Timer.Start(context.TODO(),
			xrtimer.NewOptions().
				SetScanSecondDuration(p.Bench.Json.Timer.ScanSecondDuration).
				SetScanMillisecondDuration(p.Bench.Json.Timer.ScanMillisecondDuration).
				SetTimerOutChan(p.BusChannel),
		)
		if err != nil {
			return errors.Errorf("timer Start err:%v %v ", err, xrutil.GetCodeLocation(1).String())
		}
	}

	runtime.GC()
	return nil
}

func (p *Mgr) PostInit(ctx context.Context, opts ...*Options) error {
	var err error

	p.Opt = mergeOptions(opts...)
	err = configure(p.Opt)
	if err != nil {
		return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
	}

	// 启动Etcd
	p.Bench.Json.Etcd.Key = fmt.Sprintf("/%v/%v/%v/%v/%v",
		common.ProjectName, common.EtcdWatchMsgTypeService, p.ZoneID, p.ServiceName, p.ServiceID)
	err = etcd.Start(&p.Bench.Json.Etcd, p.BusChannel, p.Opt.EtcdHandler)
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

	runtime.GC()
	return nil
}

func (p *Mgr) Stop() error {
	// 定时检查事件总线是否消费完成
	go checkGBusChannel()

	// 等待GEventChan处理结束
	p.BusChannelWaitGroup.Wait()

	if p.Bench.Json.Timer.ScanSecondDuration != nil || p.Bench.Json.Timer.ScanMillisecondDuration != nil {
		p.Timer.Stop()
		xrlog.GetInstance().Warn("GTimer stop")
	}
	if xretcd.IsEnable() {
		_ = xretcd.GetInstance().Stop()
		xrlog.GetInstance().Warn("GEtcd stop")
	}

	return nil
}

func checkGBusChannel() {
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
			GBusChannelCheckChan <- struct{}{}
			xrlog.GetInstance().Warn("send to GBusChannelCheckChan")
		}
	}
}
