package gorl

import (
	"fmt"
	"math/rand"
	"strconv"
)

func SmallDice(n, y uint8) uint8 {
	return (n << 4) | (y & 0x0E)
}

func SmallDiceRoll(ndy uint8) int {
	n := ndy >> 4
	y := ndy & 0x0E
	ret := 0
	var i uint8
	for i = 0; i < n; i++ {
		ret += rand.Intn(int(y)) + 1
	}
	return ret
}

func GetSmallDiceString(ndy uint8) string {
	return strconv.Itoa(int(ndy>>4)) + "d" + strconv.Itoa(int(ndy&0x0E))
}

func LargeDiceRoll(n, y int) int {
	ret := 0
	for i := 0; i < n; i++ {
		ret += rand.Intn(y) + 1
	}
	return ret
}

func Attack(interactive, pattacker bool, m *Map, g Graphics, attacker, defender *Critter) bool {
	attackerRoll := attacker.RollForAttack()
	defenderDef := defender.GetDefence()
	dead := false
	if attackerRoll >= defenderDef {
		damage := attacker.RollForDamage()
		defender.TakeDamage(damage)
		if defender.IsDead() {
			dead = true
		}
		if interactive && pattacker {
			g.Message(AttackMessage("You", defender.GetTheName(), attacker, defender))
		} else if interactive && !pattacker {
			g.Message(AttackMessage(attacker.GetTheName(), "you", attacker, defender))
		}
	} else if interactive {
		if pattacker {
			g.Message(MissMessage("You", defender.GetTheName(), attacker, defender))
		} else {
			g.Message(MissMessage(attacker.GetTheName(), "you", attacker, defender))
		}
	}
	return dead
}

func AttackMessage(aname, dname string, attacker, defender *Critter) string {
	if aname == "You" {
		return fmt.Sprintf("You strike %s", dname)
	} else {
		return fmt.Sprintf("%s strikes %s", aname, dname)
	}
}

func MissMessage(aname, dname string, attacker, defender *Critter) string {
	if aname == "You" {
		return fmt.Sprintf("You miss %s", dname)
	}
	return fmt.Sprintf("%s misses %s", aname, dname)
}
