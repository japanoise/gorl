package main

/*
   classic roguelike frontend, input method.
   despite the name, the implementation is left up to termbox.
*/

import (
	"strings"

	"github.com/japanoise/gorl"
	"github.com/japanoise/termbox-util"
	"github.com/nsf/termbox-go"
)

// It's beautiful.
const logo string = `   .@@@@@@@@@@@@@b                         @@@
  @@@*           /                         @@@
 @@@                                       @@@
%@@,                                       @@@
@@@        ,,,,,,,                         @@@
@@@        @@@@@@@.    @@@@@@@.   @@@@*/,  @@@
*@@#           *@@.   @@     %@   @@       @@@
 @@@.          *@@.   @@      @(  @@       @@@
  #@@@,        (@@.   @@     &@   @@       @@@
    .@@@@@@@@@@@&      &@@@@@@    @@       @@@@@@@@@@@@@`

type Curses struct {
	Sprites     map[gorl.Sprite]*CursesSprite
	TileSprites map[gorl.TileID]*CursesSprite
	Messages    []string
}

type CursesSprite struct {
	Ru rune
	Fg termbox.Attribute
	Bg termbox.Attribute
}

func NewCurses() *Curses {
	retval := &Curses{getSprites(), getTileSprites(), nil}
	return retval
}

func Draw(x, y int, spr *CursesSprite) {
	termbox.SetCell(x, y, spr.Ru, spr.Fg, spr.Bg)
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
	if dun.Tiles[x][y].Lit {
		if here == nil {
			if len(dun.Tiles[x][y].Items) == 0 {
				Draw(screenx, screeny, c.TileSprites[dun.Tiles[x][y].Id])
			} else {
				Draw(screenx, screeny, c.Sprites[dun.Tiles[x][y].Items[0].Spr])
			}
		} else {
			Draw(screenx, screeny, c.Sprites[here.GetSprite()])
		}
	} else if dun.Tiles[x][y].Disc {
		sp := c.TileSprites[dun.Tiles[x][y].Id]
		Draw(screenx, screeny, &CursesSprite{sp.Ru, termbox.ColorBlack | termbox.AttrBold, termbox.ColorDefault})
	}
}

func (c *Curses) drawMapViewport(m *gorl.Map, x, y int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()
	leftedge := x - (width / 2)
	topedge := y - (height / 2)
	for screenx := 0; screenx < width; screenx++ {
		for screeny := 0; screeny < width; screeny++ {
			x, y := screenx+leftedge, screeny+topedge
			if x >= 0 && x < m.SizeX && y >= 0 && y < m.SizeY {
				c.drawAt(m, screenx, screeny, x, y)
			}
		}
	}
	termbox.Flush()
}

func (c *Curses) Dungeon(dun *gorl.Map, x, y int, status string) {
	c.drawMapViewport(dun, x, y)
	c.flushMessages()
	PrintStatus(status)
}

func PrintStatus(status string) {
	sx, sy := termbox.Size()
	for i := 0; i <= sx; i++ {
		termbox.SetCell(i, sy-1, ' ', termbox.AttrReverse, termbox.ColorDefault)
	}
	termutil.PrintstringColored(termbox.AttrReverse, status, 0, sy-1)
	termbox.Flush()
}

func (c *Curses) Overworld(overworld *gorl.Map, x, y int, status string) {
	c.drawMapViewport(overworld, x, y)
	c.flushMessages()
	PrintStatus(status)
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
	for !empty && ret == "" {
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
	for ev.Type != termbox.EventKey {
		ev = termbox.PollEvent()
	}
	return gorl.GetBinding(termutil.ParseTermboxEvent(ev))
}

func (c *Curses) Message(str string) {
	if c.Messages == nil {
		c.Messages = []string{str}
	} else {
		c.Messages = append(c.Messages, str)
	}
}

func (c *Curses) flushMessages() {
	if c.Messages == nil || len(c.Messages) == 0 {
		return
	}
	ml := len(c.Messages)
	_, width := termbox.Size()
	for i, msg := range c.Messages {
		c.showMessage(msg)
		if i != ml-1 {
			termutil.Printstring("<more>", 0, 1)
			termbox.Flush()
			ev := termbox.PollEvent()
			for ev.Type != termbox.EventKey {
				ev = termbox.PollEvent()
			}
			clearLine(0, width)
			termbox.Flush()
		} else {
			clearLine(1, 6)
			termbox.Flush()
		}
	}
	c.Messages = nil
}

func (c *Curses) showMessage(str string) {
	width, _ := termbox.Size()
	clearLine(0, width)
	drawString(0, 0, str)
	termbox.Flush()
}

func (c *Curses) LongMessage(msgs ...string) {
	termutil.DisplayScreenMessage(msgs...)
}

func (c *Curses) GetDirection(prompt string) gorl.Direction {
	drawString(0, 0, prompt)
	termbox.Flush()
	ev := termbox.PollEvent()
	for ev.Type != termbox.EventKey {
		ev = termbox.PollEvent()
	}
	if ev.Ch == 0 {
		switch ev.Key {
		case termbox.KeyArrowUp:
			return gorl.DirNorth
		case termbox.KeyArrowDown:
			return gorl.DirSouth
		case termbox.KeyArrowLeft:
			return gorl.DirWest
		case termbox.KeyArrowRight:
			return gorl.DirEast
		default:
			return gorl.DirUp
		}
	}
	return gorl.DirUp
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

func (c *Curses) MainMenu(choices []string) int {
	sel := 0
	quip := gorl.GetQuip()
	qw := termutil.RunewidthStr(quip)
	logolines := strings.Split(logo, "\n")
	logoheight := len(logolines)
	logowidth := len(logolines[logoheight-1])
	widths := make([]int, len(choices))
	for i, choice := range choices {
		widths[i] = termutil.RunewidthStr(choice)
	}
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		sx, _ := termbox.Size()
		offsetx := (sx - logowidth) / 2
		for li, line := range logolines {
			for i, ru := range line { //it's definitely just ascii here, so it should be ok to treat it as width 1
				termbox.SetCell(offsetx+i, li+1, ru, termbox.ColorDefault, termbox.ColorDefault)
			}
		}
		for i, choice := range choices {
			if i == sel {
				offsetx = (sx - (widths[i] + 4)) / 2
				termutil.Printstring("> "+choice+" <", offsetx, logoheight+4+i)
			} else {
				offsetx = (sx - widths[i]) / 2
				termutil.Printstring(choice, offsetx, logoheight+4+i)
			}
		}
		offsetx = (sx - qw) / 2
		termutil.PrintstringColored(termbox.ColorRed, quip, offsetx, logoheight+2)
		termbox.Flush()
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			if ev.Ch == 0 {
				switch ev.Key {
				case termbox.KeyEnter:
					return sel
				case termbox.KeyArrowDown:
					if sel != len(choices)-1 {
						sel++
					}
				case termbox.KeyArrowUp:
					if sel != 0 {
						sel--
					}
				}
			}
		}
	}
}

func (c *Curses) DeathScreen(player *gorl.Critter) {
	c.flushMessages()
	termbox.PollEvent()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termutil.Printstring("Rest in peace "+player.Name, 0, 0)
	termbox.Flush()
	termbox.PollEvent()
}
