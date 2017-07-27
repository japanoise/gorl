package gorl

import "strings"

type SpellEffect uint8

const (
	SpellNothing SpellEffect = iota
	SpellFire
	SpellIce
	SpellLightning
	SpellHeal
)

type SpellData uint8

const (
	SpellSelf SpellData = 0x10 << iota
	SpellOther
	SpellArea
)

const (
	SpellSorcery SpellData = 1 << iota
	SpellRitual
	SpellHoly
)

type Spell struct {
	Name    string
	Effect  SpellEffect
	Potency uint8
	Data    SpellData
}

type SpellFunc func(Graphics, *Spell, *Critter, *Map, Point) []*Critter // Given the spell, caster, and location. Returns a list of affected critters.

var SpellFuncs map[SpellEffect]SpellFunc

func SpellDataHasFlag(data, flag SpellData) bool {
	return data&flag != 0
}

func initSpells() error {
	SpellFuncs = make(map[SpellEffect]SpellFunc)
	SpellFuncs[SpellFire] = elemental(SpellFire)
	SpellFuncs[SpellLightning] = elemental(SpellLightning)
	SpellFuncs[SpellIce] = elemental(SpellIce)
	SpellFuncs[SpellHeal] = heal
	return nil
}

// Straightforward damaging spell
func elemental(class SpellEffect) SpellFunc {
	return func(g Graphics, s *Spell, c *Critter, m *Map, p Point) []*Critter {
		if SpellDataHasFlag(s.Data, SpellSelf) {
			c.TakeDamage(SmallDiceRoll(s.Potency))
			return []*Critter{c}
		} else if SpellDataHasFlag(s.Data, SpellOther) {
			x, y := p.GetXY()
			if m.OOB(x, y) {
				g.Message("The " + s.Name + " shoots harmlessly off into the distance.")
				return nil
			} else if m.Tiles[x][y].Here != nil {
				target := m.Tiles[x][y].Here
				g.Message("The " + s.Name + " strikes " + target.GetTheName() + "!")
				target.TakeDamage(SmallDiceRoll(s.Potency))
				return []*Critter{target}
			} else {
				g.Message("The " + s.Name + " strikes " + TilesDir[m.Tiles[x][y].Id].Name + "!")
				return nil
			}
		} else if SpellDataHasFlag(s.Data, SpellArea) {
			x, y := p.GetXY()
			if m.OOB(x, y) {
				g.Message("The " + s.Name + " shoots harmlessly off into the distance.")
				return nil
			}
			g.Message("The " + s.Name + " explodes violently!")
			return forEachCritterInCircle(x, y, int(s.Potency>>4), m,
				func(crit *Critter) {
					crit.TakeDamage(SmallDiceRoll(SmallDice(1, s.Potency)))
				})
		}
		return nil
	}
}

func forEachCritterInCircle(x, y, radius int, m *Map, f func(*Critter)) []*Critter {
	critters := make([]*Critter, 0)
	r2 := radius * radius
	for j := y - radius; j < y+radius; j++ {
		di2 := (j - y) * (j - y)
		for i := x - radius; i < x+radius; i++ {
			// If this point is in the circle…
			if (i-x)*(i-x)+di2 <= r2 {
				_, c := m.GetPassable(i, j)
				if c != nil {
					f(c)
					critters = append(critters, c)
				}
			}
		}
	}
	return critters
}

