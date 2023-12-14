package runtime

import (
	"fmt"
	"runtime"
	libconsts "social/lib/consts"
)

// CodeLocation 代码位置
type CodeLocation struct {
	FileName string //文件名
	FuncName string //函数名
	Line     int    //行数
}

// Error 错误信息
func (p *CodeLocation) Error() string {
	return fmt.Sprintf("file:%v line:%v func:%v", p.FileName, p.Line, p.FuncName)
}

// String 错误信息
func (p *CodeLocation) String() string {
	return p.Error()
}

// GetCodeLocation 获取代码位置
//
//	参数:
//		skip:The argument skip is the number of stack frames to ascend, with 0 identifying the caller of Caller.
func GetCodeLocation(skip int) *CodeLocation {
	c := &CodeLocation{
		FileName: libconsts.Unknown,
		FuncName: libconsts.Unknown,
	}

	pc, fileName, line, ok := runtime.Caller(skip)

	if ok {
		c.FileName = fileName
		c.Line = line
		c.FuncName = runtime.FuncForPC(pc).Name()
	}
	return c
}

// Location 获取代码位置
func Location() string {
	return GetCodeLocation(2).String()
}
