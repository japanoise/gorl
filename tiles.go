package gorl

type TileID byte

const (
	TileVoid TileID = iota
	TileWall
	TileFloor
	TileGrass
	TileGrass2
	TileStairUp
	TileStairDown
	TileSea
	TileFreshwater
	TileOverworldDungeon
	TileOcean
	TileMountain
	TileOverworldVillage
	TileLetterI
	TileLetterN
	TileDoor
	TileDoorOpen
	TileColumn
	TileFountain
	TileAltar
	TileSprungTrap
)

type Tile struct {
	Name        string
	Passable    bool
	Transparent bool
}

var (
	TilesDir map[TileID]Tile
)

func initTiles() error {
	/* Basic Tiles */
	TilesDir = make(map[TileID]Tile)
	return loadConfigFile("tiles.json", &TilesDir)
}

func Look(m *Map, out Graphics, in Input, player *Critter) {
	tile := TileNSquaresInDirFromXY(m, 1, in.GetDirection("Look in which direction?"), player.X, player.Y)
	if tile == nil {
		out.Message("There's nothing here.")
		return
	}
	ret := ""
	if tile.Here != nil {
		ret += tile.Here.GetName() + ", "
	}
	if len(tile.Items) > 0 {
		for _, item := range tile.Items {
			ret += item.Describe() + ", "
		}
	}
	ret += TilesDir[tile.Id].Name
	out.Message(ret)
}

func TileNSquaresInDirFromXY(m *Map, squares int, dir Direction, x, y int) *MapTile {
	switch dir {
	case DirNorth:
		if y-squares < 0 {
			return nil
		} else {
			return &m.Tiles[x][y-squares]
		}
	case DirSouth:
		if y+squares >= m.SizeY {
			return nil
		} else {
			return &m.Tiles[x][y+squares]
		}
	case DirWest:
		if x-squares < 0 {
			return nil
		} else {
			return &m.Tiles[x-squares][y]
		}
	case DirEast:
		if x+squares >= m.SizeX {
			return nil
		} else {
			return &m.Tiles[x+squares][y]
		}
	case DirNE:
		if x+squares >= m.SizeX || y-squares < 0 {
			return nil
		} else {
			return &m.Tiles[x+squares][y-squares]
		}
	case DirSE:
		if x+squares >= m.SizeX || y+squares >= m.SizeY {
			return nil
		} else {
			return &m.Tiles[x+squares][y+squares]
		}
	case DirNW:
		if x-squares < 0 || y-squares < 0 {
			return nil
		} else {
			return &m.Tiles[x-squares][y-squares]
		}
	case DirSW:
		if x-squares < 0 || y+squares >= m.SizeY {
			return nil
		} else {
			return &m.Tiles[x-squares][y+squares]
		}
	case DirUp:
		return &m.Tiles[x][y]
	}
	return nil
}

/* Fast tile functions for dungen */

// Wall tile
func WallTile() MapTile {
	return TileOfClass(TileWall)
}

// Floor tile
func FloorTile() MapTile {
	return TileOfClass(TileFloor)
}

// Door tile
func DoorTile() MapTile {
	return TileOfClass(TileDoor)
}

func TileOfClass(id TileID) MapTile {
	return MapTile{nil, id, nil, nil, false, false}
}
