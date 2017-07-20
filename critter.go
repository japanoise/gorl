package gorl

/* A creature! */

type Critter struct {
	X       int
	Y       int
	Race    MonsterID
	Name    string
	Collide func(m *Map, g Graphics, this, other *Critter) bool // What to do when there's a collision; true if I should delet it
	Stats   StatBlock
	Female  bool
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
	ret.Collide = func(m *Map, out Graphics, this, other *Critter) bool {
		out.Message("\"Rargh! I'm a very scary monster!\"")
		out.Message("You slap the monster about with a large piece of fish!")
		out.Message("The monster collapses, defeated.")
		return true
	}
	return ret
}
