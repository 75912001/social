package pb

import (
	"github.com/pkg/errors"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
)

type messageMap map[uint32]*Message

// Mgr 管理器
type Mgr struct {
	messageMap messageMap
}

// Init 初始化管理器
func (p *Mgr) Init() {
	p.messageMap = make(messageMap)
}

// Register 注册消息
func (p *Mgr) Register(messageID uint32, messageSlice ...*Message) error {
	if pb := p.Find(messageID); pb != nil {
		return errors.WithMessagef(liberror.MessageIDExistent, "%v messageID:%#x %v",
			libruntime.Location(), messageID, messageID)
	}
	message := merge(messageSlice...)
	if err := configure(message); err != nil {
		return errors.WithMessagef(err, "%v messageID:%#x %v", libruntime.Location(), messageID, messageID)
	}
	p.messageMap[messageID] = message
	return nil
}

func (p *Mgr) Find(messageID uint32) *Message {
	return p.messageMap[messageID]
}

// Replace 替换/覆盖(Override)
func (p *Mgr) Replace(messageID uint32, messageEntity *Message) error {
	if err := configure(messageEntity); err != nil {
		return errors.WithMessagef(err, "%v messageID:%#x %v", libruntime.Location(), messageID, messageID)
	}
	p.messageMap[messageID] = messageEntity

	return nil
}
