package log

//使用系统log,自带锁
//使用协程操作io输出日志
//每天自动创建新的日志文件

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"runtime/debug"
	libconstant "social/lib/consts"
	liberror "social/lib/error"
	"social/lib/util"
	"strconv"
	"sync"
	"time"
)

var (
	instance *Mgr
	once     sync.Once
)

// GetInstance 获取
func GetInstance() *Mgr {
	once.Do(func() {
		instance = new(Mgr)
	})
	return instance
}

// IsEnable 是否 开启
func IsEnable() bool {
	if instance == nil {
		return false
	}
	return GetInstance().logChan != nil
}

// Mgr 日志
type Mgr struct {
	options *options

	loggerSlice     [LevelOn]*log.Logger // 日志实例 note:此处非协程安全
	logChan         chan *entry          // 日志写入通道
	waitGroupOutPut sync.WaitGroup       // 同步锁
	logDuration     int                  // 日志分割刻度 按天或者小时  e.g.:20210819或2021081901
	openFiles       []*os.File           // 当前打开的文件
	pool            *sync.Pool
}

// GetLevel 获取日志等级
func (p *Mgr) GetLevel() Level {
	if p.options == nil {
		return LevelOff
	}
	if p.options.level == nil {
		return LevelOn
	}
	return *p.options.level
}

// Start 开始
//
//	参数:
//		absPath:日志绝对路径
//		namePrefix:日志名 前缀
func (p *Mgr) Start(_ context.Context, opts ...*options) error {
	p.options = mergeOptions(opts...)
	if err := configure(p.options); err != nil {
		return errors.WithMessage(err, util.GetCodeLocation(1).String())
	}

	p.logChan = make(chan *entry, logChannelCapacity)

	// 初始化logger
	for i := LevelOff; i < LevelOn; i++ {
		p.loggerSlice[i] = log.New(os.Stdout, "", 0)
	}

	// 初始化各级别的日志输出
	if err := p.newWriters(); err != nil {
		return errors.WithMessage(err, util.GetCodeLocation(1).String())
	}

	if p.options.IsEnablePool() {
		p.pool = &sync.Pool{
			New: func() interface{} {
				return new(entry)
			},
		}
	}

	p.waitGroupOutPut.Add(1)
	go func() {
		defer func() {
			if util.IsRelease() {
				if err := recover(); err != nil {
					PrintErr(libconstant.GoroutinePanic, err, debug.Stack())
				}
			}
			p.waitGroupOutPut.Done()
			PrintInfo(libconstant.GoroutineDone)
		}()
		p.doLog()
	}()
	return nil
}

// getLogDuration 取得日志刻度
func (p *Mgr) getLogDuration(sec int64) int {
	var logFormat string
	if util.IsRelease() {
		logFormat = "2006010215" //年月日小时
	} else {
		logFormat = "20060102" //年月日
	}

	durationStr := time.Unix(sec, 0).Format(logFormat)
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		PrintfErr("strconv.Atoi sec:%v durationStr:%v err:%v", sec, durationStr, err)
	}
	return duration
}

// doLog 处理日志
func (p *Mgr) doLog() {
	for v := range p.logChan {
		p.fireHooks(v)

		// 检查自动切换日志
		if p.logDuration != p.getLogDuration(v.time.Unix()) {
			if err := p.newWriters(); err != nil {
				PrintfErr("log duration changed, init writers failed, err:%v", err)
				if p.options.IsEnablePool() {
					v.reset()
					p.pool.Put(v)
				}
				continue
			}
		}
		if *p.options.isWriteFile {
			p.loggerSlice[v.level].Print(v.formatMessage())
		}
		if p.options.IsEnablePool() {
			v.reset()
			p.pool.Put(v)
		}
	}
	// goroutine 退出,再设置chan为nil, (如果没有退出就设置为nil, 读chan == nil  会 block)
	p.logChan = nil
}

// SetLevel 设置日志等级
func (p *Mgr) SetLevel(lv Level) error {
	if lv < LevelOff || LevelOn < lv {
		return errors.WithMessage(liberror.Level, util.GetCodeLocation(1).String())
	}
	p.options.WithLevel(lv)
	return nil
}

// newWriters 初始化各级别的日志输出
func (p *Mgr) newWriters() error {
	// 检查是否要关闭文件
	for i := range p.openFiles {
		if err := p.openFiles[i].Close(); err != nil {
			return errors.WithMessage(err, util.GetCodeLocation(1).String())
		}
	}

	second := time.Now().Unix()
	duration := p.getLogDuration(second)
	accessWriter, err := newAccessFileWriter(*p.options.absPath, *p.options.namePrefix, duration)
	if err != nil {
		return errors.WithMessage(err, util.GetCodeLocation(1).String())
	}
	errorWriter, err := newErrorFileWriter(*p.options.absPath, *p.options.namePrefix, duration)
	if err != nil {
		return errors.WithMessage(err, util.GetCodeLocation(1).String())
	}
	p.logDuration = duration

	allWriter := io.MultiWriter(accessWriter, errorWriter)

	// 标准输出,标准错误重定向
	stdOut.SetOutput(accessWriter)
	stdErr.SetOutput(allWriter)

	p.loggerSlice[LevelTrace].SetOutput(accessWriter)
	p.loggerSlice[LevelDebug].SetOutput(accessWriter)
	p.loggerSlice[LevelInfo].SetOutput(accessWriter)
	p.loggerSlice[LevelWarn].SetOutput(allWriter)
	p.loggerSlice[LevelError].SetOutput(allWriter)
	p.loggerSlice[LevelFatal].SetOutput(allWriter)

	// 记录打开的文件
	p.openFiles = p.openFiles[0:0]
	p.openFiles = append(p.openFiles, accessWriter)
	p.openFiles = append(p.openFiles, errorWriter)

	return nil
}

