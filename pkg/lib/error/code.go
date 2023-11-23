package error

import "github.com/pkg/errors"

var (
	// Success 成功
	Success = &Error{Code: 0x000, Name: "Success", Desc: "success"}
	// Link 链接
	Link = &Error{Code: 0xf001, Name: "Link", Desc: "link error"}
	// System 系统
	System = &Error{Code: 0xf002, Name: "System", Desc: "system error"}
	// Param 参数
	Param = &Error{Code: 0xf003, Name: "Param", Desc: "parameter error"}
	// Packet 数据包
	Packet = &Error{Code: 0xf004, Name: "Packet", Desc: "packet error"}
	// Timeout 超时
	Timeout = &Error{Code: 0xf005, Name: "Timeout", Desc: "time out"}
	// ChannelFull 通道 满
	ChannelFull = &Error{Code: 0xf006, Name: "ChannelFull", Desc: "channel full"}
	// ChannelEmpty 通道 空
	ChannelEmpty = &Error{Code: 0xf007, Name: "ChannelEmpty", Desc: "channel empty"}
	// OutOfRange 超出范围
	OutOfRange = &Error{Code: 0xf008, Name: "OutOfRange", Desc: "out of range"}
	// InvalidValue 无效数值
	InvalidValue = &Error{Code: 0xf009, Name: "InvalidValue", Desc: "invalid value"}
	// Conflict 冲突
	Conflict = &Error{Code: 0xf00a, Name: "Conflict", Desc: "conflict"}
	// TypeMismatch 类型不匹配
	TypeMismatch = &Error{Code: 0xf00b, Name: "TypeMismatch", Desc: "type mismatch"}
	// InvalidPointer 无效指针
	InvalidPointer = &Error{Code: 0xf00c, Name: "InvalidPointer", Desc: "invalid pointer"}
	// Level 等级
	Level = &Error{Code: 0xf00d, Name: "level", Desc: "level error"}
	// NonExistent 不存在
	NonExistent = &Error{Code: 0xf00e, Name: "NonExistent", Desc: "non-existent"}
	// Exists 存在
	Exists = &Error{Code: 0xf00f, Name: "Exists", Desc: "exists"}
	// Marshal 序列化
	Marshal = &Error{Code: 0xf010, Name: "Marshal", Desc: "marshal"}
	// Unmarshal 反序列化
	Unmarshal = &Error{Code: 0xf011, Name: "Unmarshal", Desc: "unmarshal"}
	// Insert 插入
	Insert = &Error{Code: 0xf012, Name: "Insert", Desc: "insert error"}
	// Find 查找
	Find = &Error{Code: 0xf013, Name: "Find", Desc: "find error"}
	// Update 更新
	Update = &Error{Code: 0xf014, Name: "Update", Desc: "update error"}
	// Delete 删除
	Delete = &Error{Code: 0xf015, Name: "Delete", Desc: "delete error"}
	// Duplicate 重复
	Duplicate = &Error{Code: 0xf016, Name: "Duplicate", Desc: "duplicate error"}
	// Config 配置
	Config = &Error{Code: 0xf017, Name: "Config", Desc: "config error"}
	// InvalidOperation 无效操作
	InvalidOperation = &Error{Code: 0xf018, Name: "InvalidOperation", Desc: "invalid operation"}
	// IllConditioned 条件不足
	IllConditioned = &Error{Code: 0xf019, Name: "IllConditioned", Desc: "ill conditioned"}
	// PermissionDenied 没有权限
	PermissionDenied = &Error{Code: 0xf01a, Name: "PermissionDenied", Desc: "permission denied"}
	// BlockedAccount 冻结账号
	BlockedAccount = &Error{Code: 0xf01b, Name: "BlockedAccount", Desc: "blocked account"}
	// Send 发送
	Send = &Error{Code: 0xf01c, Name: "Send", Desc: "send"}
	// Configure 给配置
	Configure = &Error{Code: 0xf01d, Name: "Configure", Desc: "configure"}
	// Retry 重试
	Retry = &Error{Code: 0xf01e, Name: "Retry", Desc: "retry"}
	// MessageIDNonExistent 消息ID 不存在
	MessageIDNonExistent = &Error{Code: 0xf01f, Name: "MessageIDNonExistent", Desc: "message id non-existent"}
	// Redis 系统 Redis
	Redis = &Error{Code: 0xf020, Name: "Redis", Desc: "redis"}
	// Busy 繁忙
	Busy = &Error{Code: 0xf021, Name: "Busy", Desc: "busy"}
	// OutOfResources 资源不足
	OutOfResources = &Error{Code: 0xf022, Name: "OutOfResources", Desc: "out of resources"}
	// NATS NATS错误
	NATS = &Error{Code: 0xf023, Name: "NATS", Desc: "nats"}
	// PacketQuantityLimit 包数量限制
	PacketQuantityLimit = &Error{Code: 0xf024, Name: "PacketQuantityLimit", Desc: "packet quantity limit"}
	// OverloadWarning 过载-告警
	OverloadWarning = &Error{Code: 0xf025, Name: "OverloadWarning", Desc: "overload warning"}
	// OverloadError 过载-错误
	OverloadError = &Error{Code: 0xf026, Name: "OverloadError", Desc: "overload error"}
	// MessageIDDisable 消息ID禁用
	MessageIDDisable = &Error{Code: 0xf027, Name: "MessageIDDisable", Desc: "message id is disabled"}
	// MessageIDExistent 消息ID 存在
	MessageIDExistent = &Error{Code: 0xf028, Name: "MessageIDExistent", Desc: "message id existent"}
	// ModeMismatch 模式 不匹配
	ModeMismatch = &Error{Code: 0xf029, Name: "ModeMismatch", Desc: "mode mismatch"}
	// FormatMismatch 格式 不匹配
	FormatMismatch = &Error{Code: 0xf02a, Name: "FormatMismatch", Desc: "format mismatch"}
	// MISSING 找不到,丢失,未命中
	MISSING = &Error{Code: 0xf02b, Name: "MISSING", Desc: "missing"}
	// VersionMismatch 版本 不匹配
	VersionMismatch = &Error{Code: 0xf02c, Name: "VersionMismatch", Desc: "version mismatch"}
	// Unavailable 不可用
	Unavailable = &Error{Code: 0xf02d, Name: "Unavailable", Desc: "unavailable"}
	// Unknown 未知
	Unknown = &Error{Code: 0xf02e, Name: "Unknown", Desc: "unknown"}
	// 0xffff
)

