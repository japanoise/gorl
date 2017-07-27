package main

import (
	"github.com/japanoise/gorl"
	"github.com/nsf/termbox-go"
)

func getTileSprites() map[gorl.TileID]*CursesSprite {
	ret := make(map[gorl.TileID]*CursesSprite)
	ret[gorl.TileVoid] = &CursesSprite{' ', termbox.ColorDefault, termbox.ColorDefault}
	ret[gorl.TileWall] = &CursesSprite{'#', termbox.ColorDefault, termbox.ColorDefault}
	ret[gorl.TileFloor] = &CursesSprite{'.', termbox.ColorDefault, termbox.ColorDefault}
	ret[gorl.TileGrass] = &CursesSprite{',', termbox.ColorGreen, termbox.ColorDefault}
	ret[gorl.TileGrass2] = &CursesSprite{'`', termbox.ColorGreen, termbox.ColorDefault}
	ret[gorl.TileStairUp] = &CursesSprite{'<', termbox.ColorDefault, termbox.ColorDefault}
	ret[gorl.TileStairDown] = &CursesSprite{'>', termbox.ColorDefault, termbox.ColorDefault}
	ret[gorl.TileSea] = &CursesSprite{'~', termbox.ColorBlue | termbox.AttrBold, termbox.ColorDefault}
	ret[gorl.TileFreshwater] = &CursesSprite{'~', termbox.ColorCyan | termbox.AttrReverse, termbox.ColorDefault}
	ret[gorl.TileOverworldDungeon] = &CursesSprite{'>', termbox.ColorDefault | termbox.AttrReverse, termbox.ColorDefault}
	return ret
}

func getSprites() map[gorl.Sprite]*CursesSprite {
	ret := make(map[gorl.Sprite]*CursesSprite)
	ret[gorl.SpriteBlank] = &CursesSprite{' ', termbox.ColorDefault, termbox.ColorDefault}
	ret[gorl.SpriteHumanMale] = &CursesSprite{'@', termbox.ColorBlue, termbox.ColorDefault}
	ret[gorl.SpriteHumanFemale] = &CursesSprite{'@', termbox.ColorMagenta, termbox.ColorDefault}
	ret[gorl.SpriteInfernalMale] = &CursesSprite{'&', termbox.ColorCyan, termbox.ColorDefault}
	ret[gorl.SpriteInfernalFemale] = &CursesSprite{'&', termbox.ColorRed, termbox.ColorDefault}
	ret[gorl.SpriteKoboldMale] = &CursesSprite{'k', termbox.ColorGreen, termbox.ColorDefault}
	ret[gorl.SpriteKoboldFemale] = &CursesSprite{'k', termbox.ColorYellow, termbox.ColorDefault}
	ret[gorl.SpriteMonsterUnknown] = &CursesSprite{'?', termbox.ColorRed, termbox.ColorDefault}
	ret[gorl.SpriteItemGold] = &CursesSprite{'$', termbox.ColorYellow, termbox.ColorDefault}
	ret[gorl.SpriteItemWeaponGeneric] = &CursesSprite{')', termbox.ColorYellow, termbox.ColorDefault}
	ret[gorl.SpriteItemAppGeneric] = &CursesSprite{'[', termbox.ColorCyan, termbox.ColorDefault}
	ret[gorl.SpriteItemPotion] = &CursesSprite{'!', termbox.ColorCyan, termbox.ColorDefault}
	return ret
}
