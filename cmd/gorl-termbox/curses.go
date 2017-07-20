package main

/*
   classic roguelike frontend, input method.
   despite the name, the implementation is left up to termbox.
*/

import (
	"github.com/japanoise/gorl"
	"github.com/japanoise/termbox-util"
	"github.com/nsf/termbox-go"
)

func NewCurses() *Curses {
	retval := &Curses{getMonsterSprites(), getTileSprites()}
	return retval
}

type CursesSprite struct {
	Ru rune
	Fg termbox.Attribute
	Bg termbox.Attribute
}

func Draw(x, y int, spr *CursesSprite) {
	termbox.SetCell(x, y, spr.Ru, spr.Fg, spr.Bg)
}

type Curses struct {
	MonsterSprites map[gorl.Sprite]*CursesSprite
	TileSprites    map[gorl.TileID]*CursesSprite
}

func (c *Curses) Start() error {
	err := termbox.Init()
	return err
}

func (c *Curses) End() {
	termbox.Close()
}

func (c *Curses) drawAt(dun *gorl.Map, screenx, screeny, x, y int) {
	here := dun.Tiles[x][y].Here
	if here == nil {
		Draw(screenx, screeny, c.TileSprites[dun.Tiles[x][y].Id])
	} else {
		Draw(screenx, screeny, c.MonsterSprites[here.GetSprite()])
	}
}

func (c *Curses) Dungeon(dun *gorl.Map, x, y int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()
	leftedge := x - (width / 2)
	topedge := y - (height / 2)
	for screenx := 0; screenx < width; screenx++ {
		for screeny := 0; screeny < width; screeny++ {
			x, y := screenx+leftedge, screeny+topedge
			if x >= 0 && x < dun.SizeX && y >= 0 && y < dun.SizeY {
				c.drawAt(dun, screenx, screeny, x, y)
			}
		}
	}
	termbox.Flush()
}

func clearLine(y, width int) {
	for i := 0; i < width; i++ {
		eraseCh(i, y)
	}
}

func putCh(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, termbox.ColorDefault, termbox.ColorDefault)
}

func eraseCh(x, y int) {
	termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
}

func (c *Curses) GetString(prompt string, empty bool) string {
	ret := termutil.Prompt(prompt, nil)
	for empty && ret == "" {
		ret = termutil.Prompt(prompt, nil)
	}
	termbox.HideCursor()
	return ret
}

func clearLines(n, width int) {
	for i := 0; i <= n; i++ {
		clearLine(i, width)
	}
}

func (c *Curses) MenuIndex(prompt string, choices []string) int {
	return termutil.ChoiceIndex(prompt, choices, 0)
}

func (c *Curses) Menu(prompt string, choices []string) string {
	sel := c.MenuIndex(prompt, choices)
	return choices[sel]
}

func (c *Curses) GetAction() gorl.Control {
	ev := termbox.PollEvent()
	if ev.Type == termbox.EventKey {
		if ev.Ch == 0 {
			switch ev.Key {
			case termbox.KeyCtrlC:
				return gorl.Quit
			case termbox.KeyArrowUp:
				return gorl.PlayerUp
			case termbox.KeyArrowDown:
				return gorl.PlayerDown
			case termbox.KeyArrowLeft:
				return gorl.PlayerLeft
			case termbox.KeyArrowRight:
				return gorl.PlayerRight
			default:
				return gorl.DoNothing
			}
		} else {
			switch ev.Ch {
			case 'q':
				return gorl.Quit
			case 'w':
				return gorl.Warp
			case '<':
				return gorl.PlayerClimbUp
			case '>':
				return gorl.PlayerClimbDown
			default:
				return gorl.DoNothing
			}
		}
	} else {
		return gorl.DoNothing
	}
}

func (c *Curses) Message(str string) {
	width, _ := termbox.Size()
	clearLine(0, width)
	drawString(0, 0, str)
	termbox.Flush()
	ev := termbox.PollEvent()
	for ev.Type != termbox.EventKey {
		ev = termbox.PollEvent()
	}
	clearLine(0, width)
	termbox.Flush()
}

func drawString(x, y int, str string) {
	drawStringDetails(x, y, str, termbox.ColorDefault, termbox.ColorDefault)
}

func drawStringDetails(x, y int, str string, fg, bg termbox.Attribute) {
	os := 0
	for _, runeValue := range str {
		termbox.SetCell(x+os, y, runeValue, fg, bg)
		os += termutil.Runewidth(runeValue)
	}
}
