package gorl

import "github.com/japanoise/engutil"

var playableRaces []MonsterID = []MonsterID{
	MonsterHuman,
	MonsterKobold,
	MonsterInfernal,
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
	return player
}
