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
	OneMinute     uint32 = 1000 * 60
	OneHour       uint32 = OneMinute * 60
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

func (s *State) GetLightLevel() int {
	// According to https://www.timeanddate.com/sun/uk/london, in the UK we have a
	// variance between 16:00 and 21:00 for sunset, and 08:00 to 04:00 for sunrise.
	// All calculations assume these values.
	if s.Dungeon > 0 {
		return 20 // Again, need to give the player torches
	} else {
		// Much easier to split the calculation here.
		if s.TimeMili < 12*OneHour {
			rise := s.GetSunrise()
			if s.TimeMili > rise {
				return 40
			} else {
				// fag-packet calculation: the sun takes 4 minutes to rise/set
				// so light level follows a straightforward line
				risestart := rise - 4*OneMinute
				if s.TimeMili < risestart {
					return 3
				}
				x := s.TimeMili - risestart
				return int((37*x)/(4*OneMinute) + 3)
			}
		} else {
			set := s.GetSunset()
			if s.TimeMili > set {
				return 3
			} else {
				setstart := set - 4*OneMinute
				if s.TimeMili < setstart {
					return 40
				}
				x := s.TimeMili - setstart
				debug.Println(x)
				return int(40 - ((37 * x) / 4 * OneMinute))
			}
		}
	}
}

func (s *State) getDayLengthFactor() uint32 {
	omonth := s.TimeMonth
	if omonth > MonthRollover/2 {
		// So that the summer solstice happens in the sixth month, this value
		// stays between 0 and 6.
		omonth = MonthRollover - omonth
	}
	return uint32(omonth)
}

func (s *State) GetSunset() uint32 {
	omonth := s.getDayLengthFactor()
	return 16*OneHour + (2880000 * omonth)
}

func (s *State) GetSunrise() uint32 {
	omonth := s.getDayLengthFactor()
	return 8*OneHour - (2160000 * omonth)
}

func (s *State) IsDay() bool {
	return s.GetSunrise() < s.TimeMili && s.TimeMili < s.GetSunset()
}
