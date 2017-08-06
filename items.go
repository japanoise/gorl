package gorl

import (
	"fmt"
	"strconv"

	"github.com/japanoise/engutil"
	"github.com/ulule/deepcopier"
)

type Item struct {
	Name     string
	Spr      Sprite
	Class    ItemClassID
	DamageAc uint8 // if it's a piece of apparel, it's the AC, if it's a weapon, it's dice-damage
	Value    int
	Bcu      BCU
	Magic    *Spell
}

type BCU uint8

const (
	Uncursed BCU = iota
	Cursed
	Blessed
)

type ItemClass struct {
	Spr  Sprite
	Name string
}

type ItemClassID uint8

const (
	ItemClassCurrency ItemClassID = iota
	ItemClassWeapon
	ItemClassApp
	ItemClassPotion
	ItemClassFood
	ItemClassMercGen
	ItemClassRanged
	ItemClassAmmo
)

type InvItem struct {
	Items []*Item
}

var ItemClassDir map[ItemClassID]*ItemClass

func initItems() error {
	ItemClassDir = make(map[ItemClassID]*ItemClass)
	ItemClassDir[ItemClassCurrency] = &ItemClass{SpriteItemGold, "currency"}
	ItemClassDir[ItemClassWeapon] = &ItemClass{SpriteItemWeaponGeneric, "weapon"}
	ItemClassDir[ItemClassApp] = &ItemClass{SpriteItemAppGeneric, "apparel"}
	ItemClassDir[ItemClassPotion] = &ItemClass{SpriteItemPotion, "potion"}
	ItemClassDir[ItemClassFood] = &ItemClass{SpriteItemFoodGeneric, "food"}
	ItemClassDir[ItemClassMercGen] = &ItemClass{SpriteItemFoodGeneric, "mercantile generator"}
	ItemClassDir[ItemClassRanged] = &ItemClass{SpriteItemWeaponGeneric, "ranged weapon"}
	ItemClassDir[ItemClassAmmo] = &ItemClass{SpriteItemAmmo, "ammo"}
	return nil
}

func NewItemOfClass(name string, class ItemClassID) *Item {
	ret := &Item{}
	ret.Name = name
	ret.Class = class
	ret.Spr = ItemClassDir[class].Spr
	return ret
}

func GetGoldCoins(value int) *Item {
	ret := NewItemOfClass("gold coin", ItemClassCurrency)
	ret.Value = value
	return ret
}

func (i *Item) DoDamage() int {
	if i.Class == ItemClassWeapon {
		return SmallDiceRoll(i.DamageAc)
	} else {
		return 1
	}
}

func (i *Item) DoRangedDamage() int {
	if i.Class == ItemClassRanged {
		val := SmallDiceRoll(i.DamageAc)
		debug.Println(val)
		return val
	} else {
		return 1
	}
}

func (i *Item) GetAC() uint8 {
	if i.Class == ItemClassApp {
		return i.DamageAc
	} else {
		return 0
	}
}

func NewWeapon(name string, value int, bcu BCU, damageDice uint8) *Item {
	ret := NewItemOfClass(name, ItemClassWeapon)
	ret.Value = value
	ret.Bcu = bcu
	ret.DamageAc = damageDice
	return ret
}

func NewApparel(name string, value int, bcu BCU, ac uint8) *Item {
	ret := NewItemOfClass(name, ItemClassApp)
	ret.Value = value
	ret.Bcu = bcu
	ret.DamageAc = ac
	return ret
}

func (i *Item) Describe() string {
	switch i.Class {
	case ItemClassCurrency:
		return engutil.Numbered(i.Name, i.Value)
	default:
		return i.Name
	}
}

func (i *Item) DescribeExtra() string {
	ret := i.Describe()
	switch i.Class {
	case ItemClassWeapon:
		ret += fmt.Sprintf(" [%s]", GetSmallDiceString(i.DamageAc))
	case ItemClassApp:
		ret += fmt.Sprintf(" [AC %d]", i.DamageAc)
	case ItemClassMercGen:
		ret = "Large piles of " + ret
	}
	return ret
}

func (i *Item) GenMerchItem() *Item {
	var ret *Item = &Item{}
	deepcopier.Copy(i).To(ret)
	ret.Class = ItemClassID(ret.Bcu)
	ret.Bcu = Uncursed
	return ret
}

func ShowItemList(g Graphics, prompt string, items []*InvItem) (*Item, int) {
	choices := make([]string, len(items)+1)
	choices[0] = "<Cancel>"
	for i := range items {
		choices[i+1] = strconv.Itoa(len(items[i].Items)) + "x " +
			items[i].Items[0].DescribeExtra()
	}
	choice := g.MenuIndex(prompt, choices)
	if choice == 0 {
		return nil, -1
	} else {
		return items[choice-1].Items[0], choice - 1
	}
}

