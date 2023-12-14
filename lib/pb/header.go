package pb

type IHeader interface {
	Marshal() []byte  // 将 成员变量 -> data 中
	Unmarshal([]byte) // 将 data 数据 -> 成员变量中
}
