package gorl

type Map struct {
	Tiles       [][]MapTile
	SizeX       int
	SizeY       int
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
}

func (m *MapTile) IsPassable() bool {
	return TilesDir[m.Id].Passable
}

func Move(m *Map, who *Critter, dx, dy int) *Critter {
	x, y := who.X+dx, who.Y+dy
	passable, target := m.GetPassable(x, y)
	if passable && target == nil {
		m.Tiles[who.X][who.Y].Here = nil
		who.X = x
		who.Y = y
		m.Tiles[who.X][who.Y].Here = who
	}
	return target
}

func (m *Map) GetPassable(x, y int) (bool, *Critter) {
	if x < m.SizeX && x >= 0 && y < m.SizeY && y >= 0 {
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
