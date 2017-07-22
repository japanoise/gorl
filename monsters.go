package gorl

type MonsterID uint32

type Monster struct {
	Name string
	SprM Sprite
	SprF Sprite
}

var Bestiary map[MonsterID]Monster

const (
	MonsterUnknown MonsterID = iota
	MonsterHuman
)

// Monster definitions
func init() {
	Bestiary = make(map[MonsterID]Monster)
	Bestiary[MonsterUnknown] = Monster{
		"unknown creature", SpriteMonsterUnknown, SpriteMonsterUnknown,
	}
	Bestiary[MonsterHuman] = Monster{
		"human", SpriteHumanMale, SpriteHumanFemale,
	}
}

func GetMonster(race MonsterID) *Critter {
	return &Critter{0, 0, race, "", DefStatBlock(), false, []*Item{}}
}
