package gorl

type TileID byte

const (
	TileVoid TileID = iota
	TileWall
	TileFloor
	TileGrass
	TileStairUp
	TileStairDown
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
	TilesDir[TileVoid] = Tile{
		"the void",
		false,
	}
	TilesDir[TileWall] = Tile{
		"solid wall",
		false,
	}
	TilesDir[TileFloor] = Tile{
		"stone floor",
		true,
	}
	TilesDir[TileGrass] = Tile{
		"grass",
		true,
	}
	TilesDir[TileStairDown] = Tile{
		"stairs leading down",
		true,
	}
	TilesDir[TileStairUp] = Tile{
		"stairs leading up",
		true,
	}
}

/* Fast tile functions for dungen */

// Wall tile
func WallTile() MapTile {
	return MapTile{nil, TileWall}
}

// Floor tile
func FloorTile() MapTile {
	return MapTile{nil, TileFloor}
}
