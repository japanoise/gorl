package gorl

import (
	"math/rand"

	"github.com/japanoise/engutil"
)

var playableRaces []MonsterID = []MonsterID{
	MonsterHuman,
	MonsterKobold,
	MonsterInfernal,
}

type HungerState uint8

const (
	HungerNormal HungerState = iota
	HungerHungry
	HungerStarving
	HungerDying
	HungerDead
)

type CharFlags uint8

const (
	FlagFriendly CharFlags = 1 << iota
	FlagShopkeep
	FlagInnkeep
)

type PlayerData struct {
	TimeSinceEaten uint32
	Hunger         HungerState
}

const (
	TimeStarvation uint32 = 6.048e+8 // one week
	TimeDying             = 5.184e+8 // 6 days
	TimeStarving          = 2.592e+8 // 3 days
	TimeHungry            = 1.44e+7  // 4 hours
)

var hungerStrings map[HungerState]string

func initHunger() error {
	hungerStrings = make(map[HungerState]string)
	hungerStrings[HungerNormal] = "Not hungry."
	hungerStrings[HungerHungry] = "Hungry."
	hungerStrings[HungerStarving] = "Starving."
	hungerStrings[HungerDying] = "Dying of starvation."
	hungerStrings[HungerDead] = "Dead from starvation."
	return nil
}

func GetHungerString(h HungerState) string {
	return hungerStrings[h]
}

func CharGen(g Graphics) *Critter {
	choices := make([]string, len(playableRaces))
	for i, r := range playableRaces {
		choices[i] = engutil.ASlashAn(Bestiary[r].Name)
	}
	choice := g.MenuIndex("Which race will you be?", choices)
	female := g.Menu("Male or female?", []string{"male", "female"}) == "female"
	player := GetMonster(playableRaces[choice], female)
	player.Stats.MaxHp += 10 // Players are special and get extra health
	player.Stats.CurHp += 10
	player.Name = g.GetString("Your name?", false)
	player.Casting = 0xFF
	player.SpellBook = []*Spell{
		&Spell{"Thunderbolt", SpellLightning, SmallDice(2, 4), SpellOther | SpellSorcery},
		&Spell{"Fireball", SpellFire, SmallDice(4, 6), SpellArea | SpellRitual},
		&Spell{"Cryosis", SpellIce, SmallDice(1, 3), SpellSelf | SpellRitual},
		&Spell{"Lay On Hands", SpellHeal, SmallDice(1, 3), SpellSelf | SpellHoly},
		&Spell{"Healing Ray", SpellHeal, SmallDice(1, 3), SpellOther | SpellHoly},
		&Spell{"Healing Burst", SpellHeal, SmallDice(4, 3), SpellArea | SpellHoly},
	}
	potion := NewItemOfClass("potion of healing", ItemClassPotion)
	potion.Magic = &Spell{}
	potion.Magic.Effect = SpellHeal
	potion.Magic.Potency = SmallDice(2, 6)
	potion.Magic.Data = SpellSelf
	potion.Magic.Name = "You feel a little better."
	slime := NewItemOfClass("slime-mold", ItemClassFood)
	player.Inv = []*InvItem{NewInvItem(potion, 3), NewInvItem(slime, 1)}
	return player
}

func NewFriendly(r *rand.Rand, flags CharFlags) *Critter {
	race := playableRaces[r.Intn(len(playableRaces))]
	female := r.Intn(2) == 0
	ret := GetMonster(race, female)
	ret.Flags = FlagFriendly | flags
	ret.Dialog = 1
	return ret
}

func (c *Critter) HasFlags(flags CharFlags) bool {
	return c.Flags&flags != 0
}

func (c *Critter) GenerateName() {
	if c.Female {
		c.Name = "Alice"
	} else {
		c.Name = "Bob"
	}
	if c.HasFlags(FlagInnkeep) {
		c.Name += " the barkeep"
	}
}
