package main

import (
	"fmt"

	"github.com/japanoise/gorl"
)

func main() {
	c := NewCurses()
	err := gorl.StartGame(c, c)
	if err != nil {
		c.End()
		fmt.Println(err.Error())
	}
}
