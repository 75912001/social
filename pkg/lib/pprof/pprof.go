package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	libconstant "social/pkg/lib/constant"
	liblog "social/pkg/lib/log"
	libutil "social/pkg/lib/util"
)

// StartHTTPprof 开启http采集分析
//
//	参数:
//		addr: "0.0.0.0:8090"
func StartHTTPprof(addr string) {
	go func() {
		defer func() {
			if libutil.IsRelease() {
				if err := recover(); err != nil {
					liblog.PrintErr(libconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			liblog.PrintErr(libconstant.GoroutineDone)
		}()
		if err := http.ListenAndServe(addr, nil); err != nil {
			liblog.PrintErr(libconstant.Failure, addr, err)
		}
	}()
}
