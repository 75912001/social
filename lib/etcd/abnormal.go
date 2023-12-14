package etcd

import (
	"context"
	"github.com/pkg/errors"
	"runtime/debug"
	libconsts "social/lib/consts"
	liblog "social/lib/log"
	libruntime "social/lib/runtime"
	libutil "social/lib/util"
	"time"
)

func (p *Mgr) abnormal(ctx context.Context) {
	liblog.PrintErr("etcd lease OnRun abnormal, retrying")
	go func(ctx context.Context) {
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					liblog.PrintErr(libconsts.Retry, libconsts.GoroutinePanic, err, debug.Stack())
				}
			}
			liblog.PrintInfo(libconsts.Retry, libconsts.GoroutineDone)
		}()
		if err := p.Stop(); err != nil {
			liblog.PrintInfo(libconsts.Retry, libconsts.Failure, err)
			return
		}
		if err := p.retryStartAndRun(ctx); err != nil {
			liblog.PrintErr(libconsts.Retry, libconsts.Failure, err)
			return
		}
	}(ctx)
}

// 多次重试 Start 和 Run
func (p *Mgr) retryStartAndRun(ctx context.Context) error {
	liblog.PrintfErr("renewing etcd lease, reconfiguring.grantLeaseMaxRetries:%v, grantLeaseIntervalSecond:%v",
		*p.options.grantLeaseMaxRetries, grantLeaseRetryDuration/time.Second)
	var failedGrantLeaseAttempts = 0
	for {
		if err := p.Start(ctx, p.options); err != nil {
			failedGrantLeaseAttempts++
			if *p.options.grantLeaseMaxRetries <= failedGrantLeaseAttempts {
				return errors.WithMessagef(err, "%v exceeded max attempts to renew etcd lease %v %v",
					libruntime.GetCodeLocation(1), *p.options.grantLeaseMaxRetries, failedGrantLeaseAttempts)
			}
			liblog.PrintErr("error granting etcd lease, will retry.", err)
			time.Sleep(grantLeaseRetryDuration)
			continue
		} else {
		retryKeepAlive:
			// 续租
			if err = p.Run(ctx); err != nil {
				failedGrantLeaseAttempts++
				if *p.options.grantLeaseMaxRetries <= failedGrantLeaseAttempts {
					return errors.WithMessagef(err, "%v exceeded max attempts to renew etcd lease %v %v",
						libruntime.GetCodeLocation(1), *p.options.grantLeaseMaxRetries, failedGrantLeaseAttempts)
				}
				liblog.PrintErr("error granting etcd lease, will retry.", err)
				time.Sleep(grantLeaseRetryDuration)
				goto retryKeepAlive
			} else {
				return nil
			}
		}
	}
}
