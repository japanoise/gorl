package main

import (
	"github.com/japanoise/gorl"
)

func main() {
	c := NewCurses()
	gorl.MainLoop(c, c)
}
