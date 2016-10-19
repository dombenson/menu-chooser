package main

import (
	"menud/connpool"
	"fmt"
)

func main() {
	user, err := connpool.GetUser(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user)
	}
}