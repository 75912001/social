package item

type dependency interface {
	UserGetName() string
}

func NewItem(dep dependency) *Item {
	return &Item{
		Deps: dep,
		Name: "Item",
	}
}

type Item struct {
	Deps dependency
	Name string
}

func (p *Item) ItemGetName() string {
	return p.Name + p.Deps.UserGetName()
}
