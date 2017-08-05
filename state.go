package gorl

/* Game state */

type State struct {
	Monsters  []*Critter `json:"-"`
	CurLevel  *Map       `json:"-"`
	Dungeon   int
	In        Input    `json:"-"`
	Out       Graphics `json:"-"`
	TimeMili  uint32
	TimeDay   uint8
	TimeMonth uint8
	TimeYear  uint16
	Player    *PlayerData
}

const (
	MiliRollover  uint32 = 86400000
	DayRollover   uint8  = 30
	MonthRollover uint8  = 12
)

func (s *State) IncMili(amount uint32) {
	s.TimeMili += amount
	if s.TimeMili >= MiliRollover {
		s.TimeMili -= MiliRollover
		s.IncDay(1)
	}
}

func (s *State) IncDay(amount uint8) {
	s.TimeDay += amount
	if s.TimeDay >= DayRollover {
		s.TimeDay -= DayRollover
		s.IncMonth(1)
	}
}

func (s *State) IncMonth(amount uint8) {
	s.TimeMonth += amount
	if s.TimeMonth >= MonthRollover {
		s.TimeMonth -= MonthRollover
		s.TimeYear++
	}
}

func (s *State) UpdateTimer(player *Critter) (uint32, bool) {
	var oneturn uint32 = player.Speed
	if s.Dungeon == -1 {
		ret := oneturn * 100 * 100
		s.IncMili(ret)
		return ret, s.updateHunger(ret)
	} else if s.Dungeon == 0 {
		ret := oneturn * 100
		s.IncMili(ret)
		return ret, s.updateHunger(ret)
	} else {
		s.IncMili(oneturn)
		return oneturn, s.updateHunger(oneturn)
	}
}

func (s *State) updateHunger(time uint32) bool {
	s.Player.TimeSinceEaten += time
	if s.Player.TimeSinceEaten >= TimeStarvation {
		s.Out.Message("Unable to continue, you collapse.")
		return true
	} else if s.Player.TimeSinceEaten >= TimeDying {
		if s.Player.Hunger != HungerDying {
			s.Out.Message("You are dying of starvation!")
		}
		s.Player.Hunger = HungerDying
	} else if s.Player.TimeSinceEaten >= TimeStarving {
		if s.Player.Hunger != HungerStarving {
			s.Out.Message("You are starving!")
		}
		s.Player.Hunger = HungerStarving
	} else if s.Player.TimeSinceEaten >= TimeHungry {
		if s.Player.Hunger != HungerHungry {
			s.Out.Message("You are feeling hungry...")
		}
		s.Player.Hunger = HungerHungry
	}
	return false
}
