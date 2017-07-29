package gorl

import (
	"fmt"

	"github.com/japanoise/engutil"
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
)

var ItemClassDir map[ItemClassID]*ItemClass

func initItems() error {
	ItemClassDir = make(map[ItemClassID]*ItemClass)
	ItemClassDir[ItemClassCurrency] = &ItemClass{SpriteItemGold, "currency"}
	ItemClassDir[ItemClassWeapon] = &ItemClass{SpriteItemWeaponGeneric, "weapon"}
	ItemClassDir[ItemClassApp] = &ItemClass{SpriteItemAppGeneric, "apparel"}
	ItemClassDir[ItemClassPotion] = &ItemClass{SpriteItemPotion, "potion"}
	ItemClassDir[ItemClassFood] = &ItemClass{SpriteItemFoodGeneric, "food"}
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
	}
	return ret
}

func ShowItemList(g Graphics, gold int, items []*Item) *Item {
	choices := make([]string, len(items)+1)
	choices[0] = "<Cancel>"
	for i := range items {
		choices[i+1] = items[i].DescribeExtra()
	}
	choice := g.MenuIndex(fmt.Sprintf("Inventory (%d gold)", gold), choices)
	if choice == 0 {
		return nil
	} else {
		return items[choice-1]
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
	} else if item.Class == ItemClassWeapon {
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
			if invitem == item {
				player.Inv[i] = store
			}
		}
	} else {
		delindex := -1
		for i, invitem := range player.Inv {
			if invitem == item {
				delindex = i
			}
		}
		if delindex == -1 {
			return
		}
		player.Inv[delindex] = player.Inv[len(player.Inv)-1]
		player.Inv[len(player.Inv)-1] = nil
		player.Inv = player.Inv[:len(player.Inv)-1]
	}
}
