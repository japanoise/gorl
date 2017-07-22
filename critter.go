package gorl

import "github.com/japanoise/engutil"

/* A creature! */

type Critter struct {
	X      int
	Y      int
	Race   MonsterID
	Name   string
	Stats  StatBlock
	Female bool
	Inv    []*Item
}

type StatBlock struct {
	MaxHp int
	CurHp int
	Str   int
	Dex   int
}

func (c *Critter) DoMove(m *Map, x, y int) {
	if m.Tiles[x][y].IsPassable() {
		c.X = x
		c.Y = y
	}
}

func (c *Critter) Chase(m *Map, x, y int) {
	if x > c.X {
		c.DoMove(m, c.X+1, c.Y)
	} else if x < c.X {
		c.DoMove(m, c.X-1, c.Y)
	} else if y > c.Y {
		c.DoMove(m, c.X, c.Y+1)
	} else if y < c.Y {
		c.DoMove(m, c.X, c.Y-1)
	}
}

func (c *Critter) GetSprite() Sprite {
	if c.Female {
		return Bestiary[c.Race].SprF
	} else {
		return Bestiary[c.Race].SprM
	}
}

func (c *Critter) Delete(m *Map) {
	if m.Tiles[c.X][c.Y].Here == c {
		m.Tiles[c.X][c.Y].Here = nil
	}
}

func DefStatBlock() StatBlock {
	return StatBlock{
		10, 10, 10, 10,
	}
}

func RandomCritter(elevation int) *Critter {
	ret := GetMonster(MonsterUnknown)
	return ret
}

func (c *Critter) GetName() string {
	if c.Name != "" {
		return c.Name
	} else {
		return engutil.ASlashAn(c.GetRaceName())
	}
}

// Returns the critter's name, or "the $RACE" if it's anonymous
func (c *Critter) GetTheName() string {
	if c.Name != "" {
		return c.Name
	} else {
		return "the " + c.GetRaceName()
	}
}

func (c *Critter) GetRaceName() string {
	if c.Race == MonsterHuman {
		if c.Female {
			return "woman"
		} else {
			return "man"
		}
	} else {
		return Bestiary[c.Race].Name
	}
}

func (c *Critter) RollForAttack() int {
	return LargeDiceRoll(1, 20)
}

func (c *Critter) GetDefence() int {
	return 10
}

func (c *Critter) RollForDamage() int {
	return SmallDiceRoll(SmallDice(1, 6))
}

func (c *Critter) TakeDamage(damage int) {
	c.Stats.CurHp -= damage
}

func (c *Critter) IsDead() bool {
	return c.Stats.CurHp <= 0
}
