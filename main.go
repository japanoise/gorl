package gorl

import (
	"strconv"
)

func MainLoop(g Graphics, i Input) {
	mydun := DigDungeon(3)
	dunlevel, duncritters, _, _ := mydun.GetDunLevel(0, 1, []*Critter{})
	state := State{
		duncritters,
		1,
		i,
		g,
	}
	state.Monsters = make([]*Critter, len(duncritters), len(duncritters))
	debug.Print("Copy to from: ", state.Monsters, duncritters)
	copy(state.Monsters, duncritters)
	debug.Print("Copy to from: ", state.Monsters, duncritters)

	state.Out.Start()
	defer state.Out.End()

	player := GetMonster(MonsterHuman)
	player.Name = state.Out.GetString("Your name?", false)
	dunlevel.PlaceCritterAtUpStairs(player)

	playing := true
	state.Out.Dungeon(dunlevel, player.X, player.Y)
	for playing {
		state.Out.Dungeon(dunlevel, player.X, player.Y)
		act := state.In.GetAction()
		var target *Critter
		switch act {
		case PlayerClimbDown:
			if dunlevel.Tiles[player.X][player.Y].Id == TileStairDown {
				state.Out.Message("You climb down the stairs...")
				state.Dungeon++
				dunlevel, duncritters, _, _ = mydun.GetDunLevel(state.Dungeon-1, state.Dungeon, state.Monsters)
				state.Monsters = make([]*Critter, len(duncritters), len(duncritters))
				copy(state.Monsters, duncritters)
				dunlevel.PlaceCritterAtUpStairs(player)
			} else {
				state.Out.Message("There are no stairs here!")
			}
		case PlayerClimbUp:
			if dunlevel.Tiles[player.X][player.Y].Id == TileStairUp {
				if state.Dungeon == 1 {
					state.Out.Message("There's daylight at the top of the stairs!")
					playing = false
				} else {
					state.Out.Message("You climb up the stairs...")
					state.Dungeon--
					dunlevel, duncritters, _, _ = mydun.GetDunLevel(state.Dungeon+1, state.Dungeon, state.Monsters)
					state.Monsters = make([]*Critter, len(duncritters), len(duncritters))
					copy(state.Monsters, duncritters)
					dunlevel.PlaceCritterAtDownStairs(player)
				}
			} else {
				state.Out.Message("There are no stairs here!")
			}
		case Warp:
			x, err := strconv.Atoi(state.Out.GetString("x", false))
			y, err2 := strconv.Atoi(state.Out.GetString("y", false))
			if err != nil || err2 != nil {
				continue
			}
			target = Move(dunlevel, player, x, y)
		case PlayerDown:
			target = Move(dunlevel, player, 0, +1)
		case PlayerUp:
			target = Move(dunlevel, player, 0, -1)
		case PlayerLeft:
			target = Move(dunlevel, player, -1, 0)
		case PlayerRight:
			target = Move(dunlevel, player, +1, 0)
		case PlayerLook:
			Look(dunlevel, g, i, player)
		case Quit:
			playing = false
		}
		if target != nil && target.Collide != nil {
			delete := target.Collide(dunlevel, state.Out, target, player)
			if delete {
				target.Delete(dunlevel)
				for i, crit := range state.Monsters {
					if crit == target {
						state.Monsters[i] = nil
					}
				}
			}
		}
	}
}
