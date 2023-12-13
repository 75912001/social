package etcd

import "time"

var (
	grantLeaseRetryDuration           = time.Second * 3 // 授权租约 重试 间隔时长
	grantLeaseMaxRetriesDefault       = 600             // 授权租约 最大 重试次数 默认值
	dialTimeoutDefault                = time.Second * 5 //dialTimeout is the timeout for failing to establish a connection. 默认值
	TtlSecondDefault            int64 = 33              //默认TTL时间 秒
)
