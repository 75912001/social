package pkg

const EtcdWatchMsgTypeService string = "service"
const EtcdWatchMsgTypeCommand string = "command"

const EtcdTtlSecondDefault int64 = 33 //默认TTL时间 秒

var (
	GZoneID      uint32 // 区域ID
	GServiceName string // 服务
	GServiceID   uint32 // 服务ID
)
