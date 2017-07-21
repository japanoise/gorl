package gorl

type Overworld struct {
	M       *Map
	SavedPx int
	SavedPy int
}

type MapOverworldData struct {
	Dungeon *StateDungeon
}

func OverworldGen(player *Critter, px, py int) *Overworld {
	sizex, sizey := 100, 100
	m := GetBlankMap(-1, sizex, sizey)
	for x := 0; x < sizex; x++ {
		for y := 0; y < sizey; y++ {
			if x < 5 || y < 5 || x > sizex-5 || y > sizex-5 {
				m.Tiles[x][y] = MapTile{nil, TileSea, []*Item{}, nil}
			} else {
				m.SetGrassTile(x, y)
			}
		}
	}
	m.Tiles[10][10] = MapTile{nil, TileOverworldDungeon, []*Item{}, &MapOverworldData{DigDungeon(5)}}
	m.Tiles[15][10] = MapTile{nil, TileOverworldDungeon, []*Item{}, &MapOverworldData{DigDungeon(5)}}
	player.X = px
	player.Y = py
	m.Tiles[px][py].Here = player
	return &Overworld{m, px, py}
}

func (m *Map) SetGrassTile(x, y int) {
	if x%2 == 0 && y%2 == 0 {
		m.Tiles[x][y] = MapTile{nil, TileGrass2, []*Item{}, nil}
	} else {
		m.Tiles[x][y] = MapTile{nil, TileGrass, []*Item{}, nil}
	}
}
