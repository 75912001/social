// Package error 错误
package error

import (
	"fmt"
	"github.com/pkg/errors"
)

// 错误信息
var errMap = make(map[uint32]struct{})

// Register 注册, 为了检测是否重复
func Register(err *Error) error {
	if _, ok := errMap[err.Code]; ok { //不可重复
		return errors.WithMessage(Duplicate, getCodeLocation(1).Error())
	}
	errMap[err.Code] = struct{}{}
	return nil
}

func NewError(code uint32, name string, desc string) *Error {
	return &Error{
		Code: code,
		Name: name,
		Desc: desc,
	}
}

// Error 错误
type Error struct {
	Code         uint32 // 码
	Name         string // 名称
	Desc         string // 描述 Description
	ExtraMessage string // 附加信息
	ExtraError   error  // 附加错误
}

// 错误信息
func (p *Error) Error() string {
	if Success.Code == p.Code {
		return ""
	}
	return fmt.Sprintf("name:%v code:%v %#x description:%v extraMessage:%v extraError%v",
		p.Name, p.Code, p.Code, p.Desc, p.ExtraMessage, p.ExtraError)
}

func (p *Error) WithExtraMessage(extraMessage string) *Error {
	p.ExtraMessage = extraMessage
	return p
}

func (p *Error) WithExtraError(extraError error) *Error {
	p.ExtraError = extraError
	return p
}
