package gorl

import (
	"bytes"
	"fmt"
	"math"
)

type DMapNum uint16

const DMAPNUM_MAX = math.MaxUint16 - 10 // to prevent overflow

// Brogue-style 'Dijkstra' map
type DijkstraMap struct {
	Points [][]DMapNum
	M      *Map
}

type DMapPoint struct {
	X   int
	Y   int
	Val DMapNum
}

func (d *DMapPoint) GetXY() (int, int) {
	return d.X, d.Y
}

func BlankDMap(m *Map) *DijkstraMap {
	ret := make([][]DMapNum, m.SizeX)
	for i := range ret {
		ret[i] = make([]DMapNum, m.SizeY)
		for j := range ret[i] {
			ret[i][j] = DMAPNUM_MAX
		}
	}
	return &DijkstraMap{ret, m}
}

func (d *DijkstraMap) Calc(points ...Placed) {
	for _, point := range points {
		x, y := point.GetXY()
		d.Points[x][y] = 0
	}
	mademutation := true
	for mademutation {
		mademutation = false
		for x := range d.Points {
			for y := range d.Points[x] {
				if d.M.Tiles[x][y].IsPassable() {
					ln := d.LowestNeighbour(x, y).Val
					if d.Points[x][y] > ln+1 {
						d.Points[x][y] = ln + 1
						mademutation = true
					}
				}
			}
		}
	}
}

func (d *DijkstraMap) GetValPoint(x, y int) DMapPoint {
	if d.M.OOB(x, y) {
		return DMapPoint{x, y, DMAPNUM_MAX}
	} else {
		return DMapPoint{x, y, d.Points[x][y]}
	}
}

func (d *DijkstraMap) LowestNeighbour(x, y int) DMapPoint {
	vals := []DMapPoint{
		d.GetValPoint(x+1, y),
		d.GetValPoint(x-1, y),
		d.GetValPoint(x, y-1),
		d.GetValPoint(x, y+1),
		d.GetValPoint(x+1, y+1),
		d.GetValPoint(x+1, y-1),
		d.GetValPoint(x-1, y+1),
		d.GetValPoint(x-1, y-1),
	}
	var lv DMapNum = DMAPNUM_MAX
	ret := vals[0]
	for _, val := range vals {
		if val.Val < lv {
			lv = val.Val
			ret = val
		}
	}
	return ret
}

func (d *DijkstraMap) String() string {
	buf := bytes.Buffer{}
	for x := range d.Points {
		for y := range d.Points[x] {
			buf.WriteString(fmt.Sprintf("%6d", d.Points[x][y]))
			buf.WriteString(", ")
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}
