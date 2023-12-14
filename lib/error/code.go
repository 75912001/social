package error

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	// Success 成功
	Success = CreateError(0x000, "Success", "success")
	// Link 链接
	Link = CreateError(0xf001, "Link", "link error")
	// System 系统
	System = CreateError(0xf002, "System", "system error")
	// Param 参数
	Param = CreateError(0xf003, "Param", "parameter error")
	// Packet 数据包
	Packet = CreateError(0xf004, "Packet", "packet error")
	// Timeout 超时
	Timeout = CreateError(0xf005, "Timeout", "time out")
	// ChannelFull 通道 满
	ChannelFull = CreateError(0xf006, "ChannelFull", "channel full")
	// ChannelEmpty 通道 空
	ChannelEmpty = CreateError(0xf007, "ChannelEmpty", "channel empty")
	// OutOfRange 超出范围
	OutOfRange = CreateError(0xf008, "OutOfRange", "out of range")
	// InvalidValue 无效数值
	InvalidValue = CreateError(0xf009, "InvalidValue", "invalid value")
	// Conflict 冲突
	Conflict = CreateError(0xf00a, "Conflict", "conflict")
	// TypeMismatch 类型不匹配
	TypeMismatch = CreateError(0xf00b, "TypeMismatch", "type mismatch")
	// InvalidPointer 无效指针
	InvalidPointer = CreateError(0xf00c, "InvalidPointer", "invalid pointer")
	// Level 等级
	Level = CreateError(0xf00d, "level", "level error")
	// NonExistent 不存在
	NonExistent = CreateError(0xf00e, "NonExistent", "non-existent")
	// Exists 存在
	Exists = CreateError(0xf00f, "Exists", "exists")
	// Marshal 序列化
	Marshal = CreateError(0xf010, "Marshal", "marshal")
	// Unmarshal 反序列化
	Unmarshal = CreateError(0xf011, "Unmarshal", "unmarshal")
	// Insert 插入
	Insert = CreateError(0xf012, "Insert", "insert error")
	// Find 查找
	Find = CreateError(0xf013, "Find", "find error")
	// Update 更新
	Update = CreateError(0xf014, "Update", "update error")
	// Delete 删除
	Delete = CreateError(0xf015, "Delete", "delete error")
	// Duplicate 重复
	Duplicate = CreateError(0xf016, "Duplicate", "duplicate error")
	// Config 配置
	Config = CreateError(0xf017, "Config", "config error")
	// InvalidOperation 无效操作
	InvalidOperation = CreateError(0xf018, "InvalidOperation", "invalid operation")
	// IllConditioned 条件不足
	IllConditioned = CreateError(0xf019, "IllConditioned", "ill conditioned")
	// PermissionDenied 没有权限
	PermissionDenied = CreateError(0xf01a, "PermissionDenied", "permission denied")
	// BlockedAccount 冻结账号
	BlockedAccount = CreateError(0xf01b, "BlockedAccount", "blocked account")
	// Send 发送
	Send = CreateError(0xf01c, "Send", "send")
	// Configure 给配置
	Configure = CreateError(0xf01d, "Configure", "configure")
	// Retry 重试
	Retry = CreateError(0xf01e, "Retry", "retry")
	// MessageIDNonExistent 消息ID 不存在
	MessageIDNonExistent = CreateError(0xf01f, "MessageIDNonExistent", "message id non-existent")
	// Redis 系统 Redis
	Redis = CreateError(0xf020, "Redis", "redis")
	// Busy 繁忙
	Busy = CreateError(0xf021, "Busy", "busy")
	// OutOfResources 资源不足
	OutOfResources = CreateError(0xf022, "OutOfResources", "out of resources")
	// NATS NATS错误
	NATS = CreateError(0xf023, "NATS", "nats")
	// PacketQuantityLimit 包数量限制
	PacketQuantityLimit = CreateError(0xf024, "PacketQuantityLimit", "packet quantity limit")
	// OverloadWarning 过载-告警
	OverloadWarning = CreateError(0xf025, "OverloadWarning", "overload warning")
	// OverloadError 过载-错误
	OverloadError = CreateError(0xf026, "OverloadError", "overload error")
	// MessageIDDisable 消息ID禁用
	MessageIDDisable = CreateError(0xf027, "MessageIDDisable", "message id is disabled")
	// MessageIDExistent 消息ID 存在
	MessageIDExistent = CreateError(0xf028, "MessageIDExistent", "message id existent")
	// ModeMismatch 模式 不匹配
	ModeMismatch = CreateError(0xf029, "ModeMismatch", "mode mismatch")
	// FormatMismatch 格式 不匹配
	FormatMismatch = CreateError(0xf02a, "FormatMismatch", "format mismatch")
	// MISSING 找不到,丢失,未命中
	MISSING = CreateError(0xf02b, "MISSING", "missing")
	// VersionMismatch 版本 不匹配
	VersionMismatch = CreateError(0xf02c, "VersionMismatch", "version mismatch")
	// Unavailable 不可用
	Unavailable = CreateError(0xf02d, "Unavailable", "unavailable")
	// NotImplemented 未实现
	NotImplemented = CreateError(0xf02e, "NotImplemented", "not implemented")
	// PacketHeaderLength 数据包头长度
	PacketHeaderLength = CreateError(0xf02f, "PacketHeaderLength", "packet header length error")
	// Unknown 未知
	Unknown = CreateError(0xffff, "Unknown", "unknown")
	// 0xffff
)

// CreateError 创建错误码,初始化程序的时候创建,创建失败会panic.
// [NOTE] 不要在程序运行的时候使用.
func CreateError(code uint32, name string, desc string) *Error {
	newErr := NewError(code, name, desc)
	e := CheckForDuplicates(newErr)
	if e != nil {
		panic(
			errors.WithMessage(e, fmt.Sprintf("create error duplicates %v %#x %v %v", code, code, name, desc)),
		)
	}
	return newErr
}
