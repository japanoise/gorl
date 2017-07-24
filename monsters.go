package gorl

type MonsterID uint32

type Monster struct {
	Name       string
	SprM       Sprite
	SprF       Sprite
	BaseDamage uint8
	BaseAC     uint8
	HitDice    uint8
	Level      uint8
}

var Bestiary map[MonsterID]Monster

const (
	MonsterUnknown MonsterID = iota
	MonsterHuman
	MonsterKobold
	MonsterInfernal
)

// Monster definitions
func init() {
	Bestiary = make(map[MonsterID]Monster)
	Bestiary[MonsterUnknown] = Monster{
		"unknown creature", SpriteMonsterUnknown, SpriteMonsterUnknown,
		SmallDice(1, 6), 10, SmallDice(1, 4), 1,
	}
	Bestiary[MonsterHuman] = Monster{
		"human", SpriteHumanMale, SpriteHumanFemale,
		SmallDice(1, 4), 10, SmallDice(1, 4), 1,
	}
	Bestiary[MonsterKobold] = Monster{
		"kobold", SpriteKoboldMale, SpriteKoboldFemale,
		SmallDice(1, 4), 10, SmallDice(1, 4), 1,
	}
	Bestiary[MonsterInfernal] = Monster{
		"infernal", SpriteInfernalMale, SpriteInfernalFemale,
		SmallDice(1, 4), 10, SmallDice(1, 4), 1,
	}
}

func GetMonster(race MonsterID, female bool) *Critter {
	monst := Bestiary[race]
	return &Critter{0, 0, race, "", GenStatBlock(monst.HitDice, monst.Level),
		female, []*Item{}, 0, nil, nil, nil}
}

func RandomCritter(elevation int) *Critter {
	ret := GetMonster(MonsterUnknown, false)
	return ret
}
