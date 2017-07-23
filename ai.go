package gorl

import "math"

type AiData struct {
	Active bool
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

func AiOneTurn(state *State, player *Critter) bool {
	dead := false
	for _, monster := range state.Monsters {
		if monster == nil {
			continue
		}
		if monster.AI == nil {
			monster.AI = newAI()
		}
		if monster.AI.Active {
			target := monster.Chase(state.CurLevel, player.X, player.Y)
			if target == player {
				dead = dead || Attack(true, false, state.CurLevel, state.Out, monster, player)
			}
		}
	}
	return dead
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
	return &AiData{false}
}
