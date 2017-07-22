package gorl

import (
	"fmt"

	"github.com/japanoise/engutil"
)

/* A creature! */

type Critter struct {
	X      int
	Y      int
	Race   MonsterID
	Name   string
	Stats  StatBlock
	Female bool
	Inv    []*Item
	Gold   int
	Weapon *Item
	Armor  *Item
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
	if m.Tiles[c.X][c.Y].Items == nil {
		m.Tiles[c.X][c.Y].Items = make([]*Item, 0, len(c.Inv))
	}
	if c.Inv != nil {
		m.Tiles[c.X][c.Y].Items = append(m.Tiles[c.X][c.Y].Items, c.Inv...)
	}
	if c.Gold > 0 {
		m.Tiles[c.X][c.Y].Items = append(m.Tiles[c.X][c.Y].Items, GetGoldCoins(c.Gold))
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

func (c *Critter) GetRace() Monster {
	return Bestiary[c.Race]
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
		return c.GetRace().Name
	}
}

func (c *Critter) RollForAttack() int {
	return LargeDiceRoll(1, 20)
}

func (c *Critter) GetDefence() int {
	if c.Armor == nil {
		return int(c.GetRace().BaseAC)
	} else {
		return int(c.Armor.GetAC())
	}
}

func (c *Critter) RollForDamage() int {
	if c.Weapon == nil {
		return SmallDiceRoll(c.GetRace().BaseDamage)
	} else {
		return c.Weapon.DoDamage()
	}
}

func (c *Critter) TakeDamage(damage int) {
	c.Stats.CurHp -= damage
}

func (c *Critter) IsDead() bool {
	return c.Stats.CurHp <= 0
}

func (c *Critter) SnarfItems(items []*Item) {
	for _, item := range items {
		if item.Class == ItemClassCurrency {
			c.Gold += item.Value
		} else {
			c.Inv = append(c.Inv, item)
		}
	}
}

func (c *Critter) CompleteDescription(g Graphics) {
	wield := "Wielding nothing."
	if c.Weapon != nil {
		wield = "Wielding " + c.Weapon.DescribeExtra()
	}
	wear := "Buck naked."
	if c.Armor != nil {
		wear = "Wearing " + c.Armor.DescribeExtra()
	}
	g.LongMessage(
		fmt.Sprintf("%s the %s %s", c.Name, GetMaleFemaleStr(c.Female), c.GetRace().Name),
		wield,
		wear,
	)
}

func GetMaleFemaleStr(female bool) string {
	if female {
		return "female"
	} else {
		return "male"
	}
}
