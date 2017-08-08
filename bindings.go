package gorl

type Control string

const (
	ControlInvalid  Control = ""
	DoNothing               = "Do nothing"
	PlayerUp                = "Move up"
	PlayerDown              = "Move down"
	PlayerLeft              = "Move left"
	PlayerRight             = "Move right"
	PlayerNE                = "North East"
	PlayerSE                = "South East"
	PlayerNW                = "North West"
	PlayerSW                = "South West"
	PlayerClimbDown         = "Climb down"
	PlayerClimbUp           = "Climb up"
	PlayerLook              = "Look"
	PlayerInventory         = "Inventory"
	PlayerStats             = "Stats"
	PlayerZapSpell          = "Zap Spell"
	Quit                    = "Quit"
	DoSaveGame              = "Save"
	ExtCmd                  = "Extended command"
	ViewMessages            = "Messages"
	GetItems                = "Grab"
	DropItems               = "Drop"
	Chat                    = "Chat"
	Fire                    = "Fire"
	Open                    = "Open"
)

var bindings map[string]Control

func initBindings() error {
	bindings = make(map[string]Control)
	return loadConfigFile("bindings.json", &bindings)
}

func GetBinding(key string) Control {
	return bindings[key]
}
