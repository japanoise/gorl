package gorl

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

func StartGame(g Graphics, i Input) error {
	state := &State{
		[]*Critter{},
		nil,
		0,
		i,
		g,
		25200000, // 7 am
		12,       // on the 12 of
		5,        // May
		1432,     // 1432
		&PlayerData{},
	}
	state.Out.Start()
	ierr := initAll()
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
			player := CharGen(state.Out)
			over := OverworldGen(player, 15, 15)
			doMainLoop(state, player, over, nil)
		} else if menuitem == 1 {
			player, newstate, over, err := LoadGame(state)
			if err != nil {
				state.Out.AckMessage(err.Error())
				continue
			} else {
				state.CurLevel = newstate.CurLevel
				state.Monsters = newstate.Monsters
				state.Dungeon = newstate.Dungeon
			}
			if state.Dungeon == 0 {
				player.X = over.SavedPx
				player.Y = over.SavedPy
				state.Dungeon = 0
				state.CurLevel = over.M
				doMainLoop(state, player, over, nil)
			} else if state.Dungeon > 0 {
				doMainLoop(state, player, over, over.M.Tiles[over.SavedPx][over.SavedPy].OwData.Dungeon)
			} else {
				over.MetaOw.Tiles[player.X][player.Y].Here = player
				state.CurLevel = over.MetaOw
				doMainLoop(state, player, over, nil)
			}
		} else {
			playing = false
		}
	}
	return nil
}

