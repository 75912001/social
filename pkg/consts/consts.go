package consts

import "time"

const (
	ProjectName                         = "social"
	LogAbsPath                          = "/data/" + ProjectName + "/log"
	BusChannelNumberDefault             = 1000000 //1000000 大约占用15.6MB
	TimerScanSecondDurationDefault      = time.Millisecond * 100
	TimerScanMillisecondDurationDefault = time.Millisecond * 25
)
