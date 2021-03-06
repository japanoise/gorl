package gorl

import (
	"fmt"
	"math/rand"
	"strconv"
)

const SmallDiceYMask = 0x0F

func SmallDice(n, y uint8) uint8 {
	return (n << 4) | (y & SmallDiceYMask)
}

func SmallDiceRoll(ndy uint8) int {
	n := ndy >> 4
	y := ndy & SmallDiceYMask
	ret := 0
	var i uint8
	for i = 0; i < n; i++ {
		ret += rand.Intn(int(y)) + 1
	}
	return ret
}

func GetSmallDiceString(ndy uint8) string {
	return strconv.Itoa(int(ndy>>4)) + "d" + strconv.Itoa(int(ndy&SmallDiceYMask))
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

func RangedAttack(interactive, pattacker bool, s *State, attacker, defender *Critter) bool {
	attackerRoll := attacker.RollForRangedAttack()
	defenderDef := defender.GetDefence()
	dead := false
	if attackerRoll >= defenderDef {
		damage := attacker.RollForRangedDamage()
		defender.TakeDamage(damage)
		if defender.IsDead() {
			dead = true
		}
		if interactive && pattacker {
			s.Out.Message(RangedAttackMessage("You", defender.GetTheName(), attacker, defender))
		} else if interactive && !pattacker {
			s.Out.Message(RangedAttackMessage(attacker.GetTheName(), "you", attacker, defender))
		}
	} else if interactive {
		if pattacker {
			s.Out.Message(MissMessage("You", defender.GetTheName(), attacker, defender))
		} else {
			s.Out.Message(MissMessage(attacker.GetTheName(), "you", attacker, defender))
		}
	}
	return dead
}

func RangedAttackMessage(aname, dname string, attacker, defender *Critter) string {
	if aname == "You" {
		return fmt.Sprintf("Your shot strikes %s", dname)
	} else {
		return fmt.Sprintf("%s's shot strikes %s", aname, dname)
	}
}
