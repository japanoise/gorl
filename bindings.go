package gorl

type Control uint8

const (
	ControlInvalid Control = iota
	DoNothing
	PlayerUp
	PlayerDown
	PlayerLeft
	PlayerRight
	PlayerClimbDown
	PlayerClimbUp
	PlayerLook
	PlayerInventory
	PlayerStats
	PlayerZapSpell
	Quit
	DoSaveGame
)

var bindings map[string]Control

func initBindings() error {
	bindings = make(map[string]Control)
	return loadConfigFile("bindings.json", &bindings)
}

func GetBinding(key string) Control {
	return bindings[key]
}
