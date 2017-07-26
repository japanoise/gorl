package gorl

import (
	"fmt"

	"github.com/japanoise/engutil"
)

/* A creature! */

type Critter struct {
	X         int
	Y         int
	Race      MonsterID
	Name      string
	Stats     StatBlock
	Female    bool
	Inv       []*Item
	Gold      int
	Weapon    *Item
	Armor     *Item
	AI        *AiData
	Casting   SpellData
	SpellBook []*Spell
}

type StatBlock struct {
	MaxHp int
	CurHp int
	MaxMp int
	CurMp int
	Level uint8
	Exp   int
}

func (c *Critter) Chase(m *Map, d *DijkstraMap) *Critter {
	lv := d.LowestNeighbour(c.X, c.Y)
	return MoveAbs(m, c, lv.X, lv.Y)
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

func GenStatBlock(hitdice uint8, level uint8) StatBlock {
	hp := 0
	for i := uint8(0); i < level; i++ {
		hp += SmallDiceRoll(hitdice)
	}
	return StatBlock{
		hp, hp, 0, 0, level, 0,
	}
}

func (c *Critter) GetXY() (int, int) {
	return c.X, c.Y
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

func (c *Critter) CanCast(s *Spell) bool {
	ret := true
	if SpellDataHasFlag(s.Data, SpellRitual) {
		ret = ret && SpellDataHasFlag(c.Casting, SpellRitual)
	} else if SpellDataHasFlag(s.Data, SpellSorcery) {
		ret = ret && SpellDataHasFlag(c.Casting, SpellSorcery)
	} else if SpellDataHasFlag(s.Data, SpellHoly) {
		ret = ret && SpellDataHasFlag(c.Casting, SpellHoly)
	}
	return ret
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

func (c *Critter) Kill(state *State) {
	state.Out.Message("You have defeated " + c.GetTheName())
	c.Delete(state.CurLevel)
	for i, crit := range state.Monsters {
		if crit == c {
			state.Monsters[i] = nil
		}
	}
}

func (c *Critter) RestoreHp(hp int) {
	c.Stats.CurHp += hp
	if c.Stats.CurHp > c.Stats.MaxHp {
		c.Stats.CurHp = c.Stats.MaxHp
	}
}
