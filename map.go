package gorl

import (
	"math"
)

type Map struct {
	Tiles       [][]MapTile
	MSizeX      int
	MSizeY      int
	Elevation   int // gt 0 means dungeon, lt 0 means outdoors, 0 means inside a civilized area
	UpStairsX   int
	UpStairsY   int
	DownStairsX int
	DownStairsY int
}

type MapTile struct {
	Here   *Critter
	Id     TileID
	Items  []*Item
	OwData *MapOverworldData
	Lit    bool
	Disc   bool
}

func (m *MapTile) IsPassable() bool {
	return TilesDir[m.Id].Passable
}

func (m *Map) SizeY() int {
	return m.MSizeY
}

func (m *Map) SizeX() int {
	return m.MSizeX
}

func (m *Map) IsPassable(x, y int) bool {
	return m.Tiles[x][y].IsPassable()
}

func (m *Map) CollectItems(items []*DungeonItem) []*DungeonItem {
	for x := 0; x < m.MSizeX; x++ {
		for y := 0; y < m.MSizeY; y++ {
			if m.Tiles[x][y].IsPassable() {
				if m.Tiles[x][y].Items != nil {
					for _, item := range m.Tiles[x][y].Items {
						items = append(items, &DungeonItem{x, y, item})
					}
				}
			}
		}
	}
	return items
}

func (m *Map) PlaceItems(items []*DungeonItem) {
	for _, item := range items {
		if m.Tiles[item.X][item.Y].Items == nil {
			m.Tiles[item.X][item.Y].Items = make([]*Item, 0, 3)
		}
		m.Tiles[item.X][item.Y].Items = append(m.Tiles[item.X][item.Y].Items, item.It)
	}
}

func Move(m *Map, who *Critter, dx, dy int) *Critter {
	x, y := who.X+dx, who.Y+dy
	return MoveAbs(m, who, x, y)
}

func MoveAbs(m *Map, who *Critter, x, y int) *Critter {
	passable, target := m.GetPassable(x, y)
	if passable && target == nil {
		m.Tiles[who.X][who.Y].Here = nil
		who.X = x
		who.Y = y
		m.Tiles[who.X][who.Y].Here = who
	}
	return target
}

func (m *Map) CantSeeThrough(x, y int) bool {
	if !m.OOB(x, y) {
		return TilesDir[m.Tiles[x][y].Id].Transparent == false
	} else {
		return true
	}
}

func (m *Map) Lit(x, y int) {
	if !m.OOB(x, y) {
		m.Tiles[x][y].Lit = true
		m.Tiles[x][y].Disc = true
	}
}

func (m *Map) UnLit(x, y int) {
	if !m.OOB(x, y) {
		m.Tiles[x][y].Lit = false
	}
}

// Is a point out-of-bounds?
func (m *Map) OOB(x, y int) bool {
	if x < m.MSizeX && x >= 0 && y < m.MSizeY && y >= 0 {
		return false
	} else {
		return true
	}
}

func (m *Map) GetPassable(x, y int) (bool, *Critter) {
	if !m.OOB(x, y) {
		return m.Tiles[x][y].IsPassable(), m.Tiles[x][y].Here
	} else {
		return false, nil
	}
}

func (m *Map) PlaceCritterAtDownStairs(c *Critter) {
	c.X = m.DownStairsX
	c.Y = m.DownStairsY
	m.Tiles[c.X][c.Y].Here = c
}

func (m *Map) PlaceCritterAtUpStairs(c *Critter) {
	c.X = m.UpStairsX
	c.Y = m.UpStairsY
	m.Tiles[c.X][c.Y].Here = c
}

func (m *Map) Darken() {
	for x := 0; x < m.MSizeX; x++ {
		for y := 0; y < m.MSizeY; y++ {
			m.Tiles[x][y].Lit = false
		}
	}
}

func GetBlankMap(elevation, sizex, sizey int) *Map {
	retval := Map{
		make([][]MapTile, sizex),
		sizex,
		sizey,
		elevation,
		0, 0, 0, 0,
	}
	for i := 0; i < sizex; i++ {
		retval.Tiles[i] = make([]MapTile, sizey)
		for j := 0; j < sizey; j++ {
			retval.Tiles[i][j] = MapTile{}
		}
	}
	return &retval
}

func CalcVisibility(m *Map, player *Critter, light int) {
	clearlight(m, player, light)
	fov(m, player.X, player.Y, light)
}

func clearlight(m *Map, player *Critter, light int) {
	for x := player.X - light - 1; x < player.X+light+1; x++ {
		for y := player.Y - light - 1; y < player.Y+light+1; y++ {
			m.UnLit(x, y)
		}
	}
}

func fov(m *Map, x, y int, radius int) {
	for i := -radius; i <= radius; i++ { //iterate out of map bounds as well (radius^1)
		for j := -radius; j <= radius; j++ { //(radius^2)
			if i*i+j*j < radius*radius {
				los(m, x, y, x+i, y+j)
			}
		}
	}
}

/* Los calculation http://www.roguebasin.com/index.php?title=LOS_using_strict_definition */
func los(m *Map, x0, y0, x1, y1 int) {
	// By taking source by reference, litting can be done outside of this function which would be better made generic.
	var sx int
	var sy int
	var dx int
	var dy int
	var dist float64

	dx = x1 - x0
	dy = y1 - y0

	//determine which quadrant to we're calculating: we climb in these two directions
	if x0 < x1 { //sx = (x0 < x1) ? 1 : -1;
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 { //sy = (y0 < y1) ? 1 : -1;
		sy = 1
	} else {
		sy = -1
	}

	xnext := x0
	ynext := y0

	//calculate length of line to cast (distance from start to final tile)
	dist = sqrt(dx*dx + dy*dy)

	for xnext != x1 || ynext != y1 { //essentially casting a ray of length radius: (radius^3)
		if m.OOB(xnext, ynext) {
			return
		}
		if m.CantSeeThrough(xnext, ynext) {
			m.Tiles[xnext][ynext].Disc = true
			return
		}

		// Line-to-point distance formula < 0.5
		if abs(dy*(xnext-x0+sx)-dx*(ynext-y0))/dist < 0.5 {
			xnext += sx
		} else if abs(dy*(xnext-x0)-dx*(ynext-y0+sy))/dist < 0.5 {
			ynext += sy
		} else {
			xnext += sx
			ynext += sy
		}
	}
	m.Lit(x1, y1)
	if !m.OOB(x1, y1) && m.Tiles[x1][y1].Here != nil && m.Tiles[x1][y1].Here.AI != nil {
		m.Tiles[x1][y1].Here.AI.Active = true
	}
}

func sqrt(x int) float64 {
	return math.Sqrt(float64(x))
}

func abs(x int) float64 {
	return math.Abs(float64(x))
}
