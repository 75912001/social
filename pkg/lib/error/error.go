// Package error 错误
package error

import (
	"fmt"

	"github.com/pkg/errors"
)

// 错误信息
var errorMap = make(map[uint32]struct{})

// Register 注册, 为了检测是否重复
func Register(err *Error) error {
	if _, ok := errorMap[err.Code]; ok { //不可重复
		return errors.WithMessage(Exists, getCodeLocationInfo(1).Error())
	}
	errorMap[err.Code] = struct{}{}
	return nil
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

func NewError(e *Error) *Error {
	return &Error{
		Code: e.Code,
		Name: e.Name,
		Desc: e.Desc,
	}
}

func (p *Error) WithExtraMessage(extraMessage string) *Error {
	p.ExtraMessage = extraMessage
	return p
}

func (p *Error) WithExtraError(extraError error) *Error {
	p.ExtraError = extraError
	return p
}