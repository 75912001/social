package mongodb

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	"strings"
	"time"
)

// Options contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Options struct {
	addrs    []string //192.168.50.98:27017 slice
	userName *string  //用户名称
	pw       *string  //密码
	dbName   *string  //数据库名

	// 连接池最大数量,该数量应该与并发数量匹配.e.g.:有100个并发协程,则需要100个连接数量.If this is 0, maximum connection pool size is not limited. The default is 100.
	// specifies that maximum number of connections allowed in the driver's connection pool to each server.
	// Requests to a server will block if this maximum is reached. This can also be set through the "maxPoolSize" URI option
	// (e.g. "maxPoolSize=100"). If this is 0, maximum connection pool size is not limited. The default is 100.
	// 观察正常压力周期下最大连接数；
	// 将maxPoolSize设置为比最大连接数稍大的值(例如大30%)
	maxPoolSize *uint64
	// specifies the minimum number of connections allowed in the driver's connection pool to each server. If
	// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
	// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
	minPoolSize *uint64
	// 操作超时时间, 1000*1000*1000 = 1 Second = time.Second  e.g.:time.Second * 10
	timeoutDuration *time.Duration //default: timeoutDurationDefault
	// specifies the maximum amount of time that a connection will remain idle in a connection pool
	// before it is removed from the pool and closed. This can also be set through the "maxIdleTimeMS" URI option (e.g.
	// "maxIdleTimeMS=10000"). The default is 0, meaning a connection can remain unused indefinitely.
	maxConnIdleTime *time.Duration
	// specifies the maximum number of connections a connection pool may establish simultaneously. This can
	// also be set through the "maxConnecting" URI option (e.g. "maxConnecting=2"). If this is 0, the default is used. The
	// default is 2. Values greater than 100 are not recommended.
	maxConnecting *uint64
}

// NewOptions 新的Options
func NewOptions() *Options {
	return new(Options)
}

func (p *Options) SetAddrs(addrs []string) *Options {
	p.addrs = append(p.addrs, addrs...)
	return p
}

func (p *Options) SetUserName(userName string) *Options {
	p.userName = &userName
	return p
}

func (p *Options) SetPW(pw string) *Options {
	p.pw = &pw
	return p
}

func (p *Options) SetDBName(dbName *string) *Options {
	p.dbName = dbName
	return p
}

func (p *Options) SetMinPoolSize(minPoolSize *uint64) *Options {
	p.minPoolSize = minPoolSize
	return p
}

func (p *Options) SetMaxPoolSize(maxPoolSize *uint64) *Options {
	p.maxPoolSize = maxPoolSize
	return p
}

func (p *Options) SetTimeoutDuration(timeoutDuration *time.Duration) *Options {
	p.timeoutDuration = timeoutDuration
	return p
}

func (p *Options) SetMaxConnIdleTime(maxConnIdleTime *time.Duration) *Options {
	p.maxConnIdleTime = maxConnIdleTime
	return p
}

func (p *Options) SetMaxConnecting(maxConnecting *uint64) *Options {
	p.maxConnecting = maxConnecting
	return p
}

func (p *Options) genURI() string {
	var uri string
	userName := url.QueryEscape(*p.userName)
	pw := url.QueryEscape(*p.pw)
	//副本集+切片集+route
	var hostPort string
	for _, v := range p.addrs {
		hostPort += v + ","
	}
	hostPort = strings.TrimRight(hostPort, ",")
	//uri = fmt.Sprintf("mongodb://%v:%v@%v/%v?connect=automatic&replicaSet=replset", userName, pw, hostPort, dbName)
	uri = fmt.Sprintf("mongodb://%v:%v@%v/%v?connect=automatic", userName, pw, hostPort, *p.dbName)
	//直连模式
	//uri = fmt.Sprintf("mongodb://%v:%v@%v:%v/%v%v", userName, pw, ip, port, dbName, "?connect=direct")
	return uri
}

// mergeOptions combines the given *Options into a single *Options in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func mergeOptions(opts ...*Options) *Options {
	newOptions := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.addrs != nil {
			newOptions.SetAddrs(opt.addrs)
		}
		if opt.userName != nil {
			newOptions.userName = opt.userName
		}
		if opt.pw != nil {
			newOptions.pw = opt.pw
		}
		if opt.dbName != nil {
			newOptions.dbName = opt.dbName
		}
		if opt.maxPoolSize != nil {
			newOptions.maxPoolSize = opt.maxPoolSize
		}
		if opt.minPoolSize != nil {
			newOptions.minPoolSize = opt.minPoolSize
		}
		if opt.timeoutDuration != nil {
			newOptions.timeoutDuration = opt.timeoutDuration
		}
		if opt.maxConnIdleTime != nil {
			newOptions.maxConnIdleTime = opt.maxConnIdleTime
		}
		if opt.maxConnecting != nil {
			newOptions.maxConnecting = opt.maxConnecting
		}
	}
	return newOptions
}

// 配置
func configure(opts *Options) error {
	if len(opts.addrs) == 0 {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	if opts.userName == nil {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	if opts.pw == nil {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	if opts.dbName == nil {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	if opts.maxPoolSize == nil {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	if opts.minPoolSize == nil {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	if opts.timeoutDuration == nil {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	if opts.maxConnIdleTime == nil {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	if opts.maxConnecting == nil {
		return errors.WithMessagef(liberror.Param, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	return nil
}
