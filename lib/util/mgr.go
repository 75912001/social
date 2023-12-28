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
func (p *Mgr[TKey, TVal]) Find(key TKey) TVal {
	data, _ := p.elementMap[key]
	return data
}

// Del 删除元素
func (p *Mgr[TKey, TVal]) Del(key TKey) {
	delete(p.elementMap, key)
}

// ConditionFunc 用于定义查找条件的函数签名
type ConditionFunc[TKey comparable, TVal IObject] func(key TKey, value TVal) bool

// FindOneWithCondition 根据条件查找元素-一个
func (p *Mgr[TKey, TVal]) FindOneWithCondition(condition ConditionFunc[TKey, TVal]) (val TVal) {
	for key, value := range p.elementMap {
		if condition(key, value) {
			return value
		}
	}
	return val
}
