package gorl

/* Interface for graphics frontend */

type Graphics interface {
	Start() error                                  /* Init the frontend */
	End()                                          /* Close the frontend */
	Dungeon(dun *Map, x, y int)                    /* Draw a dungeon level (may be a viewport) */
	Overworld(overworld *Map, x, y int)            /* Draw the overworld (ditto) */
	Message(msg string)                            /* Show a message; block and return when user ack's it */
	Menu(prompt string, choices []string) string   /* Show a selection menu*/
	MenuIndex(prompt string, choices []string) int /* Show a selection menu, return the index*/
	GetString(prompt string, empty bool) string    /* Get a free string */
}

/* Input interface */

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
	Quit
	Warp
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

type Input interface {
	GetAction() Control                   /* Get one command from the input */
	GetDirection(prompt string) Direction /* Get a direction to do some action */
}
