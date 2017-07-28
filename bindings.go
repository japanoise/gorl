package gorl

type Control string

const (
	ControlInvalid  Control = ""
	DoNothing               = "Do nothing"
	PlayerUp                = "Move up"
	PlayerDown              = "Move down"
	PlayerLeft              = "Move left"
	PlayerRight             = "Move right"
	PlayerClimbDown         = "Climb down"
	PlayerClimbUp           = "Climb up"
	PlayerLook              = "Look"
	PlayerInventory         = "Inventory"
	PlayerStats             = "Stats"
	PlayerZapSpell          = "Zap Spell"
	Quit                    = "Quit"
	DoSaveGame              = "Save"
)

var bindings map[string]Control

func initBindings() error {
	bindings = make(map[string]Control)
	return loadConfigFile("bindings.json", &bindings)
}

func GetBinding(key string) Control {
	return bindings[key]
}
