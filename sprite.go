package gorl

type Sprite uint64

// All sprites a tileset or graphics implementation will have to deal with.
// Tiles are done seperately.
const (
	SpriteBlank Sprite = iota
	SpriteHumanMale
	SpriteHumanFemale
	SpriteMonsterUnknown
	SpriteItemGold
	SpriteItemWeaponGeneric
	SpriteItemAppGeneric
)
