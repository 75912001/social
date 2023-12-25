// Package func_mgr function 管理器
// 通过function id (uint32) 来绑定一个处理
package func_mgr

import (
	"github.com/pkg/errors"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
)

// Function 协议function
type Function func(arg ...interface{}) (interface{}, error)

// mgr 管理器
type mgr struct {
	functionMap functionMap
}

func NewMgr() *mgr {
	p := new(mgr)
	p.functionMap = make(functionMap)
	return p
}

// Register 注册
func (p *mgr) Register(funcID uint32, fun Function) error {
	if pb := p.Find(funcID); pb != nil {
		return errors.WithMessage(liberror.MessageIDExistent, libruntime.Location())
	}
	p.functionMap[funcID] = fun

	return nil
}

// functionMap 协议 function map
type functionMap map[uint32]Function //key: funcID, val:function

func (p *mgr) Find(funcID uint32) Function {
	return p.functionMap[funcID]
}