func init() {
	_ = Register(Success)
	_ = Register(Link)
	_ = Register(System)
	_ = Register(Param)
	_ = Register(Packet)
	_ = Register(Timeout)
	_ = Register(ChannelFull)
	_ = Register(ChannelEmpty)
	_ = Register(OutOfRange)
	_ = Register(InvalidValue)
	_ = Register(Conflict)
	_ = Register(TypeMismatch)
	_ = Register(InvalidPointer)
	_ = Register(Level)
	_ = Register(NonExistent)
	_ = Register(Exists)
	_ = Register(Marshal)
	_ = Register(Unmarshal)
	_ = Register(Insert)
	_ = Register(Find)
	_ = Register(Update)
	_ = Register(Delete)
	_ = Register(Duplicate)
	_ = Register(Config)
	_ = Register(InvalidOperation)
	_ = Register(IllConditioned)
	_ = Register(PermissionDenied)
	_ = Register(BlockedAccount)
	_ = Register(Send)
	_ = Register(Configure)
	_ = Register(Retry)
	_ = Register(MessageIDNonExistent)
	_ = Register(Redis)
	_ = Register(Busy)
	_ = Register(OutOfResources)
	_ = Register(NATS)
	_ = Register(PacketQuantityLimit)
	_ = Register(OverloadWarning)
	_ = Register(OverloadError)
	_ = Register(MessageIDDisable)
	_ = Register(MessageIDExistent)
	_ = Register(ModeMismatch)
	_ = Register(FormatMismatch)
	_ = Register(MISSING)
	_ = Register(VersionMismatch)
	_ = Register(Unavailable)

	_ = Register(Unknown)
	if len(errorMap) != (int(Unknown.Code+2) - 0xf001) { // 未能初始化 全部 错误
		panic(errors.New("failed to initialize all errors"))
	}
}
