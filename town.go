package gorl

import "math/rand"

// Generate a cute little village. Similar-ish to dungen.
func VilGen(seed int64) (*Map, []*Critter) {
	sizex, sizey := 100, 100
	ret := GetBlankMap(-2, sizex, sizey)
	for x := 0; x < sizex; x++ {
		for y := 0; y < sizey; y++ {
			ret.SetGrassTile(x, y)
		}
	}
	r := rand.New(rand.NewSource(seed))
	numroomsx := 5
	numroomsy := 5
	roomsx := sizex / numroomsx
	roomsy := sizey / numroomsy
	spawns := make([]SpawnRegion, 0)
	// Draw rooms
	for xrn := 0; xrn < numroomsx; xrn++ {
		for yrn := 0; yrn < numroomsy; yrn++ {
			anchorx := xrn * roomsx
			anchory := yrn * roomsy
			var posx, posy, cornerx, cornery int
			posx = anchorx + r.Intn(roomsx/3)
			posy = anchory + r.Intn(roomsy/3)
			cornerx = (anchorx + roomsx - 2) - r.Intn(roomsx/3)
			cornery = (anchory + roomsy - 2) - r.Intn(roomsy/3)
			spawns = append(spawns, SpawnRegion{
				posx + 1, posy + 1,
				cornerx - 1, cornery - 1,
			})
			for x := posx; x <= cornerx; x++ {
				for y := posy; y <= cornery; y++ {
					if x == posx || y == posy || x == cornerx || y == cornery {
						ret.Tiles[x][y] = WallTile()
					} else {
						ret.Tiles[x][y] = FloorTile()
					}
				}
			}
			// Make sure we can get inside the place
			ret.Tiles[posx+3][cornery] = FloorTile()
		}
	}
	// Place an inn
	inn := spawns[r.Intn(len(spawns))]
	innkeep := NewFriendly(r, FlagInnkeep)
	innkeep.GenerateName()
	innkeep.Inv = []*InvItem{
		NewMerch(ItemClassFood, 10, "Cornish Pasty"),
	}
	PlaceCritterInRoom(innkeep, ret, inn)
	ret.Tiles[inn.TopLeftX+2][inn.TopLeftY-1] = TileOfClass(TileLetterI)
	ret.Tiles[inn.TopLeftX+3][inn.TopLeftY-1] = TileOfClass(TileLetterN)
	ret.Tiles[inn.TopLeftX+4][inn.TopLeftY-1] = TileOfClass(TileLetterN)
	return ret, []*Critter{innkeep}
}
