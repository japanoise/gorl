package gorl

/* Interface for graphics frontend */

type Graphics interface {
	Start() error                                  /* Init the frontend */
	End()                                          /* Close the frontend */
	Dungeon(dun *Map, x, y int)                    /* Draw a dungeon level (may be a viewport) */
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
	Quit
	Warp
)

type Input interface {
	GetAction() Control /* Get one command from the input */
}
