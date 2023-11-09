package log

type Level int

// 日志等级
const (
	LevelOff   Level = 0 //关闭
	LevelFatal Level = 1
	LevelError Level = 2
	LevelWarn  Level = 3
	LevelInfo  Level = 4
	LevelDebug Level = 5
	LevelTrace Level = 6
	LevelOn    Level = 7 //7 全部打开
)

var levelTag = []string{
	LevelOff:   "LevelOff",
	LevelFatal: "Fatal",
	LevelError: "Error",
	LevelWarn:  "Warn",
	LevelInfo:  "Info",
	LevelDebug: "Debug",
	LevelTrace: "Trace",
	LevelOn:    "LevelOn",
}
