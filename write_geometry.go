package vt

import (
	m "github.com/murphy214/mercantile"
	"math"
	//pc "github.com/murphy214/polyclip"
	//"fmt"
)

const mercatorPole = 20037508.34

type Cursor struct {
	Geometry  []uint32
	LastPoint []int32
	Bounds    m.Extrema
	DeltaX    float64 // delta between bounds
	DeltaY    float64 // delta between bounds
	Count     uint32
	Extent    int32
	Bds       m.Extrema
}

var startbds = m.Extrema{N: -90.0, S: 90.0, E: -180.0, W: 180.0}

func NewCursor(tileid m.TileID) *Cursor {
	bound := m.Bounds(tileid)
	deltax := bound.E - bound.W
	deltay := bound.N - bound.S
	cur := Cursor{LastPoint: []int32{0, 0}, Bounds: bound, DeltaX: deltax, DeltaY: deltay, Count: 0, Extent: int32(4096), Bds: startbds}
	cur = ConvertCursor(cur)
	return &cur
}

func NewCursorExtent(tileid m.TileID, extent int32) *Cursor {
	bound := m.Bounds(tileid)
	deltax := bound.E - bound.W
	deltay := bound.N - bound.S
	cur := Cursor{LastPoint: []int32{0, 0}, Bounds: bound, DeltaX: deltax, DeltaY: deltay, Count: 0, Extent: extent, Bds: startbds}
	cur = ConvertCursor(cur)
	return &cur
}

func ConvertPoint(point []float64) []float64 {
	x := mercatorPole / 180.0 * point[0]

	y := math.Log(math.Tan((90.0+point[1])*math.Pi/360.0)) / math.Pi * mercatorPole
	y = math.Max(-mercatorPole, math.Min(y, mercatorPole))
	return []float64{x, y}
}

func cmdEnc(id uint32, count uint32) uint32 {
	return (id & 0x7) | (count << 3)
}

func moveTo(count uint32) uint32 {
	return cmdEnc(1, count)
}

func lineTo(count uint32) uint32 {
	return cmdEnc(2, count)
}

func closePath(count uint32) uint32 {
	return cmdEnc(7, count)
}

func paramEnc(value int32) int32 {
	return (value << 1) ^ (value >> 31)
}

func (cur *Cursor) MovePoint(point []int32) {
	cur.Geometry = append(cur.Geometry, moveTo(1))
	cur.Geometry = append(cur.Geometry, uint32(paramEnc(point[0]-cur.LastPoint[0])))
	cur.Geometry = append(cur.Geometry, uint32(paramEnc(point[1]-cur.LastPoint[1])))
	cur.LastPoint = point
}

func (cur *Cursor) LinePoint(point []int32) {
	deltax := point[0] - cur.LastPoint[0]
	deltay := point[1] - cur.LastPoint[1]
	if ((deltax == 0) && (deltay == 0)) == false {
		cur.Geometry = append(cur.Geometry, uint32(paramEnc(deltax)))
		cur.Geometry = append(cur.Geometry, uint32(paramEnc(deltay)))
		cur.Count = cur.Count + 1
	}
	cur.LastPoint = point
}

// makes a line pretty neatly
func (cur *Cursor) MakeLine(coords [][]int32) {
	// applying the first move to point
	startpos := len(cur.Geometry)
	cur.MovePoint(coords[0])
	cur.Geometry = append(cur.Geometry, lineTo(uint32(len(coords)-1)))

	// iterating through each point
	for _, point := range coords[1:] {
		cur.LinePoint(point)
	}

	cur.Geometry[startpos+3] = lineTo(cur.Count)
}

// makes a line pretty neatly
func (cur *Cursor) MakeLineFloat(coords [][]float64) {
	// applying the first move to point
	startpos := len(cur.Geometry)
	firstpoint := cur.SinglePoint(coords[0])
	cur.MovePoint(firstpoint)
	cur.Geometry = append(cur.Geometry, lineTo(uint32(len(coords)-1)))
	// iterating through each point
	for _, point := range coords[1:] {
		cur.LinePoint(cur.SinglePoint(point))
	}

	cur.Geometry[startpos+3] = lineTo(cur.Count)

	//return cur.Geometry

}

// reverses the coord list
func reverse(coord [][]int32) [][]int32 {
	current := len(coord) - 1
	newlist := [][]int32{}
	for current != -1 {
		newlist = append(newlist, coord[current])
		current = current - 1
	}
	return newlist
}

// asserts a winding order
func assert_winding_order(coord [][]int32, exp_orient string) [][]int32 {
	count := 0
	firstpt := coord[0]
	weight := 0.0
	var oldpt []int32
	for _, pt := range coord {
		if count == 0 {
			count = 1
		} else {
			weight += float64((pt[0] - oldpt[0]) * (pt[1] + oldpt[1]))
		}
		oldpt = pt
	}

	weight += float64((firstpt[0] - oldpt[0]) * (firstpt[1] + oldpt[1]))
	var orientation string
	if weight > 0 {
		orientation = "clockwise"
	} else {
		orientation = "counter"
	}

	if orientation != exp_orient {
		return reverse(coord)
	} else {
		return coord
	}
	return coord

}

