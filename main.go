package gorl

import (
	"strconv"
)

func StartGame(g Graphics, i Input) error {
	state := &State{
		[]*Critter{},
		nil,
		0,
		i,
		g,
	}
	state.Out.Start()
	ierr := InitDirs()
	if ierr != nil {
		return ierr
	}
	defer state.Out.End()
	playing := true
	for playing {
		menuitem := state.Out.MainMenu([]string{"New Game", "Load Game", "Quit"})
		if menuitem == 0 {
			state.CurLevel = nil
			state.Monsters = nil
			state.Dungeon = 0
			player := GetMonster(MonsterHuman)
			player.Name = state.Out.GetString("Your name?", false)
			over := OverworldGen(player, 15, 15)
			doMainLoop(state, player, over, nil)
		} else if menuitem == 1 {
			player, newstate, over, err := LoadGame(state)
			if err != nil {
				state.Out.Message(err.Error())
				continue
			} else {
				state.CurLevel = newstate.CurLevel
				state.Monsters = newstate.Monsters
				state.Dungeon = newstate.Dungeon
			}
			if state.Dungeon <= 0 {
				player.X = over.SavedPx
				player.Y = over.SavedPy
				state.Dungeon = 0
				state.CurLevel = over.M
				doMainLoop(state, player, over, nil)
			} else {
				doMainLoop(state, player, over, over.M.Tiles[over.SavedPx][over.SavedPy].OwData.Dungeon)
			}
		} else {
			playing = false
		}
	}
	return nil
}

func doMainLoop(state *State, player *Critter, over *Overworld, stdun *StateDungeon) {
	mydun := stdun
	playing := true
	msg := ""
	showmsg := false
	if state.Dungeon > 0 {
		state.CurLevel.Tiles[player.X][player.Y].Here = player
	} else {
		state.CurLevel = over.M
	}
	for playing {
		// Draw the level
		if state.Dungeon > 0 {
			CalcVisibility(state.CurLevel, player, 20) //Eventually, this will be torch level.
			state.Out.Dungeon(state.CurLevel, player.X, player.Y)
		} else {
			CalcVisibility(state.CurLevel, player, 40)
			state.Out.Overworld(over.M, player.X, player.Y)
		}
		if showmsg { // show any delayed messages
			state.Out.Message(msg)
			showmsg = false
		}
		var target *Critter
		act := state.In.GetAction() // Poll for an action
		switch act {                // Act on the action
		case PlayerClimbUp:
			if state.Dungeon <= 0 {
				state.Out.Message("There are no stairs to climb up here!")
			} else {
				dungeonclimbup(state, player, over, mydun)
			}
		case PlayerClimbDown:
			if state.Dungeon <= 0 {
				mydun = overdown(state, player, over)
			} else {
				dungeonclimbdown(state, player, over, mydun)
			}
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
			Look(state.CurLevel, state.Out, state.In, player)
		case PlayerInventory:
			UseItem(state.Out, player, ShowItemList(state.Out, player.Gold, player.Inv))
		case PlayerStats:
			player.CompleteDescription(state.Out)
		case DoSaveGame:
			over.SavedPx = player.X
			over.SavedPy = player.Y
			err := SaveGame(player, state, over)
			if err == nil {
				state.Out.Message("Game saved.")
			} else {
				state.Out.Message(err.Error())
			}
		case Quit:
			playing = false
		}

		// End of actions, now act on the consequences
		if target != nil {
			delete := Attack(true, true, state.CurLevel, state.Out, player, target)
			if delete {
				state.Out.Message("You have defeated " + target.GetTheName())
				target.Delete(state.CurLevel)
				for i, crit := range state.Monsters {
					if crit == target {
						state.Monsters[i] = nil
					}
				}
			}
		}

		// If the player is standing on any items, snarf them to the player's inventory
		if state.CurLevel != nil && state.CurLevel.Tiles[player.X][player.Y].Items != nil {
			player.SnarfItems(state.CurLevel.Tiles[player.X][player.Y].Items)
			msg = ""
			showmsg = false
			for _, item := range state.CurLevel.Tiles[player.X][player.Y].Items {
				msg += item.DescribeExtra() + ","
				showmsg = true // Delay showing the message until after the screen is redrawn
			}
			state.CurLevel.Tiles[player.X][player.Y].Items = []*Item{}
		}

		// Make the monsters act
		playerdead := AiOneTurn(state, player)
		if playerdead {
			state.Out.Message("You died!")
			state.Out.DeathScreen(player)
			playing = false
		}
	}
}

