package subbench

import (
	"encoding/json"
	"github.com/pkg/errors"
	xrerror "social/pkg/lib/error"
	xrutil "social/pkg/lib/util"
	"sync"
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
	Gate Gate `json:"gate"`
}

type Gate struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

func (p *mgr) Load(strJson string) error {
	if err := json.Unmarshal([]byte(strJson), &p); err != nil {
		return errors.WithMessagef(err, "%v %v", strJson, xrutil.GetCodeLocation(1).String())
	}
	//base
	if len(p.Gate.IP) == 0 {
		return errors.WithMessage(xrerror.Config, xrutil.GetCodeLocation(1).String())
	}
	if 0 == p.Gate.Port {
		return errors.WithMessage(xrerror.Config, xrutil.GetCodeLocation(1).String())
	}
	return nil
}
