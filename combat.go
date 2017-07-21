package gorl

import "math/rand"

func Dice(n, y uint8) uint8 {
	return (n << 4) | (y & 0x0E)
}

func DiceRoll(ndy uint8) int {
	n := ndy >> 4
	y := ndy & 0x0E
	ret := 0
	var i uint8
	for i = 0; i < n; i++ {
		ret += rand.Intn(int(y)) + 1
	}
	return ret
}
