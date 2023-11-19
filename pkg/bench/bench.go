package bench

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"runtime"
	"social/pkg/common"
	xrerror "social/pkg/lib/error"
	xrlog "social/pkg/lib/log"
	xrutil "social/pkg/lib/util"
	"social/pkg/server"
	"time"
)

//bench.json 配置文件.
//该配置文件与可执行程序在同一目录下.

type Mgr struct {
	Json benchJson
}

type benchJson struct {
	Base   Base   `json:"base"`
	Etcd   Etcd   `json:"etcd"`
	Timer  Timer  `json:"timer"`
	Server Server `json:"server"`
}

type Base struct {
	Version       string         `json:"version"`
	PprofHttpPort uint32         `json:"pprofHttpPort"` //pprof性能分析 http端口 default:0 不使用
	LogLevel      int            `json:"logLevel"`      //日志等级 default:7
	LogAbsPath    string         `json:"logAbsPath"`    //日志绝对路径 default:/data/xxx/log
	GoMaxProcess  int            `json:"goMaxProcess"`  //default:runtime.NumCPU()
	RunMode       xrutil.RunMode `json:"runMode"`       //运行模式 0:release 1:debug default:0,release
}

type Server struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type Timer struct {
	//秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为100毫秒 default:100000000
	ScanSecondDuration *time.Duration `json:"scanSecondDuration"`
	//毫秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为25毫秒 default:25000000
	ScanMillisecondDuration *time.Duration `json:"scanMillisecondDuration"`
}

type Etcd struct {
	Addrs []string      `json:"addrs"`
	TTL   int64         `json:"ttl"` //ttl 秒 [默认为 common.EtcdTtlSecondDefault 秒, e.g.:系统每10秒续约一次,该参数至少为11秒]
	Key   string        //`json:"key"`   //common.ProjectName/common.EtcdWatchMsgTypeService/zoneID/serviceName/serviceID
	Value EtcdValueJson `json:"value"` //有:直接使用. default:{"ip":"192.168.50.10","port":3021, "version":version}
}

type EtcdValueJson struct {
	ServiceNetTCP ServiceNetJson `json:"serviceNetTCP,omitempty"` //有:直接使用. 没有:使用 server 属性生成ip, port
	Version       string         `json:"version,omitempty"`       //有:直接使用. 没有:使用 base.version 生成
	AvailableLoad uint32         `json:"availableLoad,omitempty"` //可用负载, 可用资源数 [默认为 common.BenchJsonAvailableLoadMaxDefault]
}

// ServiceNetJson 服务 网络 接口
type ServiceNetJson struct {
	IP   string `json:"ip,omitempty"`
	Port uint16 `json:"port,omitempty"`
}

// String 显示服务信息
func (p *benchJson) String() string {
	return fmt.Sprintf("version:%v", p.Base.Version)
}

// Parse 解析, bench.json
func (p *Mgr) Parse(pathFile string) error {
	if err := json.Unmarshal([]byte(pathFile), &p.Json); err != nil {
		return errors.WithMessagef(err, "%v %v", pathFile, xrutil.GetCodeLocation(1).String())
	}
	//base
	if len(p.Json.Base.Version) == 0 {
		return errors.WithMessage(xrerror.Config, xrutil.GetCodeLocation(1).String())
	}
	if 0 == p.Json.Base.LogLevel {
		p.Json.Base.LogLevel = int(xrlog.LevelOn)
	}
	if len(p.Json.Base.LogAbsPath) == 0 {
		p.Json.Base.LogAbsPath = common.LogAbsPath
	}
	if 0 == p.Json.Base.GoMaxProcess {
		p.Json.Base.GoMaxProcess = runtime.NumCPU()
	}
	xrutil.GRunMode = p.Json.Base.RunMode
	//server
	if len(p.Json.Server.IP) == 0 {
		return errors.WithMessage(xrerror.Config, xrutil.GetCodeLocation(1).String())
	}
	if p.Json.Server.Port == 0 {
		return errors.WithMessage(xrerror.Config, xrutil.GetCodeLocation(1).String())
	}
	//etcd
	if len(p.Json.Etcd.Addrs) == 0 {
		return errors.WithMessage(xrerror.Config, xrutil.GetCodeLocation(1).String())
	}
	if p.Json.Etcd.TTL == 0 {
		p.Json.Etcd.TTL = common.EtcdTtlSecondDefault
	}
	if len(p.Json.Etcd.Key) == 0 {
		p.Json.Etcd.Key = fmt.Sprintf("%v/%v/%v/%v/%v",
			common.ProjectName, common.EtcdWatchMsgTypeService,
			server.GMgr.ZoneID, server.GMgr.ServiceName, server.GMgr.ServiceID)
	}
	if len(p.Json.Etcd.Value.ServiceNetTCP.IP) == 0 {
		p.Json.Etcd.Value.ServiceNetTCP.IP = p.Json.Server.IP
	}
	if p.Json.Etcd.Value.ServiceNetTCP.Port == 0 {
		p.Json.Etcd.Value.ServiceNetTCP.Port = p.Json.Server.Port
	}
	//timer
	if nil == p.Json.Timer.ScanSecondDuration {
		t := time.Millisecond * 100
		p.Json.Timer.ScanSecondDuration = &t
	}
	if nil == p.Json.Timer.ScanMillisecondDuration {
		t := time.Millisecond * 25
		p.Json.Timer.ScanMillisecondDuration = &t
	}
	return nil
}
