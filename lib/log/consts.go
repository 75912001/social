package log

const (
	logChannelCapacity = 100000            // 日志通道最大容量
	logTimeFormat      = "15:04:05.000000" // 日志时间格式 时:分:秒.微秒
	callerInfoFormat   = "Line:%d %s"      // 堆栈信息格式 例如 Line:69 server/xxx/xx/x/log.TestExample
	TraceIDKey         = "TraceID"         // 日志traceId key
	UserIDKey          = "UID"             // 日志用户ID key
	bufferCapacity     = 300               // buffer 容量
	calldepth1         = 1
	calldepth2         = calldepth1 + 1 //打印堆栈时候跳过的层数
	calldepth3         = calldepth2 + 1
)

const (
	accessLogFileBaseName = "access.log"
	errorLogFileBaseName  = "error.log"
)
