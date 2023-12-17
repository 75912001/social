package robot

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
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
		return errors.WithMessagef(err, "%v %v", pathFile, libruntime.Location())
	} else {
		if err = json.Unmarshal(data, &p); err != nil {
			return errors.WithMessagef(err, "%v %v", pathFile, libruntime.Location())
		}
	}
	//base
	if len(p.Gate.IP) == 0 {
		return errors.WithMessage(liberror.Config, libruntime.Location())
	}
	if 0 == p.Gate.Port {
		return errors.WithMessage(liberror.Config, libruntime.Location())
	}
	return nil
}
