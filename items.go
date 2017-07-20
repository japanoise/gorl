package gorl

type Item struct {
	Name  string
	Spr   Sprite
	Class ItemClassID
}

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
	return &Item{name, ItemClassDir[class].Spr, class}
}

func GetGoldCoin() *Item {
	return NewItemOfClass("gold coin", ItemClassCurrency)
}