func dungeonclimbup(state *State, player *Critter, over *Overworld, mydun *StateDungeon) {
	if state.CurLevel.Tiles[player.X][player.Y].Id == TileStairUp {
		if state.Dungeon == 1 {
			state.Out.Message("There's daylight at the top of the stairs!")
		} else {
			state.Out.Message("You climb up the stairs...")
		}
		state.Dungeon--
		items := make([]*DungeonItem, 0, 20)
		items = state.CurLevel.CollectItems(items)
		var duncritters []*Critter
		state.CurLevel, duncritters, _, _ = mydun.GetDunLevel(state.Dungeon+1, state.Dungeon, state.Monsters, items)
		state.Monsters = make([]*Critter, len(duncritters), len(duncritters))
		copy(state.Monsters, duncritters)
		if state.Dungeon == 0 {
			player.X = over.SavedPx
			player.Y = over.SavedPy
			state.CurLevel = over.M
		} else {
			state.CurLevel.PlaceCritterAtDownStairs(player)
		}
	} else if state.CurLevel.Tiles[player.X][player.Y].Id == TileStairDown {
		state.Out.Message("These stairs only lead down!")
	} else {
		state.Out.Message("There are no stairs here!")
	}
}

func overdown(state *State, player *Critter, over *Overworld) *StateDungeon {
	if state.CurLevel.Tiles[player.X][player.Y].Id == TileOverworldDungeon {
		tile := state.CurLevel.Tiles[player.X][player.Y]
		if tile.OwData == nil || tile.OwData.Dungeon == nil {
			state.Out.Message("The entrance has caved in.")
		} else {
			state.Out.Message("You make your way down into the murky depths...")
			over.SavedPx = player.X
			over.SavedPy = player.Y
			state.Dungeon = 1
			mydun := tile.OwData.Dungeon
			dunlevel, duncritters, _, _ := mydun.GetDunLevel(0, 1, []*Critter{}, nil)
			state.Monsters = make([]*Critter, len(duncritters), len(duncritters))
			debug.Print("Copy to from: ", state.Monsters, duncritters)
			copy(state.Monsters, duncritters)
			debug.Print("Copy to from: ", state.Monsters, duncritters)
			dunlevel.PlaceCritterAtUpStairs(player)
			state.CurLevel = dunlevel
			return mydun
		}
	}
	return nil
}

func dungeonclimbdown(state *State, player *Critter, over *Overworld, mydun *StateDungeon) {
	if state.CurLevel.Tiles[player.X][player.Y].Id == TileStairDown {
		state.Out.Message("You climb down the stairs...")
		state.Dungeon++
		items := make([]*DungeonItem, 0, 20)
		if state.Dungeon > 1 {
			items = state.CurLevel.CollectItems(items)
		}
		var duncritters []*Critter
		state.CurLevel, duncritters, _, _ = mydun.GetDunLevel(state.Dungeon-1, state.Dungeon, state.Monsters, items)
		state.Monsters = make([]*Critter, len(duncritters), len(duncritters))
		copy(state.Monsters, duncritters)
		state.CurLevel.PlaceCritterAtUpStairs(player)
	} else if state.CurLevel.Tiles[player.X][player.Y].Id == TileStairUp {
		state.Out.Message("These stairs only lead up!")
	} else {
		state.Out.Message("There are no stairs here!")
	}
}
