package gorl

import (
	"math"

	dmap "github.com/japanoise/dmap"
)

type AiData struct {
	Active bool
	TimeEl uint32
}

// Interface for something placed into the world
type Placed interface {
	GetXY() (int, int) // What the placed thing thinks its position is
}

// General implementation of Placed
type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}

func (p *Point) GetXY() (int, int) {
	return p.X, p.Y
}

func AiOneTurn(state *State, player *Critter, pdjmap *dmap.DijkstraMap) (bool, string) {
	turnlen, starved := state.UpdateTimer(player)
	if starved {
		return true, "starvation"
	}
	dead := false
	killer := ""
	for _, monster := range state.Monsters {
		if monster == nil {
			continue
		}
		if monster.AI == nil {
			monster.AI = newAI()
		}
		if monster.AI.Active && !monster.HasFlags(FlagFriendly) {
			monster.AI.TimeEl += turnlen
			target := monster.Chase(state.CurLevel, pdjmap)
			for monster.AI.TimeEl >= monster.Speed {
				if target == player && !dead {
					dead = dead || Attack(true, false, state.CurLevel, state.Out, monster, player)
					if dead {
						killer = monster.GetName()
					}
				}
				monster.AI.TimeEl -= monster.Speed
			}
		}
	}
	return dead, killer
}

func CoordsToFloat(x, y int) (float64, float64) {
	return float64(x), float64(y)
}

// Distance between two placed things, using Pythagorean Theorem
func Dist(c1, c2 Placed) int {
	// Apparently FPUs make casting fast. I hope so…
	x1, y1 := CoordsToFloat(c1.GetXY())
	x2, y2 := CoordsToFloat(c2.GetXY())
	return int(math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2)))
}

func newAI() *AiData {
	return &AiData{false, 0}
}