func initAll() error {
	initFuncs := []func() error{
		initDirs,
		startLogging,
		initRng,
		initMonsters,
		initItems,
		initTiles,
		initSpells,
		initBindings,
		initHunger,
	}
	for _, f := range initFuncs {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

func initDirs() error {
	if configdir != "" && datadir != "" {
		return nil
	}
	configdir = os.Getenv("XDG_CONFIG_HOME")
	datadir = os.Getenv("XDG_DATA_HOME")
	if configdir == "" {
		h, err := homedir.Dir()
		if err != nil {
			return err
		}
		configdir = h + string(os.PathSeparator) + ".config" + string(os.PathSeparator) + "gorl"
	} else {
		configdir += string(os.PathSeparator) + "gorl"
	}
	if datadir == "" {
		h, err := homedir.Dir()
		if err != nil {
			return err
		}
		datadir = h + string(os.PathSeparator) + ".local" + string(os.PathSeparator) + "share" + string(os.PathSeparator) + "gorl"
	} else {
		datadir += string(os.PathSeparator) + "gorl"
	}
	err := os.MkdirAll(datadir, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(configdir, 0755)
	return err
}

func startLogging() error {
	out, err := os.Create("debug.log")
	debug = log.New(out, "", log.Lshortfile)
	return err
}

func doMainLoop(state *State, player *Critter, over *Overworld, stdun *StateDungeon) {
	if state.Dungeon > 0 {
		state.CurLevel.Tiles[player.X][player.Y].Here = player
	} else if state.Dungeon == 0 {
		state.CurLevel = over.M
	} else if state.Dungeon == -1 {
		state.CurLevel = over.MetaOw
	} else if state.Dungeon == -2 {
		state.CurLevel, state.Monsters = VilGen(over.M.Tiles[over.SavedPx][over.SavedPy].OwData.OWSeed)
		state.CurLevel.Tiles[player.X][player.Y].Here = player
	}
	mydun := stdun
	playing := true
	pmoved := true
	pdjmap := BlankDMap(state.CurLevel)
	pdjmap.Calc(player)
	for playing {
		// Recalculate the status bar
		status := CalcStatus(state, player)
		// Draw the level
		if state.Dungeon > 0 {
			CalcVisibility(state.CurLevel, player, 20) //Eventually, this will be torch level.
			state.Out.Dungeon(state.CurLevel, player.X, player.Y, status)
		} else {
			CalcVisibility(state.CurLevel, player, 40)
			state.Out.Overworld(state.CurLevel, player.X, player.Y, status)
		}
		var target *Critter
		act := state.In.GetAction() // Poll for an action
		for act == ExtCmd {
			act = Control(state.Out.GetString("Command", true))
		}
		switch act { // Act on the action
		case PlayerClimbUp:
			if state.Dungeon == -1 {
				state.Out.Message("There are no stairs to climb up here!")
			} else if state.Dungeon == 0 {
				pmoved = overup(state, player, over)
			} else if state.Dungeon > 0 {
				pmoved = dungeonclimbup(state, player, over, mydun)
			} else {
				pmoved = returntoow(state, player, over)
			}
		case PlayerClimbDown:
			if state.Dungeon == 0 {
				mydun, pmoved = overdown(state, player, over)
			} else if state.Dungeon == -1 {
				pmoved = metadown(state, player, over)
			} else if state.Dungeon > 0 {
				pmoved = dungeonclimbdown(state, player, over, mydun)
			} else {
				state.Out.Message("There are no stairs to climb down here!")
			}
		case PlayerDown:
			target = Move(state.CurLevel, player, 0, +1)
			pmoved = true
		case PlayerUp:
			target = Move(state.CurLevel, player, 0, -1)
			pmoved = true
		case PlayerLeft:
			target = Move(state.CurLevel, player, -1, 0)
			pmoved = true
		case PlayerRight:
			target = Move(state.CurLevel, player, +1, 0)
			pmoved = true
		case PlayerNE:
			target = Move(state.CurLevel, player, 1, -1)
			pmoved = true
		case PlayerNW:
			target = Move(state.CurLevel, player, -1, -1)
			pmoved = true
		case PlayerSE:
			target = Move(state.CurLevel, player, 1, 1)
			pmoved = true
		case PlayerSW:
			target = Move(state.CurLevel, player, -1, 1)
			pmoved = true
		case PlayerZapSpell:
			c := ZapSpell(state, player, state.CurLevel)
			for _, crit := range c {
				if crit == player {
					state.Out.Message("The spell hits you!")
					if player.IsDead() {
						state.Out.Message("You died!")
						state.Out.DeathScreen(player, "their own magic")
						return
					}
				} else if crit.IsDead() {
					crit.Kill(state)
				}
			}
		case PlayerLook:
			Look(state.CurLevel, state.Out, state.In, player)
		case PlayerInventory:
			Inventory(state, player)
		case PlayerStats:
			player.CompleteDescription(state.Out, GetHungerString(state.Player.Hunger))
		case DoSaveGame:
			if state.Dungeon == 0 {
				over.SavedPx = player.X
				over.SavedPy = player.Y
			} else if state.Dungeon > 0 {
				items := make([]*DungeonItem, 0, 20)
				mydun.Items[state.Dungeon] = state.CurLevel.CollectItems(items)
			}
			err := SaveGame(player, state, over)
			if err == nil {
				state.Out.Message("Game saved.")
			} else {
				state.Out.Message(err.Error())
			}
		case Quit:
			playing = false
		case ViewMessages:
			state.Out.ShowMessageLog()
		case GetItems:
			// If the player is standing on any items, snarf them to the player's inventory
			if state.CurLevel.Tiles[player.X][player.Y].Items != nil && len(state.CurLevel.Tiles[player.X][player.Y].Items) > 0 {
				Grab(state, player)
			} else {
				state.Out.Message("You don't see anything here.")
			}
		case DropItems:
			DropItem(state, player)
		case DoNothing:
		default:
			state.Out.Message("Key unbound or unknown command.")
			continue
		}

		// End of actions, now act on the consequences
		if target != nil {
			if !target.HasFlags(FlagFriendly) || state.Out.YN("Are you sure you want to attack "+target.GetTheName()+"?") {
				target.Flags &= (0xFF ^ FlagFriendly) // Unset friendly flag
				delete := Attack(true, true, state.CurLevel, state.Out, player, target)
				if delete {
					target.Kill(state)
				}
			}
		}

		// If the player's moved, recalculate the Dijkstra map and tell her what's here/loot gold
		if pmoved {
			pdjmap = BlankDMap(state.CurLevel)
			pdjmap.Calc(player)
			pmoved = false
			if state.CurLevel.Tiles[player.X][player.Y].Items != nil {
				showmsg := false
				loot := false
				msg := ""
				for _, item := range state.CurLevel.Tiles[player.X][player.Y].Items {
					msg += item.DescribeExtra() + ","
					showmsg = true
					if item.Class == ItemClassCurrency {
						loot = true
					}
				}
				if loot {
					Grab(state, player)
				} else if showmsg {
					state.Out.Message(msg)
				}
			}
		}

		// Make the monsters act
		playerdead, killer := AiOneTurn(state, player, pdjmap)
		if playerdead {
			state.Out.Message("You died!")
			state.Out.DeathScreen(player, killer)
			playing = false
		}
	}
}

func CalcStatus(state *State, player *Critter) string {
	if state.Dungeon > 0 {
		return fmt.Sprintf("[%d/%d hp] [%d/%d mp] %s, level %d, on level %d",
			player.Stats.CurHp, player.Stats.MaxHp, player.Stats.CurMp, player.Stats.MaxMp,
			player.GetName(), player.Stats.Level, state.Dungeon)
	} else if state.Dungeon == -1 {
		return fmt.Sprintf("[%d/%d hp] [%d/%d mp] %s, level %d, on the fast travel map",
			player.Stats.CurHp, player.Stats.MaxHp, player.Stats.CurMp, player.Stats.MaxMp,
			player.GetName(), player.Stats.Level)
	} else if state.Dungeon == -2 {
		return fmt.Sprintf("[%d/%d hp] [%d/%d mp] %s, level %d, in a civilized village",
			player.Stats.CurHp, player.Stats.MaxHp, player.Stats.CurMp, player.Stats.MaxMp,
			player.GetName(), player.Stats.Level)
	} else {
		return fmt.Sprintf("[%d/%d hp] [%d/%d mp] %s, level %d, wandering the wilderness",
			player.Stats.CurHp, player.Stats.MaxHp, player.Stats.CurMp, player.Stats.MaxMp,
			player.GetName(), player.Stats.Level)
	}
}

func dungeonclimbup(state *State, player *Critter, over *Overworld, mydun *StateDungeon) bool {
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
		return true
	} else if state.CurLevel.Tiles[player.X][player.Y].Id == TileStairDown {
		state.Out.Message("These stairs only lead down!")
		return false
	} else {
		state.Out.Message("There are no stairs here!")
		return false
	}
}

func overdown(state *State, player *Critter, over *Overworld) (*StateDungeon, bool) {
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
			return mydun, true
		}
	} else if state.CurLevel.Tiles[player.X][player.Y].Id == TileOverworldVillage {
		tile := state.CurLevel.Tiles[player.X][player.Y]
		if tile.OwData == nil {
			state.Out.Message("This village doesn't welcome adventurers.")
		} else {
			state.Out.Message("You make your way into the little village...")
			over.SavedPx = player.X
			over.SavedPy = player.Y
			state.Dungeon = -2
			state.CurLevel, state.Monsters = VilGen(tile.OwData.OWSeed)
			player.X, player.Y = 50, 99
			state.CurLevel.Tiles[player.X][player.Y].Here = player
		}
	}
	return nil, false
}

