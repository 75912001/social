package util

// NewMgr 创建 Mgr 实例
func NewMgr[TKey comparable, TVal IObject]() *Mgr[TKey, TVal] {
	return &Mgr[TKey, TVal]{
		elementMap: make(map[TKey]TVal),
	}
}

type Mgr[TKey comparable, TVal IObject] struct {
	elementMap map[TKey]TVal
}

// Add 添加元素
func (p *Mgr[TKey, TVal]) Add(key TKey, value TVal) {
	p.elementMap[key] = value
}

// Find 查找元素
func (p *Mgr[TKey, TVal]) Find(key TKey) (TVal, bool) {
	data, ok := p.elementMap[key]
	return data, ok
}

// Del 删除元素
func (p *Mgr[TKey, TVal]) Del(key TKey) {
	delete(p.elementMap, key)
}
