package error

import (
	"fmt"
	"runtime"
	xrconstant "social/pkg/lib/constant"
)

// codeLocationInfo 代码位置信息
type codeLocationInfo struct {
	FileName string // 文件名
	FuncName string // 函数名
	Line     int    // 行数
}

// 错误信息
func (p *codeLocationInfo) Error() string {
	return fmt.Sprintf("file:%v line:%v func:%v", p.FileName, p.Line, p.FuncName)
}

// getCodeLocationInfo 获取代码位置信息
func getCodeLocationInfo(skip int) *codeLocationInfo {
	c := &codeLocationInfo{}
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		c.FileName = xrconstant.Unknown
		c.Line = 0
		c.FuncName = xrconstant.Unknown
	} else {
		c.FileName = file
		c.Line = line
		c.FuncName = runtime.FuncForPC(pc).Name()
	}
	return c
}