// Stop 停止
func (p *Mgr) Stop() error {
	if p.logChan != nil {
		// close chan, for range 读完chan会退出.
		close(p.logChan)

		// 等待logChan 的for range 退出.
		p.waitGroupOutPut.Wait()
	}

	// 检查是否要关闭文件
	if len(p.openFiles) > 0 {
		for i := range p.openFiles {
			_ = p.openFiles[i].Close()
		}
		p.openFiles = p.openFiles[0:0]
	}
	return nil
}

// fireHooks 处理钩子
func (p *Mgr) fireHooks(entry *entry) {
	if 0 == len(p.options.hooks) {
		return
	}

	err := p.options.hooks.fire(entry.level, entry)
	if err != nil {
		PrintfErr("failed to fire hook. err:%v", err)
	}
}

// WithField 由field创建日志信息 默认大小2(cap:2*2=4)
func (p *Mgr) WithField(key string, value interface{}) *entry {
	entry := newEntry(p)
	entry.fields = make(fields, 0, 4)
	return entry.WithField(key, value)
}

// WithFields 由fields创建日志信息 默认大小4(cap:4*2=8)
func (p *Mgr) WithFields(f fields) *entry {
	entry := newEntry(p)
	entry.fields = make(fields, 0, 8)
	return entry.WithFields(f)
}

// WithContext 由ctx创建日志信息
func (p *Mgr) WithContext(ctx context.Context) *entry {
	entry := newEntry(p)
	return entry.WithContext(ctx)
}

// log 记录日志
func (p *Mgr) log(level Level, v ...interface{}) {
	entry := newEntry(p)
	entry.log(level, calldepth3, v...)
}

// logf 记录日志
func (p *Mgr) logf(level Level, format string, v ...interface{}) {
	entry := newEntry(p)
	entry.logf(level, calldepth3, format, v...)
}

// Trace 踪迹日志
func (p *Mgr) Trace(v ...interface{}) {
	if *p.options.level < LevelTrace {
		return
	}
	p.log(LevelTrace, v...)
}

// Tracef 踪迹日志
func (p *Mgr) Tracef(format string, v ...interface{}) {
	if *p.options.level < LevelTrace {
		return
	}
	p.logf(LevelTrace, format, v...)
}

// Debug 调试日志
func (p *Mgr) Debug(v ...interface{}) {
	if *p.options.level < LevelDebug {
		return
	}
	p.log(LevelDebug, v...)
}

// Debugf 调试日志
func (p *Mgr) Debugf(format string, v ...interface{}) {
	if *p.options.level < LevelDebug {
		return
	}
	p.logf(LevelDebug, format, v...)
}

// Info 信息日志
func (p *Mgr) Info(v ...interface{}) {
	if *p.options.level < LevelInfo {
		return
	}
	p.log(LevelInfo, v...)
}

// Infof 信息日志
func (p *Mgr) Infof(format string, v ...interface{}) {
	if *p.options.level < LevelInfo {
		return
	}
	p.logf(LevelInfo, format, v...)
}

// Warn 警告日志
func (p *Mgr) Warn(v ...interface{}) {
	if *p.options.level < LevelWarn {
		return
	}
	p.log(LevelWarn, v...)
}

// Warnf 警告日志
func (p *Mgr) Warnf(format string, v ...interface{}) {
	if *p.options.level < LevelWarn {
		return
	}
	p.logf(LevelWarn, format, v...)
}

// Error 错误日志
func (p *Mgr) Error(v ...interface{}) {
	if *p.options.level < LevelError {
		return
	}
	p.log(LevelError, v...)
}

// Errorf 错误日志
func (p *Mgr) Errorf(format string, v ...interface{}) {
	if *p.options.level < LevelError {
		return
	}
	p.logf(LevelError, format, v...)
}

// Fatal 致命日志
func (p *Mgr) Fatal(v ...interface{}) {
	if *p.options.level < LevelFatal {
		return
	}
	p.log(LevelFatal, v...)
}

// Fatalf 致命日志
func (p *Mgr) Fatalf(format string, v ...interface{}) {
	if *p.options.level < LevelFatal {
		return
	}
	p.logf(LevelFatal, format, v...)
}
