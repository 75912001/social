package bench

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"os"
	"runtime"
	pkgcommon "social/pkg/common"
	liberror "social/pkg/lib/error"
	liblog "social/pkg/lib/log"
	libutil "social/pkg/lib/util"
	"sync"
	"time"
)

//bench.json 配置文件.
//该配置文件与可执行程序在同一目录下.

var (
	instance *Mgr
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *Mgr {
	once.Do(func() {
		instance = new(Mgr)
	})
	return instance
}

type Mgr struct {
	Base   Base   `json:"base"`
	Etcd   Etcd   `json:"etcd"`
	Timer  Timer  `json:"timer"`
	Server Server `json:"server"`
}

type Base struct {
	Version          string          `json:"version"`
	PprofHttpPort    uint32          `json:"pprofHttpPort"`    //pprof性能分析 http端口 default:0 不使用
	LogLevel         int             `json:"logLevel"`         //日志等级 default:7
	LogAbsPath       string          `json:"logAbsPath"`       //日志绝对路径 default:/data/xxx/log
	GoMaxProcess     int             `json:"goMaxProcess"`     //default:runtime.NumCPU()
	AvailableLoad    uint32          `json:"availableLoad"`    //可用负载, 可用资源数 default:math.MaxUint32
	BusChannelNumber uint32          `json:"busChannelNumber"` //事件chan数量. default:1000000 大约占用15.6MB
	RunMode          libutil.RunMode `json:"runMode"`          //运行模式 0:release 1:debug default:0,release
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
	TTL   int64         `json:"ttl"`   // ttl 秒 [默认为 common.EtcdTtlSecondDefault 秒, e.g.:系统每10秒续约一次,该参数至少为11秒]
	Key   string        `json:"key"`   // common.ProjectName/common.EtcdWatchMsgTypeService/zoneID/serviceName/serviceID
	Value EtcdValueJson `json:"value"` // 有:直接使用. default:{"serviceNetTCP":{"ip":"0.0.0.0","port":5101},"version":"Beta.0.0.1","availableLoad":4294967295}
}

type EtcdValueJson struct {
	ServiceNetTCP ServiceNetJson `json:"serviceNetTCP,omitempty"` //有:直接使用. 没有:使用 server 属性生成ip, port
	Version       string         `json:"version,omitempty"`       //有:直接使用. 没有:使用 base.version 生成
	AvailableLoad uint32         `json:"availableLoad,omitempty"` //可用负载, 可用资源数
}

// ServiceNetJson 服务 网络 接口
type ServiceNetJson struct {
	IP   string `json:"ip,omitempty"`
	Port uint16 `json:"port,omitempty"`
}

// String 显示服务信息
func (p *Mgr) String() string {
	return fmt.Sprintf("version:%v", p.Base.Version)
}

// Parse 解析, bench.json
func (p *Mgr) Parse(pathFile string, zoneID uint32, serviceName string, serviceID uint32) error {
	if data, err := os.ReadFile(pathFile); err != nil {
		return errors.WithMessagef(err, "%v %v", pathFile, libutil.GetCodeLocation(1).String())
	} else {
		if err = json.Unmarshal(data, &p); err != nil {
			return errors.WithMessagef(err, "%v %v", pathFile, libutil.GetCodeLocation(1).String())
		}
	}
	//base
	if len(p.Base.Version) == 0 {
		return errors.WithMessage(liberror.Config, libutil.GetCodeLocation(1).String())
	}
	if 0 == p.Base.LogLevel {
		p.Base.LogLevel = int(liblog.LevelOn)
	}
	if len(p.Base.LogAbsPath) == 0 {
		p.Base.LogAbsPath = pkgcommon.LogAbsPath
	}
	if 0 == p.Base.GoMaxProcess {
		p.Base.GoMaxProcess = runtime.NumCPU()
	}
	if 0 == p.Base.AvailableLoad {
		p.Base.AvailableLoad = math.MaxUint32
	}
	if 0 == p.Base.BusChannelNumber {
		//1000000 大约占用15.6MB
		p.Base.BusChannelNumber = 1000000
	}
	libutil.GRunMode = p.Base.RunMode
	//server
	if len(p.Server.IP) == 0 {
		return errors.WithMessage(liberror.Config, libutil.GetCodeLocation(1).String())
	}
	if p.Server.Port == 0 {
		return errors.WithMessage(liberror.Config, libutil.GetCodeLocation(1).String())
	}
	//etcd
	if len(p.Etcd.Addrs) == 0 {
		return errors.WithMessage(liberror.Config, libutil.GetCodeLocation(1).String())
	}
	if p.Etcd.TTL == 0 {
		p.Etcd.TTL = pkgcommon.EtcdTtlSecondDefault
	}
	if len(p.Etcd.Key) == 0 {
		p.Etcd.Key = fmt.Sprintf("%v/%v/%v/%v/%v",
			pkgcommon.ProjectName, pkgcommon.EtcdWatchMsgTypeService,
			zoneID, serviceName, serviceID)
	}
	if len(p.Etcd.Value.ServiceNetTCP.IP) == 0 {
		p.Etcd.Value.ServiceNetTCP.IP = p.Server.IP
	}
	if p.Etcd.Value.ServiceNetTCP.Port == 0 {
		p.Etcd.Value.ServiceNetTCP.Port = p.Server.Port
	}
	if len(p.Etcd.Value.Version) == 0 {
		p.Etcd.Value.Version = p.Base.Version
	}
	//timer
	if nil == p.Timer.ScanSecondDuration {
		t := time.Millisecond * 100
		p.Timer.ScanSecondDuration = &t
	}
	if nil == p.Timer.ScanMillisecondDuration {
		t := time.Millisecond * 25
		p.Timer.ScanMillisecondDuration = &t
	}
	return nil
}
