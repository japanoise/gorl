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
