package pb

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"reflect"
	"runtime"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	"strings"
)

// Handler 处理函数
type Handler func(header IHeader, message proto.Message, obj interface{}) *liberror.Error

// NewPBMessage 创建新的 proto.Message
type NewPBMessage func() proto.Message

// Message contains options to configure a server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type Message struct {
	handler      Handler      // 消息处理函数
	newPBMessage NewPBMessage // 创建新的 proto.Message
	name         string       // [optional] default:"Unknown" 名称
}

func (p *Message) GetName() string {
	return p.name
}

// NewMessage 创建 Message
func NewMessage() *Message {
	return new(Message)
}

// Unmarshal 反序列化
//
//	message: 反序列化 得到的 消息
func (p *Message) Unmarshal(data []byte) (message proto.Message, err error) {
	message = p.newPBMessage()
	err = proto.Unmarshal(data, message)
	if err != nil {
		return nil, errors.WithMessage(err, libruntime.Location())
	}
	return message, nil
}

// Handler 处理
func (p *Message) Handler(header IHeader, message proto.Message, obj interface{}) *liberror.Error {
	return p.handler(header, message, obj)
}

// MutableCopy 深拷贝
func MutableCopy(src proto.Message, dst proto.Message) error {
	data, err := proto.Marshal(src)
	if err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	err = proto.Unmarshal(data, dst)
	if err != nil {
		return errors.WithMessage(err, libruntime.Location())
	}
	return nil
}

// 生成函数名称
func getFuncName(i interface{}, seps ...rune) string {
	if i == nil {
		return ""
	}
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	fields := strings.FieldsFunc(funcName, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return "Unknown"
}

func (p *Message) SetHandler(handler Handler) *Message {
	p.handler = handler
	p.name = getFuncName(handler, '.')
	return p
}

func (p *Message) SetNewPBMessage(newPBMessage NewPBMessage) *Message {
	p.newPBMessage = newPBMessage
	return p
}

// merge combines the given *Options into a single *Options in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func merge(opts ...*Message) *Message {
	so := NewMessage()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.handler != nil {
			so.handler = opt.handler
		}
		if opt.newPBMessage != nil {
			so.newPBMessage = opt.newPBMessage
		}
		if 0 < len(opt.name) {
			so.name = opt.name
		}
	}
	return so
}

// 配置
func configure(opts *Message) error {
	if opts.handler == nil { // 没有消息处理函数
		return errors.WithMessage(liberror.Param, libruntime.Location())
	}
	if opts.newPBMessage == nil { // 没有创建消息函数
		return errors.WithMessage(liberror.Param, libruntime.Location())
	}
	return nil
}
