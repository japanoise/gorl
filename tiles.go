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
)

type Tile struct {
	Name     string
	Passable bool
}

var (
	TilesDir map[TileID]Tile
)

func init() {
	/* Basic Tiles */
	TilesDir = make(map[TileID]Tile)
	TilesDir[TileVoid] = Tile{"the void", false}
	TilesDir[TileWall] = Tile{"solid wall", false}
	TilesDir[TileFloor] = Tile{"stone floor", true}
	TilesDir[TileGrass] = Tile{"grass", true}
	TilesDir[TileGrass2] = Tile{"grass", true}
	TilesDir[TileStairDown] = Tile{"stairs leading down", true}
	TilesDir[TileStairUp] = Tile{"stairs leading up", true}
	TilesDir[TileSea] = Tile{"seawater", false}
	TilesDir[TileFreshwater] = Tile{"freshwater", false}
	TilesDir[TileOverworldDungeon] = Tile{"cave entrance", true}
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
	return MapTile{nil, TileWall, []*Item{}, nil}
}

// Floor tile
func FloorTile() MapTile {
	return MapTile{nil, TileFloor, []*Item{}, nil}
}
