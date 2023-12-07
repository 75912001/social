package item

import (
	"social/test/cycle/module"
)

type dependency interface {
	UserGetName() string
}

func NewItem(dep module.Modules) *Item {
	return &Item{
		Deps: dep,
		Name: "Item",
	}
}

type Item struct {
	Deps module.Modules
	Name string
}

func (p *Item) ItemGetName() string {
	return p.Name + p.Deps.UserGetName()
}
