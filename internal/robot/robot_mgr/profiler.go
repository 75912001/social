package robot_mgr

import (
	xrlog "dawn-server/impl/xr/lib/log"
	"fmt"
	"github.com/VividCortex/gohistogram"
	"runtime"
)

var GProfiler = NewProfiler()

type Profiler struct {
	CurrentUserNum      int // 当前在线用户数
	TotalUserNum        int // 本次测试用户总数
	currentUserTotalNum int // 当前用户总数（包含下线过的用户）
	connectFailedNum    int // 连接失败次数
	sendNum             int // 发送请求总数
	successNum          int // 发送请求成功总数
	errorNum            int // 发送请求异常总数
	cycleSendNum        int // 单位时间内发送请求数
	cycleSuccessSendNum int // 单位时间内发送成功请求数
	cycleResponseNum    int // 单位时间内请求相应数
	costTimeMin         int // 最短耗时-毫秒
	costTimeMax         int // 最长耗时-毫秒

	ChanUserOffline         chan string
	ChanSendNum             chan int
	ChanSuccessNum          chan int
	ChanErrorNum            chan int
	ChanCycleSuccessSendNum chan int
	ChanCycleResponseNum    chan int
	ChanCycleSendNum        chan int
	ChanCostTime            chan int
	ChanConnectFailedNum    chan int

	histogram *gohistogram.NumericHistogram // 统计耗时指标
	RunTicks  int                           // 运行次数
}

func NewProfiler() *Profiler {
	return &Profiler{
		ChanErrorNum:            make(chan int, 100),
		ChanUserOffline:         make(chan string, 100),
		ChanSendNum:             make(chan int, 100),
		ChanSuccessNum:          make(chan int, 100),
		ChanCycleSendNum:        make(chan int, 100),
		ChanCycleSuccessSendNum: make(chan int, 100),
		ChanCycleResponseNum:    make(chan int, 100),
		ChanCostTime:            make(chan int, 100),
		ChanConnectFailedNum:    make(chan int, 100),
		histogram:               gohistogram.NewHistogram(10000), // 设置最大有10000个桶
	}
}

func (p *Profiler) ShowLog() {
	cycleTime := 5
	if p.RunTicks%cycleTime == 0 {
		qps := p.cycleResponseNum / cycleTime
		s := p.cycleSendNum / cycleTime

		p.resetCycleCount()

		message := fmt.Sprintf("[%dm%ds]连接失败:%d__发出:%d(%d/s)__异常响应:%d__正常响应:%d(%d/s)__在线用户数:%d"+
			"__已上线用户数:%d__协程数:%d__最短耗时:%dms__最长耗时:%dms__平均耗时:%dms__P95:%dms__P99:%dms",
			p.RunTicks/60, p.RunTicks%60, p.connectFailedNum, p.sendNum, s, p.errorNum, p.successNum, qps, p.CurrentUserNum,
			p.currentUserTotalNum, runtime.NumGoroutine(), p.costTimeMin, p.costTimeMax,
			int(p.histogram.Mean()), int(p.histogram.Quantile(0.95)), int(p.histogram.Quantile(0.99)))

		xrlog.GetInstance().Info(message)
		fmt.Println(message)
	}
}

func (p *Profiler) resetCycleCount() {
	p.cycleSendNum = 0
	p.cycleSuccessSendNum = 0
	p.cycleResponseNum = 0
}

func (p *Profiler) AddConnectFailedNum(n int) {
	p.connectFailedNum = p.connectFailedNum + n
}

func (p *Profiler) AddCurrentUserNum(n int) {
	p.CurrentUserNum = p.CurrentUserNum + n
}

func (p *Profiler) AddCurrentUserTotalNum(n int) {
	p.currentUserTotalNum = p.currentUserTotalNum + n
}

func (p *Profiler) AddErrorNum(n int) {
	p.errorNum = p.errorNum + n
}

func (p *Profiler) AddSuccessNum(n int) {
	p.successNum = p.successNum + n
}

func (p *Profiler) AddSendNum(n int) {
	p.sendNum = p.sendNum + n
}

func (p *Profiler) AddCycleSendNum(n int) {
	p.cycleSendNum = p.cycleSendNum + n
}

func (p *Profiler) AddCycleSuccessSendNum(n int) {
	p.cycleSuccessSendNum = p.cycleSuccessSendNum + n
}
func (p *Profiler) AddCycleResponseNum(n int) {
	p.cycleResponseNum = p.cycleResponseNum + n
}

func (p *Profiler) ExecuteCostTime(costMs int) {
	xrlog.GetInstance().Debugf("ExecuteCostTime:%d", costMs)

	// 最小精度为1ms
	if 0 == costMs {
		costMs = 1
	}

	if 0 == p.costTimeMin || p.costTimeMin > costMs {
		p.costTimeMin = costMs
	}
	if 0 == p.costTimeMax || p.costTimeMax < costMs {
		p.costTimeMax = costMs
	}

	// 添加耗时信息
	p.histogram.Add(float64(costMs))
}

func (p *Profiler) Watch() {
	for {
		select {
		//用户下线
		case <-p.ChanUserOffline:
			p.AddCurrentUserNum(-1)
		//记录发送请求
		case n := <-p.ChanSendNum:
			p.AddSendNum(n)
		//记录连接失败次数
		case n := <-p.ChanConnectFailedNum:
			p.AddConnectFailedNum(n)
		//记录成功请求
		case n := <-p.ChanSuccessNum:
			p.AddSuccessNum(n)
		//记录Error请求
		case n := <-p.ChanErrorNum:
			p.AddErrorNum(n)
		//记录单位时间内发送成功的请求数
		case n := <-p.ChanCycleSuccessSendNum:
			p.AddCycleSuccessSendNum(n)
		//记录单位时间内响应请求数
		case n := <-p.ChanCycleResponseNum:
			p.AddCycleResponseNum(n)
		//记录单位时间内发送的请求数
		case n := <-p.ChanCycleSendNum:
			p.AddCycleSendNum(n)
		//处理耗时信息
		case costTime := <-p.ChanCostTime:
			p.ExecuteCostTime(costTime)
		}
	}
}
