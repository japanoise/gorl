package gorl

import (
	"math"
	"math/rand"

	ds "github.com/xackery/diamondsquare"
)

type Overworld struct {
	M       *Map
	MetaOw  *Map
	SavedPx int
	SavedPy int
	MetaPx  int
	MetaPy  int
}

type MapOverworldData struct {
	Dungeon *StateDungeon
	OWSeed  int64
}

func OverworldGen(player *Critter, px, py int) *Overworld {
	mow, mpx, mpy := diamondSquare(7)
	mow.Tiles[mpx][mpy].Here = player
	m := mow.MetaOWGenMap(mpx, mpy)
	player.X = px
	player.Y = py
	m.Tiles[px][py].Here = player
	return &Overworld{m, mow, px, py, mpx, mpy}
}

func (mow *Map) MetaOWGenMap(mpx, mpy int) *Map {
	var seed int64
	if mow.Tiles[mpx][mpy].OwData == nil {
		seed = NewSeed()
		mow.Tiles[mpx][mpy].OwData = &MapOverworldData{nil, NewSeed()}
	} else {
		seed = mow.Tiles[mpx][mpy].OwData.OWSeed
	}
	return GenOWMap(seed, mow.Tiles[mpx][mpy-1].Id,
		mow.Tiles[mpx+1][mpy].Id,
		mow.Tiles[mpx-1][mpy].Id,
		mow.Tiles[mpx][mpy+1].Id) // For now, we don't need to bounds check because there's no passable tiles on the edge.
}

func GenOWMap(seed int64, north, east, west, south TileID) *Map {
	// Set up the easy stuff
	sizex, sizey := 100, 100
	m := GetBlankMap(-1, sizex, sizey)
	r := rand.New(rand.NewSource(seed))

	// Calculate the passable area
	passn, passw := 0, 0
	passs, passe := sizey, sizex
	if north == TileSea || north == TileMountain {
		passn = 10
	}
	if east == TileSea || east == TileMountain {
		passe = sizex - 10
	}
	if west == TileSea || west == TileMountain {
		passw = 10
	}
	if south == TileSea || south == TileMountain {
		passs = sizey - 10
	}

	// Draw the terrain
	for x := 0; x < sizex; x++ {
		for y := 0; y < sizey; y++ {
			if x > passe {
				if east == TileSea {
					m.Tiles[x][y] = MapTile{nil, TileSea, []*Item{}, nil, false, false}
				} else {
					m.Tiles[x][y] = MapTile{nil, TileMountain, []*Item{}, nil, false, false}
				}
			} else if y < passn {
				if north == TileSea {
					m.Tiles[x][y] = MapTile{nil, TileSea, []*Item{}, nil, false, false}
				} else {
					m.Tiles[x][y] = MapTile{nil, TileMountain, []*Item{}, nil, false, false}
				}
			} else if y > passs {
				if south == TileSea {
					m.Tiles[x][y] = MapTile{nil, TileSea, []*Item{}, nil, false, false}
				} else {
					m.Tiles[x][y] = MapTile{nil, TileMountain, []*Item{}, nil, false, false}
				}
			} else if x < passw {
				if west == TileSea {
					m.Tiles[x][y] = MapTile{nil, TileSea, []*Item{}, nil, false, false}
				} else {
					m.Tiles[x][y] = MapTile{nil, TileMountain, []*Item{}, nil, false, false}
				}
			} else {
				m.SetGrassTile(x, y)
			}
		}
	}

	// Place the dungeons
	ranx, rany := passe-passw, passs-passn
	duns := 5 + r.Intn(5)
	for i := 0; i < duns; i++ {
		looping := true
		for looping {
			rx, ry := passw+r.Intn(ranx), passn+r.Intn(rany)
			if m.Tiles[rx][ry].Id == TileGrass || m.Tiles[rx][ry].Id == TileGrass2 {
				m.Tiles[rx][ry] = MapTile{nil, TileOverworldDungeon, nil, &MapOverworldData{DigDungeon(5 + r.Intn(5)), 0}, false, false}
				looping = false
			}
		}
	}
	return m
}

func (m *Map) SetGrassTile(x, y int) {
	if x%2 == 0 && y%2 == 0 {
		m.Tiles[x][y] = MapTile{nil, TileGrass2, []*Item{}, nil, false, false}
	} else {
		m.Tiles[x][y] = MapTile{nil, TileGrass, []*Item{}, nil, false, false}
	}
}

func diamondSquare(size int) (*Map, int, int) {
	sizeof := int(math.Pow(2, float64(size))) + 1
	ret := GetBlankMap(-1, sizeof, sizeof)
	data, err := ds.Generate(sizeof, int64(sizeof), 10)
	if err != nil {
		debug.Println(err.Error())
	}
	max := 0.0
	min := math.MaxFloat64
	for x := 0; x < sizeof; x++ {
		for y := 0; y < sizeof; y++ {
			if data[x][y] > max {
				max = data[x][y]
			}
			if data[x][y] < min {
				min = data[x][y]
			}
		}
	}
	debug.Println("Max:", max, "min:", min)
	for x := 0; x < sizeof; x++ {
		for y := 0; y < sizeof; y++ {
			ret.dsTile(data[x][y], max, min, x, y)
		}
	}
	for {
		// Pretty dumb, but it gets us somewhere safe to spawn
		px, py := rand.Intn(sizeof), rand.Intn(sizeof)
		if pass, _ := ret.GetPassable(px, py); pass {
			return ret, px, py
		}
	}
}

func (m *Map) dsTile(v, max, min float64, x, y int) {
	ran := max - min
	valinrange := v - min
	pc := (valinrange / ran) * 100
	if pc >= 75 {
		if pc >= 90 {
			m.Tiles[x][y] = MapTile{nil, TileMountain, []*Item{}, nil, false, false}
		} else if pc >= 78 {
			m.SetGrassTile(x, y)
		} else {
			m.Tiles[x][y] = MapTile{nil, TileSea, []*Item{}, nil, false, false}
		}
	} else {
		m.Tiles[x][y] = MapTile{nil, TileOcean, []*Item{}, nil, false, false}
	}
}
