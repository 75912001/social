package bench

import "time"

type ISubBench interface {
	Load(strJson string) error // 加载 配置文件
}

type MongoDB struct {
	Addrs           []string       `json:"addrs"`
	User            string         `json:"user"`
	Pwd             string         `json:"pwd"`
	DBName          *string        `json:"dbName"`          // 数据库名称 default: common.MongodbDatabaseNameDefault
	MaxPoolSize     *uint64        `json:"maxPoolSize"`     // 连接池最大数量,该数量应该与并发数量匹配 default: common.MongodbMaxPoolSizeDefault
	MinPoolSize     *uint64        `json:"minPoolSize"`     // 池最小数量 default: common.MongodbMinPoolSizeDefault
	TimeoutDuration *time.Duration `json:"timeoutDuration"` // 操作超时时间 default: common.MongodbTimeoutDurationDefault
	MaxConnIdleTime *time.Duration `json:"maxConnIdleTime"` // 指定连接在连接池中保持空闲的最长时间 default: common.MongodbMaxConnIdleTimeDefault
	MaxConnecting   *uint64        `json:"maxConnecting"`   // 指定连接池可以同时建立的最大连接数 default: common.MongodbMaxConnectingDefault
}

type Redis struct {
	Addrs    []string `json:"addrs"`
	Password string   `json:"password"`
}

type DBAsync struct {
	ChanCnt      uint32 `json:"chanCnt"`              //DB异步消费chan数量. 为0或者没有则不开启异步消费
	Model        uint32 `json:"model"`                //DB异步消费模式 default consumer.ModelAsyncOne
	BulkWriteMax uint32 `json:"bulkWriteMax"`         //DB合并写 单个集合最大合批数量 default common.BulkWriteMax
	BulkWriteMs  uint32 `json:"bulkWriteMillisecond"` //DB合并写周期  单位毫秒 default common.BatchExecMaxMilliSecond
}
