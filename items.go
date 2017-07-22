package gorl

import "github.com/japanoise/engutil"

type Item struct {
	Name      string
	Spr       Sprite
	Class     ItemClassID
	DamageAc  uint8 // if it's a piece of apparel, it's the AC, if it's a weapon, it's dice-damage
	Value     int
	MagEffect MagicID
	MagLevel  int8
	Bcu       BCU
}

type MagicID uint8

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
)

var ItemClassDir map[ItemClassID]*ItemClass

func init() {
	ItemClassDir = make(map[ItemClassID]*ItemClass)
	ItemClassDir[ItemClassCurrency] = &ItemClass{SpriteItemGold, "currency"}
	ItemClassDir[ItemClassWeapon] = &ItemClass{SpriteItemWeaponGeneric, "weapon"}
	ItemClassDir[ItemClassApp] = &ItemClass{SpriteItemAppGeneric, "apparel"}
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

func NewWeapon(name string, value int, mag MagicID, magl int8, bcu BCU, damageDice uint8) *Item {
	ret := NewItemOfClass(name, ItemClassWeapon)
	ret.Value = value
	ret.MagEffect = mag
	ret.MagLevel = magl
	ret.Bcu = bcu
	ret.DamageAc = damageDice
	return ret
}

func NewApparel(name string, value int, mag MagicID, magl int8, bcu BCU, ac uint8) *Item {
	ret := NewItemOfClass(name, ItemClassApp)
	ret.Value = value
	ret.MagEffect = mag
	ret.MagLevel = magl
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
