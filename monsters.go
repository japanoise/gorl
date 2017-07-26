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

func initMonsters() error {
	Bestiary = make(map[MonsterID]Monster)
	return loadConfigFile("monsters.json", &Bestiary)
}

func GetMonster(race MonsterID, female bool) *Critter {
	monst := Bestiary[race]
	return &Critter{0, 0, race, "", GenStatBlock(monst.HitDice, monst.Level),
		female, []*Item{}, 0, nil, nil, nil, 0, nil}
}

func RandomCritter(elevation int) *Critter {
	ret := GetMonster(MonsterUnknown, false)
	return ret
}
