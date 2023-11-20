package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	xrconstant "social/pkg/lib/constant"
	xrlog "social/pkg/lib/log"
	xrutil "social/pkg/lib/util"
)

// StartHTTPprof 开启http采集分析
//
//	参数:
//		addr: "0.0.0.0:8090"
func StartHTTPprof(addr string) {
	go func() {
		defer func() {
			if xrutil.IsRelease() {
				if err := recover(); err != nil {
					xrlog.PrintErr(xrconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			xrlog.PrintErr(xrconstant.GoroutineDone)
		}()
		if err := http.ListenAndServe(addr, nil); err != nil {
			xrlog.PrintErr(xrconstant.Failure, addr, err)
		}
	}()
}
