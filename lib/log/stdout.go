package log

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	libconsts "social/lib/consts"
	libtime "social/lib/time"
)

var stdOut = log.New(os.Stdout, "", 0)

// PrintInfo 输出到os.Stdout
func PrintInfo(v ...interface{}) {
	if GetInstance().GetLevel() < LevelInfo {
		return
	}
	if IsEnable() { // 日志已启用,需要放入日志 channel 中
		GetInstance().log(LevelInfo, v...)
	} else {
		pc, _, line, ok := runtime.Caller(calldepth1)
		funcName := libconsts.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		var buf bytes.Buffer
		buf.Grow(bufferCapacity)
		// 格式为  [时间][日志级别][UID:xxx][堆栈信息]自定义内容
		buf.WriteString(fmt.Sprint("[", libtime.NowTime().Format(logTimeFormat), "]"))
		buf.WriteString(fmt.Sprint("[", levelName[LevelInfo], "]"))
		buf.WriteString("[UID:0]")
		buf.WriteString(fmt.Sprint("[", fmt.Sprintf(callerInfoFormat, line, funcName), "]"))
		buf.WriteString(fmt.Sprint(v...))
		_ = stdOut.Output(calldepth2, buf.String())
	}
}

// PrintfInfo 输出到os.Stdout
func PrintfInfo(format string, v ...interface{}) {
	if GetInstance().GetLevel() < LevelInfo {
		return
	}
	if IsEnable() { // 日志已启用,需要放入日志 channel 中
		GetInstance().logf(LevelInfo, format, v...)
	} else {
		pc, _, line, ok := runtime.Caller(calldepth1)
		funcName := libconsts.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		var buf bytes.Buffer
		buf.Grow(bufferCapacity)
		// 格式为  [时间][日志级别][UID:xxx][堆栈信息]自定义内容
		buf.WriteString(fmt.Sprint("[", libtime.NowTime().Format(logTimeFormat), "]"))
		buf.WriteString(fmt.Sprint("[", levelName[LevelInfo], "]"))
		buf.WriteString("[UID:0]")
		buf.WriteString(fmt.Sprint("[", fmt.Sprintf(callerInfoFormat, line, funcName), "]"))
		buf.WriteString(fmt.Sprintf(format, v...))
		_ = stdOut.Output(calldepth2, buf.String())
	}
}