func UseItem(state *State, player *Critter, item *Item) {
	g := state.Out
	if item == nil {
		return
	}
	var store *Item = nil
	if item.Class == ItemClassApp {
		if player.Armor != nil {
			g.Message("You remove " + item.DescribeExtra())
			store = player.Armor
		}
		g.Message("You don " + item.DescribeExtra())
		player.Armor = item
	} else if item.Class == ItemClassWeapon || item.Class == ItemClassRanged {
		if player.Weapon != nil {
			g.Message("You stow " + item.DescribeExtra())
			store = player.Weapon
		}
		g.Message("You ready " + item.DescribeExtra())
		player.Weapon = item
	} else if item.Class == ItemClassPotion {
		g.Message("You quaff " + item.Name)
		DoCastSpell(state, player, state.CurLevel, item.Magic)
	} else if item.Class == ItemClassFood {
		g.Message("You munch the " + item.Name)
		state.Player.TimeSinceEaten = 0
		state.Player.Hunger = HungerNormal
	} else {
		g.Message("There doesn't seem to be much use for that item.")
		return
	}
	if store != nil {
		for i, invitem := range player.Inv {
			if invitem.Items[0].SameAs(item) && len(invitem.Items) == 1 {
				player.Inv[i] = NewInvItem(store, 1)
				return
			}
		}
		player.AddInventoryItem(store)
	} else {
		delindex := -1
		for i, invitem := range player.Inv {
			if invitem.Items[0].SameAs(item) {
				delindex = i
			}
		}
		if delindex == -1 {
			return
		}
		player.DeleteOneInvItem(delindex)
	}
}

func (this *Item) SameAs(other *Item) bool {
	if this.Name != other.Name {
		return false
	} else if this.DescribeExtra() != other.DescribeExtra() {
		return false
	}
	return true
}

func Inventory(state *State, player *Critter) {
	choice, _ := ShowItemList(state.Out, fmt.Sprintf("Inventory (%d gold)", player.Gold), player.Inv)
	UseItem(state, player, choice)
}

func (c *Critter) AddInventoryItem(item *Item) {
	for i, items := range c.Inv {
		if items != nil && items.Items[0].SameAs(item) {
			c.Inv[i].Items = append(c.Inv[i].Items, item)
			return
		}
	}
	c.Inv = append(c.Inv, NewInvItem(item, 1))
}

func NewMerch(id ItemClassID, value int, name string) *InvItem {
	ret := NewInvItem(NewItemOfClass(name, ItemClassMercGen), 1)
	ret.Items[0].Value = value
	ret.Items[0].Bcu = BCU(id)
	return ret
}

func NewInvItem(item *Item, quantity int) *InvItem {
	ret := make([]*Item, quantity)
	ret[0] = item
	if quantity != 1 {
		for i := 1; i < quantity; i++ {
			p2 := &Item{}
			deepcopier.Copy(item).To(p2)
			ret[i] = p2
		}
	}
	return &InvItem{ret}
}

func (c *Critter) DeleteOneInvItem(delindex int) {
	if len(c.Inv[delindex].Items) == 1 {
		c.Inv[delindex] = c.Inv[len(c.Inv)-1]
		c.Inv[len(c.Inv)-1] = nil
		c.Inv = c.Inv[:len(c.Inv)-1]
	} else {
		ii := c.Inv[delindex]
		ii.Items[len(ii.Items)-1] = nil
		c.Inv[delindex].Items = ii.Items[:len(ii.Items)-1]
	}
}

func DropItem(state *State, player *Critter) {
	choice, choiceindex := ShowItemList(state.Out, "Drop which item?", player.Inv)
	if state.CurLevel.Tiles[player.X][player.Y].Items == nil {
		state.CurLevel.Tiles[player.X][player.Y].Items = []*Item{choice}
	} else {
		state.CurLevel.Tiles[player.X][player.Y].Items =
			append(state.CurLevel.Tiles[player.X][player.Y].Items, choice)
	}
	state.Out.Message("You drop " + choice.DescribeExtra())
	player.DeleteOneInvItem(choiceindex)
}

func Grab(state *State, player *Critter) {
	player.SnarfItems(state.CurLevel.Tiles[player.X][player.Y].Items)
	msg := "You take "
	for _, item := range state.CurLevel.Tiles[player.X][player.Y].Items {
		msg += item.DescribeExtra() + ","
	}
	state.Out.Message(msg)
	state.CurLevel.Tiles[player.X][player.Y].Items = []*Item{}
}

func Shoot(state *State, player *Critter) (*Critter, bool) {
	if player.Weapon == nil {
		state.Out.Message("You're not carrying a weapon!")
	} else if player.Weapon.Class != ItemClassRanged {
		state.Out.Message(player.Weapon.Describe() + " is not a ranged weapon!")
	} else {
		dir := state.In.GetDirection("Fire at what (which direction)?")
		p := getEndPoint(state.CurLevel, player, dir)
		if p.X == player.X && p.Y == player.Y {
			state.Out.Message("That's not a solution, " + player.GetName())
		} else if state.CurLevel.OOB(p.X, p.Y) {
			state.Out.Message("The arrow shoots off into the distance!")
		} else if state.CurLevel.Tiles[p.X][p.Y].Here == nil {
			state.Out.Message("The arrow strikes " + TilesDir[state.CurLevel.Tiles[p.X][p.Y].Id].Name)
		} else {
			return state.CurLevel.Tiles[p.X][p.Y].Here, RangedAttack(true, true, state, player, state.CurLevel.Tiles[p.X][p.Y].Here)
		}
	}
	return nil, false
}