// asserts a winding order
func (cur *Cursor) AssertConvert(coord [][]float64, exp_orient string) {
	count := 0
	firstpt := cur.SinglePoint(coord[0])
	weight := 0.0
	var oldpt []int32
	newlist := [][]int32{firstpt}

	// iterating through each float point
	for _, floatpt := range coord {
		pt := cur.SinglePoint(floatpt)
		newlist = append(newlist, pt)
		if count == 0 {
			count = 1
		} else {
			weight += float64((pt[0] - oldpt[0]) * (pt[1] + oldpt[1]))
		}
		oldpt = pt
	}

	weight += float64((firstpt[0] - oldpt[0]) * (firstpt[1] + oldpt[1]))
	var orientation string
	if weight > 0 {
		orientation = "clockwise"
	} else {
		orientation = "counter"
	}

	if orientation != exp_orient {
		newlist = reverse(newlist)
	}

	newcur := Cursor{LastPoint: cur.LastPoint, Bounds: cur.Bounds, DeltaX: cur.DeltaX, DeltaY: cur.DeltaY}
	newcur.MakeLine(newlist)
	newgeom := newcur.Geometry
	newgeom = append(newgeom, closePath(1))
	cur.Geometry = append(cur.Geometry, newgeom...)
	cur.LastPoint = newlist[len(newlist)-1]

}

// makes a polygon
func (cur *Cursor) MakePolygon(coords [][][]int32) []uint32 {
	// applying the first ring
	coord := coords[0]
	coord = assert_winding_order(coord, "clockwise")
	cur.MakeLine(coord)
	cur.Geometry = append(cur.Geometry, cur.Geometry...)
	cur.Geometry = append(cur.Geometry, closePath(1))

	// if multiple rings exist proceed to add those also
	if len(coords) > 1 {
		for _, coord := range coords[1:] {
			coord = assert_winding_order(coord, "counter")
			cur.MakeLine(coord)
			cur.Geometry = append(cur.Geometry, cur.Geometry...)
			cur.Geometry = append(cur.Geometry, closePath(1))

		}
	}

	return cur.Geometry
}

// makes a polygon
func (cur *Cursor) MakePolygonFloat(coords [][][]float64) {
	// applying the first ring
	cur.AssertConvert(coords[0], "clockwise")

	// if multiple rings exist proceed to add those also
	if len(coords) > 1 {
		for _, coord := range coords[1:] {
			cur.AssertConvert(coord, "counter")

		}
	}
	//return cur.Geometry
}

// converts a single point from a coordinate to a tile point
func (cur *Cursor) SinglePoint(point []float64) []int32 {
	if cur.Bounds.N < point[1] {
		cur.Bounds.N = point[1]
	} else if cur.Bounds.S > point[1] {
		cur.Bounds.S = point[1]
	}
	if cur.Bounds.E < point[0] {
		cur.Bounds.E = point[0]
	} else if cur.Bounds.W > point[0] {
		cur.Bounds.W = point[0]
	}
	// converting to sperical coordinates
	point = ConvertPoint(point)

	// getting factors to multiply by
	factorx := (point[0] - cur.Bounds.W) / cur.DeltaX
	factory := (cur.Bounds.N - point[1]) / cur.DeltaY

	xval := int32(factorx * float64(cur.Extent))
	yval := int32(factory * float64(cur.Extent))

	if xval >= cur.Extent {
		xval = cur.Extent
	}

	if yval >= cur.Extent {
		yval = cur.Extent
	}

	if xval < 0 {
		xval = 0
	}
	if yval < 0 {
		yval = 0
	}

	return []int32{xval, yval}
}

func (cur *Cursor) MakePointFloat(point []float64) {
	newpoint := cur.SinglePoint(point)

	coords := []int32{newpoint[0], newpoint[1]}
	cur.Geometry = []uint32{moveTo(uint32(1))}
	cur.LinePoint(coords)

}

func (cur *Cursor) MakeMultiPointFloat(points [][]float64) {
	cur.Geometry = []uint32{moveTo(uint32(len(points)))}
	for _, point := range points {
		newpoint := cur.SinglePoint(point)
		cur.LinePoint(newpoint)
	}
}

func (cur *Cursor) MakeMultiLineFloat(lines [][][]float64) {
	for _, line := range lines {
		cur.MakeLineFloat(line)
	}
}

func (cur *Cursor) MakeMultiPolygonFloat(lines [][][][]float64) {
	for _, line := range lines {
		cur.MakePolygonFloat(line)
	}
}

// converts a cursor to world points
func ConvertCursor(cur Cursor) Cursor {
	// getting bounds
	bounds := cur.Bounds

	// getting ne point
	en := []float64{bounds.E, bounds.N} // east, north point
	ws := []float64{bounds.W, bounds.S} // west, south point

	// converting these
	en = ConvertPoint(en)
	ws = ConvertPoint(ws)

	// gettting north east west south
	east := en[0]
	north := en[1]
	west := ws[0]
	south := ws[1]
	bounds = m.Extrema{N: north, E: east, S: south, W: west}

	// getting deltax and deltay
	deltax := east - west
	deltay := north - south

	// setting the new values
	cur.Bounds = bounds
	cur.DeltaX = deltax
	cur.DeltaY = deltay
	return cur
}
