package gorl

import (
	"errors"
	"math/rand"
	"time"
)

type StateDungeon struct {
	Seeds    []int64
	Depth    int
	Visited  []bool
	Monsters [][]*Critter
	Items    [][]*DungeonItem
}

type DungeonItem struct {
	X  int
	Y  int
	It *Item
}

type SpawnRegion struct {
	TopLeftX  int
	TopLeftY  int
	BotRightX int
	BotRightY int
}

func initRng() error {
	seed := NewSeed()
	rand.Seed(seed)
	debug.Print(seed)
	return nil
}

func DigDungeon(d int) *StateDungeon {
	ret := StateDungeon{
		make([]int64, d+1),
		d,
		make([]bool, d+1),
		make([][]*Critter, d+1),
		make([][]*DungeonItem, d+1),
	}
	return &ret
}

func (d *StateDungeon) GetDunLevel(oldelevation, elevation int, monlist []*Critter, items []*DungeonItem) (*Map, []*Critter, []SpawnRegion, error) {
	debug.Print("Called get dunlevel with args", oldelevation, elevation, monlist)
	d.Visited[oldelevation] = true
	d.Monsters[oldelevation] = make([]*Critter, len(monlist))
	copy(d.Monsters[oldelevation], monlist)
	debug.Print("Copy d.Monsters, monlist: ", d.Monsters[oldelevation], monlist)
	if oldelevation != 0 && items != nil {
		d.Items[oldelevation] = make([]*DungeonItem, len(items))
		copy(d.Items[oldelevation], items)
	}
	if elevation <= 0 || elevation > d.Depth {
		return nil, []*Critter{}, []SpawnRegion{}, errors.New("Outside of dungeon range")
	}
	if !d.Visited[elevation] {
		d.Seeds[elevation] = NewSeed()
	}
	m, sp := d.GenLevel(elevation)
	if !d.Visited[elevation] {
		d.Monsters[elevation] = Populate(m, sp, elevation)
		d.Visited[elevation] = true
	} else {
		for _, mon := range d.Monsters[elevation] {
			if mon != nil {
				m.Tiles[mon.X][mon.Y].Here = mon
			}
		}
		m.PlaceItems(d.Items[elevation])
	}
	return m, d.Monsters[elevation], sp, nil
}

func (d *StateDungeon) GetDunLevelFromStorage(elevation int) (*Map, []*Critter) {
	m, _ := d.GenLevel(elevation)
	for _, mon := range d.Monsters[elevation] {
		if mon != nil {
			m.Tiles[mon.X][mon.Y].Here = mon
		}
	}
	m.PlaceItems(d.Items[elevation])
	return m, d.Monsters[elevation]
}

func NewSeed() int64 {
	return time.Now().UTC().UnixNano()
}

func (d *StateDungeon) GenLevel(elevation int) (*Map, []SpawnRegion) {
	r := rand.New(rand.NewSource(d.Seeds[elevation]))
	m, sp := dunGen(elevation, r)
	dunAddFeatures(m, sp, elevation, d.Depth, r)
	return m, sp
}

func dunGen(elevation int, r *rand.Rand) (*Map, []SpawnRegion) {
	sizex, sizey := 100, 100
	retval := GetBlankMap(elevation, sizex, sizey)
	numroomsx := 5
	numroomsy := 5
	roomsx := sizex / numroomsx
	roomsy := sizey / numroomsy
	xh := roomsx / 2
	yh := roomsy / 2
	spawns := make([]SpawnRegion, 0)
	// Draw rooms
	for xrn := 0; xrn < numroomsx; xrn++ {
		for yrn := 0; yrn < numroomsy; yrn++ {
			anchorx := xrn * roomsx
			anchory := yrn * roomsy
			var posx, posy, cornerx, cornery int
			if r.Intn(3) == 0 {
				// Build an intersection
				posx = anchorx + xh - 1
				posy = anchory + yh - 1
				cornerx = posx + 2
				cornery = posy + 2
			} else {
				// Build a room and add it to the spawn list
				posx = anchorx + r.Intn(roomsx/3)
				posy = anchory + r.Intn(roomsy/3)
				cornerx = (anchorx + roomsx - 2) - r.Intn(roomsx/3)
				cornery = (anchory + roomsy - 2) - r.Intn(roomsy/3)
				spawns = append(spawns, SpawnRegion{
					posx + 1, posy + 1,
					cornerx - 1, cornery - 1,
				})
			}
			for x := posx; x <= cornerx; x++ {
				for y := posy; y <= cornery; y++ {
					if x == posx || y == posy || x == cornerx || y == cornery {
						retval.Tiles[x][y] = WallTile()
					} else {
						retval.Tiles[x][y] = FloorTile()
					}
				}
			}
		}
	}
	// If there aren't enough rooms, there might be hangs when placing things later.
	if len(spawns) < 5 {
		return dunGen(elevation, r)
	}
	// Draw Corridoors
	for xrn := 0; xrn < numroomsx; xrn++ {
		for yrn := 0; yrn < numroomsy; yrn++ {
			anchorx := xrn * roomsx
			anchory := yrn * roomsy
			if xrn != numroomsx-1 {
				y := anchory + yh
				for x := anchorx + xh; x <= anchorx+xh+roomsx && x < sizex; x++ {
					if !retval.Tiles[x][y].IsPassable() {
						if retval.Tiles[x][y+1].Id == TileWall && x != anchorx+xh+1 && x < anchorx+xh+roomsx-1 {
							retval.Tiles[x][y] = DoorTile()
						} else {
							retval.Tiles[x][y+1] = WallTile()
							retval.Tiles[x][y-1] = WallTile()
							retval.Tiles[x][y] = FloorTile()
						}
					}
				}
			}
			if yrn != numroomsy-1 {
				x := anchorx + xh
				for y := anchory + yh; y <= anchory+yh+roomsy && y < sizey; y++ {
					if !retval.Tiles[x][y].IsPassable() {
						if retval.Tiles[x+1][y].Id == TileWall && y != anchory+yh+1 && y < anchory+yh+roomsy-1 {
							retval.Tiles[x][y] = DoorTile()
						} else {
							retval.Tiles[x+1][y] = WallTile()
							retval.Tiles[x-1][y] = WallTile()
							retval.Tiles[x][y] = FloorTile()
						}
					}
				}
			}
		}
	}
	return retval, spawns
}

