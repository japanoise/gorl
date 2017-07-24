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
	player.Name = g.GetString("Your name?", false)
	return player
}
