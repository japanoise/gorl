package gorl

import (
	"fmt"
	"strconv"
)

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
			func(*State, *Critter, *Critter) string { return "Want to trade?" },
			2,
			[]func(*State, *Critter, *Critter) bool{func(s *State, p *Critter, me *Critter) bool {
				return me.HasFlags(FlagShopkeep|FlagInnkeep) && me.Inv != nil && len(me.Inv) != 0
			}},
			doShop,
		}, &DialogLink{
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
	Nodes[3] = &DialogNode{
		func(s *State, p, me *Critter) []string {
			return []string{me.GetTheName() + ": A pleasure doing business with you."}
		},
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

func doShop(s *State, player *Critter, shopkeep *Critter) {
	prompt := "How can I help you?"
	for {
		choice := s.Out.MenuIndex(prompt, []string{"Buy", "Sell", "Done"})
		switch choice {
		case 0:
			doBuy(s, player, shopkeep)
		case 1:
			doSell(s, player, shopkeep)
		case 2:
			return
		}
		prompt = "Is there anything else I can help you with?"
	}
}

func doBuy(s *State, player *Critter, shopkeep *Critter) {
	prompt := "What do you want to buy?"
	for {
		choices := make([]string, len(shopkeep.Inv)+1)
		choices[0] = "Done."
		for i, item := range shopkeep.Inv {
			if item.Items[0].Class == ItemClassMercGen {
				choices[i+1] = ""
			} else {
				choices[i+1] = strconv.Itoa(len(item.Items)) + "x "
			}
			choices[i+1] += item.Items[0].DescribeExtra()
		}
		choice := s.Out.MenuIndex(fmt.Sprintf("%s [%d AU]", prompt, player.Gold), choices)
		if choice == 0 {
			return
		} else if player.Gold >= shopkeep.Inv[choice-1].Items[0].Value {
			sell := shopkeep.Inv[choice-1]
			player.Gold -= sell.Items[0].Value
			shopkeep.Gold += sell.Items[0].Value
			if sell.Items[0].Class == ItemClassMercGen {
				player.AddInventoryItem(sell.Items[0].GenMerchItem())
			} else {
				player.AddInventoryItem(sell.Items[0])
				shopkeep.DeleteOneInvItem(choice - 1)
			}
			prompt = "Thanks! Anything else?"
		} else {
			prompt = "You can't afford that."
		}
	}
}

func doSell(s *State, player *Critter, shopkeep *Critter) {
	prompt := "What do you want to sell?"
	for {
		choices := make([]string, len(player.Inv)+1)
		choices[0] = "Done."
		for i, item := range player.Inv {
			if item.Items[0].Class == ItemClassMercGen {
				choices[i+1] = ""
			} else {
				choices[i+1] = strconv.Itoa(len(item.Items)) + "x "
			}
			choices[i+1] += item.Items[0].DescribeExtra()
		}
		choice := s.Out.MenuIndex(fmt.Sprintf("%s [Player: %d AU %s: %d AU]", prompt,
			player.Gold, shopkeep.GetTheName(), shopkeep.Gold), choices)
		if choice == 0 {
			return
		} else if shopkeep.Gold >= player.Inv[choice-1].Items[0].Value {
			sell := player.Inv[choice-1]
			shopkeep.Gold -= sell.Items[0].Value
			player.Gold += sell.Items[0].Value
			shopkeep.AddInventoryItem(sell.Items[0])
			player.DeleteOneInvItem(choice - 1)
			prompt = "Thanks! Anything else?"
		} else {
			prompt = "I can't afford that."
		}
	}
}
