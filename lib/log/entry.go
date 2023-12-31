package log

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	libconsts "social/lib/consts"
	libtime "social/lib/time"
	"strconv"
	"time"
)

//日志条目

// fields 日志数据字段
type fields []interface{}

// Entry 日志数据信息
type Entry struct {
	logger     *Mgr
	time       time.Time //生成日志的时间
	level      Level     //日志级别
	callerInfo string    //调用堆栈信息
	message    string    //日志消息
	ctx        context.Context
	fields     fields //key,value
}

func (p *Entry) reset() {
	p.logger = nil
	p.fields = nil
	p.level = LevelOff
	p.callerInfo = ""
	p.message = ""
	p.ctx = nil
}

// newEntry 创建
func newEntry(logger *Mgr) *Entry {
	if logger.options.IsEnablePool() {
		e := logger.pool.Get().(*Entry)
		e.logger = logger
		return e
	} else {
		return &Entry{
			logger: logger,
		}
	}
}

// WithContext 由ctx创建Entry
func (p *Entry) WithContext(ctx context.Context) *Entry {
	p.ctx = ctx
	return p
}

// WithField 由field创建Entry
func (p *Entry) WithField(key string, value interface{}) *Entry {
	p.fields = append(p.fields, key, value)
	return p
}

// WithFields 由多个field创建Entry
func (p *Entry) WithFields(fields fields) *Entry {
	p.fields = append(p.fields, fields...)
	return p
}

// formatMessage 格式化日志信息
func (p *Entry) formatMessage() string {
	// 格式为  [时间][日志级别][UID:xxx][堆栈信息][TraceID:xxx]{fields-key:fields:val}...{}自定义内容
	var buf bytes.Buffer
	buf.Grow(bufferCapacity)

	// 时间
	buf.WriteString(fmt.Sprint("[", p.time.Format(logTimeFormat), "]"))

	// 日志级别
	buf.WriteString(fmt.Sprint("[", levelName[p.level], "]"))

	// UID 优先从ctx查找 其次查找field 当UID不存在时也需要占位 值为0
	var uid uint64
	if p.ctx != nil {
		uidVal := p.ctx.Value(UserIDKey)
		if uidVal != nil {
			uid, _ = uidVal.(uint64)
		}
	}
	if 0 == uid { //没有找到UID,从field中查找
		for _, v := range p.fields {
			str, ok := v.(string)
			if ok && str == UserIDKey { //找到
				uid, _ = v.(uint64)
				break
			}
		}
	}
	buf.WriteString(fmt.Sprint("[", UserIDKey, ":", strconv.FormatUint(uid, 10), "]"))

	// 堆栈
	if len(p.callerInfo) > 0 {
		buf.WriteString(fmt.Sprint("[", p.callerInfo, "]"))
	}

	// 处理ctx TraceID
	if p.ctx != nil {
		traceIdVal := p.ctx.Value(TraceIDKey)
		if traceIdVal != nil {
			buf.WriteString(fmt.Sprint("[", TraceIDKey, ":", traceIdVal.(string), "]"))
		}
	}

	// 处理fields
	for idx, v := range p.fields {
		if idx%2 == 0 { //key
			buf.WriteString("{")
			str, ok := v.(string)
			if ok {
				buf.WriteString(str)
			} else {
				buf.WriteString(fmt.Sprint(v))
			}
			buf.WriteString(":")
		} else { //val
			str, ok := v.(string)
			if ok {
				buf.WriteString(str)
			} else {
				buf.WriteString(fmt.Sprint(v))
			}
			buf.WriteString("}")
		}
	}

	// 自定义内容
	buf.WriteString(p.message)

	return buf.String()
}

// log 记录日志
func (p *Entry) log(level Level, skip int, v ...interface{}) {
	p.level = level
	p.time = libtime.NowTime()

	if *p.logger.options.isReportCaller {
		pc, _, line, ok := runtime.Caller(skip)
		funcName := libconsts.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		p.callerInfo = fmt.Sprintf(callerInfoFormat, line, funcName)
	}
	p.message = fmt.Sprintln(v...)

	p.logger.logChan <- p
}

// log 记录日志
func (p *Entry) logf(level Level, skip int, format string, v ...interface{}) {
	p.level = level
	p.time = libtime.NowTime()

	if *p.logger.options.isReportCaller {
		pc, _, line, ok := runtime.Caller(skip)
		funcName := libconsts.Unknown
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		p.callerInfo = fmt.Sprintf(callerInfoFormat, line, funcName)
	}
	p.message = fmt.Sprintf(format, v...)

	p.logger.logChan <- p
}

// Trace 追踪日志
func (p *Entry) Trace(v ...interface{}) {
	if *p.logger.options.level < LevelTrace {
		return
	}
	p.log(LevelTrace, calldepth2, v...)
}

// Tracef 追踪日志
func (p *Entry) Tracef(format string, v ...interface{}) {
	if *p.logger.options.level < LevelTrace {
		return
	}
	p.logf(LevelTrace, calldepth2, format, v...)
}

// Debug 调试日志
func (p *Entry) Debug(v ...interface{}) {
	if *p.logger.options.level < LevelDebug {
		return
	}
	p.log(LevelDebug, calldepth2, v...)
}

// Debugf 调试日志
func (p *Entry) Debugf(format string, v ...interface{}) {
	if *p.logger.options.level < LevelDebug {
		return
	}
	p.logf(LevelDebug, calldepth2, format, v...)
}

// Info 信息日志
func (p *Entry) Info(v ...interface{}) {
	if *p.logger.options.level < LevelInfo {
		return
	}
	p.log(LevelInfo, calldepth2, v...)
}

// Infof 信息日志
func (p *Entry) Infof(format string, v ...interface{}) {
	if *p.logger.options.level < LevelInfo {
		return
	}
	p.logf(LevelInfo, calldepth2, format, v...)
}

// Warn 警告日志
func (p *Entry) Warn(v ...interface{}) {
	if *p.logger.options.level < LevelWarn {
		return
	}
	p.log(LevelWarn, calldepth2, v...)
}

// Warnf 警告日志
func (p *Entry) Warnf(format string, v ...interface{}) {
	if *p.logger.options.level < LevelWarn {
		return
	}
	p.logf(LevelWarn, calldepth2, format, v...)
}

// Error 错误日志
func (p *Entry) Error(v ...interface{}) {
	if *p.logger.options.level < LevelError {
		return
	}
	p.log(LevelError, calldepth2, v...)
}

// Errorf 错误日志
func (p *Entry) Errorf(format string, v ...interface{}) {
	if *p.logger.options.level < LevelError {
		return
	}
	p.logf(LevelError, calldepth2, format, v...)
}

// Fatal 致命日志
func (p *Entry) Fatal(v ...interface{}) {
	if *p.logger.options.level < LevelFatal {
		return
	}
	p.log(LevelFatal, calldepth2, v...)
}

// Fatalf 致命日志
func (p *Entry) Fatalf(format string, v ...interface{}) {
	if *p.logger.options.level < LevelFatal {
		return
	}
	p.logf(LevelFatal, calldepth2, format, v...)
}
