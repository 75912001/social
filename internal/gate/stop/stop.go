package stop

import (
	xrlog "social/pkg/lib/log"
	"social/pkg/server"
)

// PreStop 关闭前处理
func PreStop() {
	// todo menglingchao 关机前处理...

	// todo menglingchao 关闭grpc服务 拒绝新连接
	xrlog.GetInstance().Warn("grpc Service stop")

	// 设置为关闭中
	server.GetInstance().SetStopping()
}
