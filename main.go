package gorl

import (
	"strconv"
)

func MainLoop(g Graphics, i Input) {
	state := State{
		[]*Critter{},
		nil,
		1,
		i,
		g,
	}
	state.Out.Start()
	defer state.Out.End()
	player := GetMonster(MonsterHuman)
	player.Name = state.Out.GetString("Your name?", false)
	over := OverworldGen(player, 15, 15)
	MainLoopOverworld(g, i, state, player, over)
}

func MainLoopOverworld(g Graphics, i Input, state State, player *Critter, over *Overworld) {
	state.CurLevel = over.M
	playing := true
	for playing {
		state.Out.Overworld(over.M, player.X, player.Y)
		act := state.In.GetAction()
		var target *Critter
		switch act {
		case PlayerClimbDown:
			if state.CurLevel.Tiles[player.X][player.Y].Id == TileOverworldDungeon {
				tile := state.CurLevel.Tiles[player.X][player.Y]
				if tile.OwData == nil || tile.OwData.Dungeon == nil {
					state.Out.Message("The entrance has caved in.")
				} else {
					over.SavedPx = player.X
					over.SavedPy = player.Y
					playing = !MainLoopDungeon(g, i, state, player, tile.OwData.Dungeon)
					player.X = over.SavedPx
					player.Y = over.SavedPy
				}
			}
		case PlayerClimbUp:
			state.Out.Message("There are no stairs to climb up here!")
		case Warp:
			x, err := strconv.Atoi(state.Out.GetString("x", false))
			y, err2 := strconv.Atoi(state.Out.GetString("y", false))
			if err != nil || err2 != nil {
				continue
			}
			target = Move(state.CurLevel, player, x, y)
		case PlayerDown:
			target = Move(state.CurLevel, player, 0, +1)
		case PlayerUp:
			target = Move(state.CurLevel, player, 0, -1)
		case PlayerLeft:
			target = Move(state.CurLevel, player, -1, 0)
		case PlayerRight:
			target = Move(state.CurLevel, player, +1, 0)
		case PlayerLook:
			Look(state.CurLevel, g, i, player)
		case Quit:
			playing = false
		}
		if target != nil && target.Collide != nil {
			delete := target.Collide(state.CurLevel, state.Out, target, player)
			if delete {
				target.Delete(state.CurLevel)
				for i, crit := range state.Monsters {
					if crit == target {
						state.Monsters[i] = nil
					}
				}
			}
		}
	}
}

func MainLoopDungeon(g Graphics, i Input, state State, player *Critter, mydun *StateDungeon) bool {
	dunlevel, duncritters, _, _ := mydun.GetDunLevel(0, 1, []*Critter{}, nil)
	state.Monsters = make([]*Critter, len(duncritters), len(duncritters))
	debug.Print("Copy to from: ", state.Monsters, duncritters)
	copy(state.Monsters, duncritters)
	debug.Print("Copy to from: ", state.Monsters, duncritters)
	dunlevel.PlaceCritterAtUpStairs(player)

	playing := true
	for playing {
		state.Out.Dungeon(dunlevel, player.X, player.Y)
		act := state.In.GetAction()
		var target *Critter
		switch act {
		case PlayerClimbDown:
			if dunlevel.Tiles[player.X][player.Y].Id == TileStairDown {
				state.Out.Message("You climb down the stairs...")
				state.Dungeon++
				items := make([]*DungeonItem, 0, 20)
				if state.Dungeon > 1 {
					items = dunlevel.CollectItems(items)
				}
				dunlevel, duncritters, _, _ = mydun.GetDunLevel(state.Dungeon-1, state.Dungeon, state.Monsters, items)
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
				}
				state.Dungeon--
				items := make([]*DungeonItem, 0, 20)
				items = dunlevel.CollectItems(items)
				dunlevel, duncritters, _, _ = mydun.GetDunLevel(state.Dungeon+1, state.Dungeon, state.Monsters, items)
				state.Monsters = make([]*Critter, len(duncritters), len(duncritters))
				copy(state.Monsters, duncritters)
				if state.Dungeon != 0 {
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
			return true
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
	return false
}
