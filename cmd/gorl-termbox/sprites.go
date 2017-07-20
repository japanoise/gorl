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
	ret[gorl.TileGrass] = &CursesSprite{',', termbox.ColorGreen | termbox.AttrReverse, termbox.ColorDefault}
	ret[gorl.TileStairUp] = &CursesSprite{'<', termbox.ColorDefault, termbox.ColorDefault}
	ret[gorl.TileStairDown] = &CursesSprite{'>', termbox.ColorDefault, termbox.ColorDefault}
	return ret
}

func getMonsterSprites() map[gorl.Sprite]*CursesSprite {
	ret := make(map[gorl.Sprite]*CursesSprite)
	ret[gorl.SpriteBlank] = &CursesSprite{
		' ', termbox.ColorDefault, termbox.ColorDefault,
	}
	ret[gorl.SpriteHumanMale] = &CursesSprite{
		'@', termbox.ColorBlue, termbox.ColorDefault,
	}
	ret[gorl.SpriteHumanFemale] = &CursesSprite{
		'@', termbox.ColorMagenta, termbox.ColorDefault,
	}
	ret[gorl.SpriteMonsterUnknown] = &CursesSprite{
		'?', termbox.ColorRed, termbox.ColorDefault,
	}
	return ret
}
