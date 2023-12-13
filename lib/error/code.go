package error

var (
	// Success 成功
	Success = createError(0x000, "Success", "success")
	// Link 链接
	Link = createError(0xf001, "Link", "link error")
	// System 系统
	System = createError(0xf002, "System", "system error")
	// Param 参数
	Param = createError(0xf003, "Param", "parameter error")
	// Packet 数据包
	Packet = createError(0xf004, "Packet", "packet error")
	// Timeout 超时
	Timeout = createError(0xf005, "Timeout", "time out")
	// ChannelFull 通道 满
	ChannelFull = createError(0xf006, "ChannelFull", "channel full")
	// ChannelEmpty 通道 空
	ChannelEmpty = createError(0xf007, "ChannelEmpty", "channel empty")
	// OutOfRange 超出范围
	OutOfRange = createError(0xf008, "OutOfRange", "out of range")
	// InvalidValue 无效数值
	InvalidValue = createError(0xf009, "InvalidValue", "invalid value")
	// Conflict 冲突
	Conflict = createError(0xf00a, "Conflict", "conflict")
	// TypeMismatch 类型不匹配
	TypeMismatch = createError(0xf00b, "TypeMismatch", "type mismatch")
	// InvalidPointer 无效指针
	InvalidPointer = createError(0xf00c, "InvalidPointer", "invalid pointer")
	// Level 等级
	Level = createError(0xf00d, "level", "level error")
	// NonExistent 不存在
	NonExistent = createError(0xf00e, "NonExistent", "non-existent")
	// Exists 存在
	Exists = createError(0xf00f, "Exists", "exists")
	// Marshal 序列化
	Marshal = createError(0xf010, "Marshal", "marshal")
	// Unmarshal 反序列化
	Unmarshal = createError(0xf011, "Unmarshal", "unmarshal")
	// Insert 插入
	Insert = createError(0xf012, "Insert", "insert error")
	// Find 查找
	Find = createError(0xf013, "Find", "find error")
	// Update 更新
	Update = createError(0xf014, "Update", "update error")
	// Delete 删除
	Delete = createError(0xf015, "Delete", "delete error")
	// Duplicate 重复
	Duplicate = createError(0xf016, "Duplicate", "duplicate error")
	// Config 配置
	Config = createError(0xf017, "Config", "config error")
	// InvalidOperation 无效操作
	InvalidOperation = createError(0xf018, "InvalidOperation", "invalid operation")
	// IllConditioned 条件不足
	IllConditioned = createError(0xf019, "IllConditioned", "ill conditioned")
	// PermissionDenied 没有权限
	PermissionDenied = createError(0xf01a, "PermissionDenied", "permission denied")
	// BlockedAccount 冻结账号
	BlockedAccount = createError(0xf01b, "BlockedAccount", "blocked account")
	// Send 发送
	Send = createError(0xf01c, "Send", "send")
	// Configure 给配置
	Configure = createError(0xf01d, "Configure", "configure")
	// Retry 重试
	Retry = createError(0xf01e, "Retry", "retry")
	// MessageIDNonExistent 消息ID 不存在
	MessageIDNonExistent = createError(0xf01f, "MessageIDNonExistent", "message id non-existent")
	// Redis 系统 Redis
	Redis = createError(0xf020, "Redis", "redis")
	// Busy 繁忙
	Busy = createError(0xf021, "Busy", "busy")
	// OutOfResources 资源不足
	OutOfResources = createError(0xf022, "OutOfResources", "out of resources")
	// NATS NATS错误
	NATS = createError(0xf023, "NATS", "nats")
	// PacketQuantityLimit 包数量限制
	PacketQuantityLimit = createError(0xf024, "PacketQuantityLimit", "packet quantity limit")
	// OverloadWarning 过载-告警
	OverloadWarning = createError(0xf025, "OverloadWarning", "overload warning")
	// OverloadError 过载-错误
	OverloadError = createError(0xf026, "OverloadError", "overload error")
	// MessageIDDisable 消息ID禁用
	MessageIDDisable = createError(0xf027, "MessageIDDisable", "message id is disabled")
	// MessageIDExistent 消息ID 存在
	MessageIDExistent = createError(0xf028, "MessageIDExistent", "message id existent")
	// ModeMismatch 模式 不匹配
	ModeMismatch = createError(0xf029, "ModeMismatch", "mode mismatch")
	// FormatMismatch 格式 不匹配
	FormatMismatch = createError(0xf02a, "FormatMismatch", "format mismatch")
	// MISSING 找不到,丢失,未命中
	MISSING = createError(0xf02b, "MISSING", "missing")
	// VersionMismatch 版本 不匹配
	VersionMismatch = createError(0xf02c, "VersionMismatch", "version mismatch")
	// Unavailable 不可用
	Unavailable = createError(0xf02d, "Unavailable", "unavailable")
	// NotImplemented 未实现
	NotImplemented = createError(0xf02e, "NotImplemented", "not implemented")
	// Unknown 未知
	Unknown = createError(0xffff, "Unknown", "unknown")
	// 0xffff
)

// 内部创建使用
func createError(code uint32, name string, desc string) *Error {
	newErr := NewError(code, name, desc)
	e := CheckForDuplicates(newErr)
	if e != nil {
		panic(e)
	}
	return newErr
}
