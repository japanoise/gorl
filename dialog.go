package gorl

/* Implementation of https://gamedev.stackexchange.com/a/40524 with some changes
   to make it idiomatic Go */

type NodeID int

// Nodes in the 'tree'
type DialogNode struct {
	Lines func(*State, *Critter, *Critter) []string // NPC's response, greeting, etc
	Links []*DialogLink                             // Where we can go from here
}

// Directed link between nodes
type DialogLink struct {
	Lines      func(*State, *Critter, *Critter) string // What gets said
	To         NodeID                                  // Where this link goes
	Conditions []func(*State, *Critter, *Critter) bool // Conditions that have to be met before following this link
	Script     func(*State, *Critter, *Critter)        // What to do if this link is followed
}

var Nodes map[NodeID]*DialogNode

func init() {
	Nodes = make(map[NodeID]*DialogNode)
	Nodes[1] = &DialogNode{
		func(s *State, p, me *Critter) []string { return []string{me.GetTheName() + ": Well met."} },
		[]*DialogLink{&DialogLink{
			func(*State, *Critter, *Critter) string { return "Goodbye." },
			2,
			nil,
			nil,
		}},
	}
	Nodes[2] = &DialogNode{
		func(s *State, p, me *Critter) []string { return []string{me.GetTheName() + ": Farewell."} },
		nil,
	}
}

func (d *DialogLink) IsValid(s *State, p *Critter, me *Critter) bool {
	if d.Conditions == nil {
		return true
	}
	for _, c := range d.Conditions {
		if !c(s, p, me) {
			return false
		}
	}
	return true
}

func (d *DialogNode) HasLinks(s *State, p *Critter, me *Critter) bool {
	if d.Links == nil || len(d.Links) == 0 {
		return false
	}
	for _, link := range d.Links {
		if link.IsValid(s, p, me) {
			return true
		}
	}
	return false
}

func (d *DialogNode) Run(s *State, player *Critter, me *Critter) {
	for _, str := range d.Lines(s, player, me) {
		s.Out.ChatMessage(str)
	}
	if d.HasLinks(s, player, me) {
		ls := make([]*DialogLink, 0, len(d.Links))
		for _, lin := range d.Links {
			if lin.IsValid(s, player, me) {
				ls = append(ls, lin)
			}
		}
		choices := make([]string, len(ls))
		for i, lin := range ls {
			choices[i] = lin.Lines(s, player, me)
		}
		choice := s.Out.MenuIndex("Speaking to "+me.GetTheName(), choices)
		ls[choice].Follow(s, player, me).Run(s, player, me)
	}
}

func (l *DialogLink) Follow(s *State, player *Critter, me *Critter) *DialogNode {
	s.Out.ChatMessage(player.GetName() + ": " + l.Lines(s, player, me))
	if l.Script != nil {
		l.Script(s, player, me)
	}
	return GetNode(l.To)
}

func GetNode(id NodeID) *DialogNode {
	return Nodes[id]
}

func DoChat(state *State, player *Critter) {
	p := getEndPoint(state.CurLevel, player, state.In.GetDirection("Talk to who (in which direction)?"))
	x, y := p.GetXY()
	if state.CurLevel.OOB(x, y) {
		state.Out.Message("There's no-one there.")
		return
	}
	critter := state.CurLevel.Tiles[x][y].Here
	if critter == player {
		state.Out.Message("Talking to yourself is the first sign of madness...")
	} else if critter == nil {
		state.Out.Message("There's no-one there.")
	} else if critter.Dialog == 0 || !critter.HasFlags(FlagFriendly) {
		state.Out.Message(critter.GetTheName() + " doesn't seem interested in talking.")
	} else {
		GetNode(critter.Dialog).Run(state, player, critter)
	}
}
