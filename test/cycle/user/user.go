package user

import (
	"social/test/cycle/item"
)

func NewUser() *User {
	user := &User{
		Name: "User",
	}
	user.Item = item.NewItem(user)
	return user
}

type User struct {
	Item *item.Item
	Name string
}

func (p *User) UserGetName() string {
	return p.Name
}