// Bog-standard
func heal(g Graphics, s *Spell, c *Critter, m *Map, p Point) []*Critter {
	if SpellDataHasFlag(s.Data, SpellSelf) {
		c.RestoreHp(SmallDiceRoll(s.Potency))
		return []*Critter{c}
	} else if SpellDataHasFlag(s.Data, SpellOther) {
		x, y := p.GetXY()
		if m.OOB(x, y) {
			g.Message("The " + s.Name + " shoots off into the distance.")
			return nil
		} else if m.Tiles[x][y].Here != nil {
			target := m.Tiles[x][y].Here
			g.Message("The " + s.Name + " strikes " + target.GetTheName() + "!")
			target.RestoreHp(SmallDiceRoll(s.Potency))
			return []*Critter{target}
		} else {
			g.Message("The " + s.Name + " strikes " + TilesDir[m.Tiles[x][y].Id].Name + "!")
			return nil
		}
	} else if SpellDataHasFlag(s.Data, SpellArea) {
		x, y := p.GetXY()
		if m.OOB(x, y) {
			g.Message("The " + s.Name + " shoots off into the distance.")
			return nil
		}
		g.Message("The " + s.Name + " explodes in a warm ray of light!")
		return forEachCritterInCircle(x, y, int(s.Potency>>4), m,
			func(crit *Critter) {
				crit.RestoreHp(SmallDiceRoll(s.Potency))
			})
	}
	return nil
}

func ZapSpell(state *State, player *Critter, m *Map) []*Critter {
	s, i := PickSpell("Which spell will you cast?", state.Out, player)
	if s == nil {
		state.Out.Message("No spell to cast.")
		return nil
	} else if !player.CanCast(s) {
		state.Out.Message("You cannot cast this spell.")
		return nil
	}
	if SpellDataHasFlag(s.Data, SpellRitual) {
		// Delete the spell from the spellbook
		copy(player.SpellBook[i:], player.SpellBook[i+1:])
		player.SpellBook[len(player.SpellBook)-1] = nil
		player.SpellBook = player.SpellBook[:len(player.SpellBook)-1]
	}
	return DoCastSpell(state, player, m, s)
}

func PickSpell(prompt string, g Graphics, player *Critter) (*Spell, int) {
	if player.SpellBook == nil || player.Casting == 0 {
		return nil, 0
	}
	choices := make([]string, len(player.SpellBook))
	for i, s := range player.SpellBook {
		choices[i] = s.Name
	}
	i := g.MenuIndex(prompt, choices)
	return player.SpellBook[i], i
}

func DoCastSpell(state *State, player *Critter, m *Map, s *Spell) []*Critter {
	p := Point{-1, -1}
	fn := SpellFuncs[s.Effect]
	if fn == nil {
		state.Out.Message("The spell fails.")
		return nil
	}
	if SpellDataHasFlag(s.Data, SpellArea|SpellOther) {
		if state.Dungeon < 0 {
			state.Out.Dungeon(state.CurLevel, player.X, player.Y)
		} else {
			state.Out.Overworld(state.CurLevel, player.X, player.Y)
		}
		d := state.In.GetDirection("Fire " + s.Name + " in which direction?")
		p = getEndPoint(m, player, d)
	}
	if strings.Contains(s.Name, "You") {
		state.Out.Message(s.Name)
	} else {
		state.Out.Message("You cast " + s.Name + "!")
	}
	return fn(state.Out, s, player, m, p)
}

func getEndPoint(m *Map, start Placed, d Direction) Point {
	dx, dy := 0, 0
	switch d {
	case DirEast:
		dx = 1
	case DirWest:
		dx = -1
	case DirSouth:
		dy = 1
	case DirNorth:
		dy = -1
	case DirNE:
		dx = 1
		dy = -1
	case DirSE:
		dx = 1
		dy = 1
	case DirNW:
		dx = -1
		dy = -1
	case DirSW:
		dx = -1
		dy = 1
	}
	x, y := start.GetXY()
	for !m.OOB(x+dx, y+dy) {
		passable, here := m.GetPassable(x+dx, y+dy)
		if here != nil || !passable {
			return Point{x + dx, y + dy}
		}
		if dx > 0 {
			dx++
		} else if dx < 0 {
			dx--
		}
		if dy > 0 {
			dy++
		} else if dy < 0 {
			dy--
		}
	}
	return Point{-1, -1}
}
