package gorl

import "math/rand"

var quips []string

func init() {
	quips = []string{
		"Alpha release coming Real Soon Now™",
		"10% fruitier than other leading brands!",
		"Does my bum look big in this?",
		"WHERE'S YOUR @ AT",
		"remind me to always run my architecture decisions by you :P",
		"taking suggestions for quips btw, they have to be 80 characters or less wide",
		"GUARDS! SIEZE HIM!!",
		"wizard bad! you hero! go kill!",
	}
}

// Get Minecraft-esque quips for the main menu
func GetQuip() string {
	return quips[rand.Intn(len(quips))]
}

/* Interface for graphics frontend */
type Graphics interface {
	Start() error                                  /* Init the frontend */
	End()                                          /* Close the frontend */
	Dungeon(dun *Map, x, y int)                    /* Draw a dungeon level (may be a viewport) */
	Overworld(overworld *Map, x, y int)            /* Draw the overworld (ditto) */
	Message(msg string)                            /* Show a message; block and return when user ack's it */
	LongMessage(msgs ...string)                    /* Multiline message */
	Menu(prompt string, choices []string) string   /* Show a selection menu*/
	MenuIndex(prompt string, choices []string) int /* Show a selection menu, return the index*/
	MainMenu(choices []string) int                 /* Show a selection menu, return the index. Implementor is requested to make it fancy. */
	GetString(prompt string, empty bool) string    /* Get a free string */
	DeathScreen(player *Critter)                   /* Called when the player dies. Extra points for a tombstone. */
}

/* Input interface */
type Input interface {
	GetAction() Control                   /* Get one command from the input */
	GetDirection(prompt string) Direction /* Get a direction to do some action */
}

type Control uint8

const (
	DoNothing Control = iota
	PlayerUp
	PlayerDown
	PlayerLeft
	PlayerRight
	PlayerClimbDown
	PlayerClimbUp
	PlayerLook
	PlayerInventory
	PlayerStats
	Quit
	DoSaveGame
)

type Direction uint8

const (
	DirNorth Direction = iota
	DirSouth
	DirWest
	DirEast
	DirNE
	DirSE
	DirNW
	DirSW
	DirUp
)
