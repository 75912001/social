package main

import (
	"fmt"
	"social/test/cycle/user"
)

func main() {
	user := user.NewUser()
	fmt.Println(user.UserGetName())
	fmt.Println(user.Item.ItemGetName())
}
