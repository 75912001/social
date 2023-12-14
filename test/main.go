package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	WithStack()
}

// todo menglingchao 可以用于关键信息记录
func WithStack() {
	// 创建一个基本错误
	baseError := errors.New("This is a basic error")

	// 使用 WithStack 将调用栈信息添加到错误中
	errorWithStack := errors.WithStack(baseError)

	// 打印错误信息（包括调用栈信息）
	fmt.Printf("%+v\n", errorWithStack)
}
