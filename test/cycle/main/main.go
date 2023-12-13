package main

import (
	"fmt"
	"social/test/cycle/user"
)

func main() {
	u := user.NewUser()
	fmt.Println(u.UserGetName())
	fmt.Println(u.Item.ItemGetName())
}
