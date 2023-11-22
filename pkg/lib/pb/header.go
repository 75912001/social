package pb

type IHeader interface {
	Pack() []byte  // 将 成员变量 -> data 中
	Unpack([]byte) // 将 data 数据 -> 成员变量中
}

type IBody interface {
}

type Body struct {
	RawData []byte
}
