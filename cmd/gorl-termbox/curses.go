package main

/*
   classic roguelike frontend, input method.
   despite the name, the implementation is left up to termbox.
*/

import (
	"fmt"
	"strings"
	"time"

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
	MessageLog  []string
}

type CursesSprite struct {
	Ru rune
	Fg termbox.Attribute
	Bg termbox.Attribute
}

func NewCurses() *Curses {
	retval := &Curses{getSprites(), getTileSprites(), nil, []string{"Beginning of log."}}
	return retval
}

func Draw(x, y int, spr *CursesSprite) {
	termbox.SetCell(x, y, spr.Ru, spr.Fg, spr.Bg)
}

func (c *Curses) Start() error {
	err := termbox.Init()
	termbox.SetInputMode(termbox.InputAlt)
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
	c.MessageLog = append(c.MessageLog, c.Messages...)
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
	for {
		ev := termbox.PollEvent()
		for ev.Type != termbox.EventKey {
			ev = termbox.PollEvent()
		}
		binding := gorl.GetBinding(termutil.ParseTermboxEvent(ev))
		switch binding {
		case gorl.PlayerUp:
			return gorl.DirNorth
		case gorl.PlayerDown:
			return gorl.DirSouth
		case gorl.PlayerLeft:
			return gorl.DirWest
		case gorl.PlayerRight:
			return gorl.DirEast
		case gorl.PlayerNE:
			return gorl.DirNE
		case gorl.PlayerSE:
			return gorl.DirSE
		case gorl.PlayerNW:
			return gorl.DirNW
		case gorl.PlayerSW:
			return gorl.DirSW
		case gorl.PlayerClimbUp:
			return gorl.DirUp
		case gorl.PlayerClimbDown:
			return gorl.DirDown
		case gorl.DoNothing:
			return gorl.DirSelf
		}
		sx, sy := termbox.Size()
		clearLine(sy-2, sx)
		drawString(0, sy-2, "Invalid direction "+termutil.ParseTermboxEvent(ev))
		termbox.Flush()
	}
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
				}
			}
			binding := gorl.GetBinding(termutil.ParseTermboxEvent(ev))
			switch binding {
			case gorl.PlayerDown:
				if sel != len(choices)-1 {
					sel++
				}
			case gorl.PlayerUp:
				if sel != 0 {
					sel--
				}
			}
		}
	}
}

func (c *Curses) DeathScreen(player *gorl.Critter, killer string) {
	c.flushMessages()
	termbox.PollEvent()
	looping := true
	for looping {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		sx, sy := termbox.Size()
		nl := termutil.RunewidthStr(player.Name)
		kl := termutil.RunewidthStr(killer)
		gs := fmt.Sprintf("%d AU", player.Gold)
		gl := termutil.RunewidthStr(gs)
		tsh := 14
		tsw := 10
		if nl > tsw {
			tsw = nl
			if tsw%2 != 0 {
				tsw++
			}
		}
		if kl > tsw {
			tsw = kl
			if tsw%2 != 0 {
				tsw++
			}
		}
		if gl > tsw {
			tsw = gl
			if tsw%2 != 0 {
				tsw++
			}
		}
		anc := (sx - (tsw + 8)) / 2
		yanc := (sy - tsh) / 2
		for i := 0; i < tsw-2; i++ {
			termutil.PrintRune(anc+5+i, yanc+0, '_', termbox.ColorDefault)
		}
		termutil.Printstring("/", anc+4, yanc+1)
		termutil.Printstring("\\", anc+3+tsw, yanc+1)
		termutil.Printstring("/", anc+3, yanc+2)
		termutil.Printstring("REST", anc+4+((tsw-4)/2), yanc+2)
		termutil.Printstring("\\", anc+4+tsw, yanc+2)
		termutil.Printstring("/", anc+2, yanc+3)
		termutil.Printstring("IN", anc+4+((tsw-2)/2), yanc+3)
		termutil.Printstring("\\", anc+5+tsw, yanc+3)
		termutil.Printstring("/", anc+1, yanc+4)
		termutil.Printstring("PEACE!", anc+4+((tsw-6)/2), yanc+4)
		termutil.Printstring("\\", anc+6+tsw, yanc+4)
		for i := 5; i < tsh; i++ {
			termutil.Printstring("|", anc+7+tsw, yanc+i)
			termutil.Printstring("|", anc, yanc+i)
		}
		termutil.Printstring(player.Name, anc+4+((tsw-nl)/2), yanc+6)
		termutil.Printstring(gs, anc+4+((tsw-gl)/2), yanc+7)
		termutil.Printstring("killed by", anc+4+((tsw-10)/2), yanc+8)
		termutil.Printstring(killer, anc+4+((tsw-kl)/2), yanc+9)
		termutil.Printstring(time.Now().Format("06-01-02"), anc+4+((tsw-8)/2), yanc+11)
		for i := 0; i < anc; i++ {
			termutil.Printstring("_", i, yanc+13)
		}
		for i := anc + 1; i < anc+tsw+7; i++ {
			termutil.Printstring("_", i, yanc+13)
		}
		for i := anc + 8 + tsw; i <= sx; i++ {
			termutil.Printstring("_", i, yanc+13)
		}
		termbox.Flush()
		looping = termbox.PollEvent().Type != termbox.EventKey
	}
}

func (c *Curses) ShowMessageLog() {
	termutil.DisplayScreenMessage(c.MessageLog...)
}