// Adds features to the dungeon
func dunAddFeatures(m *Map, spawnrooms []SpawnRegion, elevation, maxdepth int, r *rand.Rand) {
	// Add stairs leading up
	room := spawnrooms[r.Intn(len(spawnrooms))]
	x, y := FindPlaceInRoom(m, room)
	m.Tiles[x][y].Id = TileStairUp
	m.UpStairsX = x
	m.UpStairsY = y
	// Add stairs leading down
	if elevation != maxdepth {
		room := spawnrooms[r.Intn(len(spawnrooms))]
		x, y := FindPlaceInRoom(m, room)
		m.Tiles[x][y].Id = TileStairDown
		m.DownStairsX = x
		m.DownStairsY = y
	}

	for _, room := range spawnrooms {
		// Extra features
		roll := r.Intn(40) - 19
		if roll > 15 {
			x, y := FindPlaceInRoom(m, room)
			switch roll {
			case 16:
				m.Tiles[x][y].Id = TileColumn
			case 17:
				m.Tiles[x][y].Id = TileFountain
			case 18:
				m.Tiles[x][y].Id = TileAltar
			default:
				m.Tiles[x][y].Id = TileSprungTrap
			}
		}
	}
}

// Populate our beautiful dungeon with treasure and monsters!!!
func Populate(dungeon *Map, spawnrooms []SpawnRegion, elevation int) []*Critter {
	ret := make([]*Critter, len(spawnrooms))
	for i, room := range spawnrooms {
		mons := RandomCritter(dungeon.Elevation)
		ret[i] = mons
		PlaceCritterInRoom(mons, dungeon, room)
		PlaceItemInRoom(GetGoldCoins(LargeDiceRoll(elevation, 6)), dungeon, room)
	}
	return ret
}

func FindPlaceInRoom(dungeon *Map, room SpawnRegion) (int, int) {
	roomw := room.BotRightX - room.TopLeftX
	roomh := room.BotRightY - room.TopLeftY
	x := room.TopLeftX + rand.Intn(roomw)
	y := room.TopLeftY + rand.Intn(roomh)
	for pass, here := dungeon.GetPassable(x, y); !(notStairs(dungeon, x, y) && pass && here == nil); {
		x = room.TopLeftX + rand.Intn(roomw)
		y = room.TopLeftY + rand.Intn(roomh)
	}
	return x, y
}

func PlaceCritterInRoom(mons *Critter, dungeon *Map, room SpawnRegion) {
	if dungeon.Tiles[mons.X][mons.Y].Here == mons {
		dungeon.Tiles[mons.X][mons.Y].Here = nil
	}
	mons.X, mons.Y = FindPlaceInRoom(dungeon, room)
	dungeon.Tiles[mons.X][mons.Y].Here = mons
}

func PlaceItemInRoom(item *Item, dungeon *Map, room SpawnRegion) {
	x, y := FindPlaceInRoom(dungeon, room)
	if dungeon.Tiles[x][y].Items == nil {
		dungeon.Tiles[x][y].Items = []*Item{item}
	} else {
		dungeon.Tiles[x][y].Items = append(dungeon.Tiles[x][y].Items, item)
	}
}

func notStairs(dungeon *Map, x, y int) bool {
	return dungeon.Tiles[x][y].Id != TileStairDown && dungeon.Tiles[x][y].Id != TileStairUp
}

// Place the critter at a random point in the dungeon
func PlaceCritter(mons *Critter, dungeon *Map, spawnrooms []SpawnRegion) {
	room := spawnrooms[rand.Intn(len(spawnrooms))]
	PlaceCritterInRoom(mons, dungeon, room)
}
