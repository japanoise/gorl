package gorl

import "math"

type AiData struct {
	Active bool
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
		} else {
			dist := Dist(player, monster)
			if -10 < dist && dist < 10 {
				monster.AI.Active = true
			}
		}
	}
	return dead
}

// Distance between two critters, using Pythagorean Theorem
func Dist(c1, c2 *Critter) int {
	// Apparently FPUs make casting fast. I hope so…
	x1 := float64(c1.X)
	y1 := float64(c1.Y)
	x2 := float64(c2.X)
	y2 := float64(c2.Y)
	return int(math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2)))
}

func newAI() *AiData {
	return &AiData{false}
}
