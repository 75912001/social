package mongodb

import "time"

var MaxPoolSizeDefault uint64 = 24 // 可以设置为 2*P + 30% 的数量

var MinPoolSizeDefault uint64 = 8 // 可以设置为 P 的数量

var TimeoutDurationDefault = time.Minute

var MaxConnIdleTimeDefault = time.Minute * 5

var MaxConnectingDefault uint64 = 4 // 设置为 P/2
