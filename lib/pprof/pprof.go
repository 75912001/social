package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	libconsts "social/lib/consts"
	liblog "social/lib/log"
	libutil "social/lib/util"
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
					liblog.PrintErr(libconsts.GoroutinePanic, err, debug.Stack())
				}
			}
			liblog.PrintErr(libconsts.GoroutineDone)
		}()
		if err := http.ListenAndServe(addr, nil); err != nil {
			liblog.PrintErr(libconsts.Failure, addr, err)
		}
	}()
}
