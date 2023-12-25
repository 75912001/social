package friend

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	libmongodb "social/lib/mongodb"
	libruntime "social/lib/runtime"
	pkgconsts "social/pkg/consts"
	pkgsubbench "social/pkg/subbench"
)

type subBenchMgr struct {
	ZoneMongoDB pkgsubbench.MongoDB `json:"zoneMongoDB"`
}

// Parse 加载 配置文件
func (p *subBenchMgr) Parse(pathFile string) error {
	if data, err := os.ReadFile(pathFile); err != nil {
		return errors.WithMessagef(err, "%v %v", pathFile, libruntime.Location())
	} else {
		if err = json.Unmarshal(data, &p); err != nil {
			return errors.WithMessagef(err, "%v %v", pathFile, libruntime.Location())
		}
	}
	//check
	if len(p.ZoneMongoDB.Addrs) == 0 {
		return errors.Errorf("mongodb address is error %v", p.ZoneMongoDB)
	}
	if len(p.ZoneMongoDB.User) == 0 || len(p.ZoneMongoDB.Password) == 0 {
		return errors.Errorf("mongodb user or password is error %v", p.ZoneMongoDB)
	}
	//带默认值
	if p.ZoneMongoDB.DBName == nil {
		v := pkgconsts.ProjectName
		p.ZoneMongoDB.DBName = &v
	}
	if p.ZoneMongoDB.MaxPoolSize == nil {
		p.ZoneMongoDB.MaxPoolSize = &libmongodb.MaxPoolSizeDefault
	}
	if p.ZoneMongoDB.MinPoolSize == nil {
		p.ZoneMongoDB.MinPoolSize = &libmongodb.MinPoolSizeDefault
	}
	if p.ZoneMongoDB.TimeoutDuration == nil {
		p.ZoneMongoDB.TimeoutDuration = &libmongodb.TimeoutDurationDefault
	}
	if p.ZoneMongoDB.MaxConnIdleTime == nil {
		p.ZoneMongoDB.MaxConnIdleTime = &libmongodb.MaxConnIdleTimeDefault
	}
	if p.ZoneMongoDB.MaxConnecting == nil {
		p.ZoneMongoDB.MaxConnecting = &libmongodb.MaxConnectingDefault
	}
	return nil
}
