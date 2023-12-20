package actor

// IActor 是 Actor 接口
type IActor[TKey comparable] interface {
	GetKey() TKey
}
