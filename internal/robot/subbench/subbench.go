package subbench

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	liberror "social/lib/error"
	libutil "social/lib/util"
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

func (p *mgr) Load(pathFile string) error {
	if data, err := os.ReadFile(pathFile); err != nil {
		return errors.WithMessagef(err, "%v %v", pathFile, libutil.GetCodeLocation(1).String())
	} else {
		if err = json.Unmarshal(data, &p); err != nil {
			return errors.WithMessagef(err, "%v %v", pathFile, libutil.GetCodeLocation(1).String())
		}
	}
	//base
	if len(p.Gate.IP) == 0 {
		return errors.WithMessage(liberror.Config, libutil.GetCodeLocation(1).String())
	}
	if 0 == p.Gate.Port {
		return errors.WithMessage(liberror.Config, libutil.GetCodeLocation(1).String())
	}
	return nil
}