func overup(state *State, player *Critter, over *Overworld) bool {
	over.SavedPx = player.X
	over.SavedPy = player.Y
	state.Dungeon = -1
	state.CurLevel = over.MetaOw
	player.X = over.MetaPx
	player.Y = over.MetaPy
	over.MetaOw.Tiles[player.X][player.Y].Here = player
	return true
}

func dungeonclimbdown(state *State, player *Critter, over *Overworld, mydun *StateDungeon) bool {
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
		return true
	} else if state.CurLevel.Tiles[player.X][player.Y].Id == TileStairUp {
		state.Out.Message("These stairs only lead up!")
	} else {
		state.Out.Message("There are no stairs here!")
	}
	return false
}

func metadown(state *State, player *Critter, over *Overworld) bool {
	if player.X == over.MetaPx && player.Y == over.MetaPy {
		player.X = over.SavedPx
		player.Y = over.SavedPy
		state.CurLevel = over.M
	} else {
		over.MetaPx = player.X
		over.MetaPy = player.Y
		state.CurLevel = over.MetaOw.MetaOWGenMap(player.X, player.Y)
		over.M = state.CurLevel
		player.X, player.Y = 50, 50
		state.CurLevel.Tiles[player.X][player.Y].Here = player
	}
	state.Dungeon = 0
	return true
}

func returntoow(state *State, player *Critter, over *Overworld) bool {
	state.Out.Message("You disembark.")
	player.X = over.SavedPx
	player.Y = over.SavedPy
	state.CurLevel = over.M
	state.Dungeon = 0
	return true
}

func loadConfigFile(configfile string, v interface{}) error {
	file, err := os.Open(configdir + string(os.PathSeparator) + configfile)
	if err != nil {
		// Attempt to load the file from the bindata
		bytes, err := Asset("bindata/" + configfile)
		if err != nil {
			return err
		}
		return json.Unmarshal(bytes, v)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(v)
}
