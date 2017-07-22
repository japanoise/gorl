package gorl

type MonsterID uint32

type Monster struct {
	Name       string
	SprM       Sprite
	SprF       Sprite
	BaseDamage uint8
	BaseAC     uint8
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
		"unknown creature", SpriteMonsterUnknown, SpriteMonsterUnknown, SmallDice(1, 6), 10,
	}
	Bestiary[MonsterHuman] = Monster{
		"human", SpriteHumanMale, SpriteHumanFemale, SmallDice(1, 4), 10,
	}
}

func GetMonster(race MonsterID) *Critter {
	return &Critter{0, 0, race, "", DefStatBlock(), false, []*Item{
		NewWeapon("mace", 10, 0, 0, Uncursed, SmallDice(1, 8)),
	}, 10, nil, nil, nil}
}
