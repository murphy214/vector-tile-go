package vt

import (
	m "github.com/murphy214/mercantile"
	"math"
	//pc "github.com/murphy214/polyclip"
)

const mercatorPole = 20037508.34

// the cursor structure for rendering a geometry 
// admittedly to much happens here, but I figure its more logical to 
// have a master geometry structure that can be dictated by hgiher order layer or feature structures
type Cursor struct {
	Geometry   []uint32 // holds the geomtry
	LastPoint  []int32 // holds the last point
	Elevations []uint32 // holds elevations
	Bounds     m.Extrema // extrema structure
	DeltaX     float64 // delta between bounds
	DeltaY     float64 // delta between bounds
	Count      uint32 // arbitary counter?
	Extent     int32 // the extent 
	Bds        m.Extrema  
	ExtentBool bool // bool for whether or not to trim by the extent
	ZBool bool // zbool to indicate whether to render the third dimmension
	SplineKnots []int // the spline knots currently not supported
	SplineDegree int // the spline degree currently not supported
	Scaling *Scaling // the scaling structure for the given feature
	CurrentElevation float64 // the current elevation
	IsTrimmed bool // boolean for whether or not a polygon is trimmed.
	GeometricAttributesBool bool // geometric attributes map
	Position int // arbitary position
	GeometricAttributesMap map[string][]interface{} // the original geomtirc attributes map
	NewGeometricAttributesMap map[string][]interface{} // the new geometric attributes map
	MovePointBool bool // a bool for the first move point
	AttributePosition int // the attribute position we are currently on
	ClosePathAttributePosition int // the position of the geometric attribute list when close path is used
}

var startbds = m.Extrema{N: -90.0, S: 90.0, E: -180.0, W: 180.0}

func TrimPolygonFloat(lines [][][]float64) [][][]float64 {
	for pos, line := range lines {
		f, l := line[0], line[len(line)-1]
		if !(f[0] == l[0] && l[1] == f[1]) {
			line = append(line, line[0])
		}
		lines[pos] = line
	}
	return lines
}

func TrimPolygon(lines [][][]int32) [][][]int32 {
	for pos, line := range lines {
		f, l := line[0], line[len(line)-1]
		if !(f[0] == l[0] && l[1] == f[1]) {
			line = append(line, line[0])
		}
		lines[pos] = line
	}
	return lines
}

func TrimMultiPolygonFloat(polygons [][][][]float64) [][][][]float64 {
	for pos, polygon := range polygons {
		polygons[pos] = TrimPolygonFloat(polygon)
	}
	return polygons
}

func TrimMultiPolygon(polygons [][][][]int32) [][][][]int32 {
	for pos, polygon := range polygons {
		polygons[pos] = TrimPolygon(polygon)
	}
	return polygons
}

// creates a new cursor about a tile
func NewCursor(tileid m.TileID) *Cursor {
	bound := m.Bounds(tileid)
	deltax := bound.E - bound.W
	deltay := bound.N - bound.S
	cur := Cursor{
		LastPoint: []int32{0, 0},
		MovePointBool:true, 
		Bounds: bound, 
		DeltaX: deltax, 
		DeltaY: deltay, 
		Count: 0, 
		Extent: int32(4096),
		Bds: startbds,
		NewGeometricAttributesMap:map[string][]interface{}{},
		GeometricAttributesMap:map[string][]interface{}{},
	}
	cur = ConvertCursor(cur)
	return &cur
}

// creates a cursor about the extent
func NewCursorExtent(tileid m.TileID, extent int32) *Cursor {
	bound := m.Bounds(tileid)
	deltax := bound.E - bound.W
	deltay := bound.N - bound.S
	cur := Cursor{
		LastPoint: []int32{0, 0},
		MovePointBool:true,
		Bounds: bound, 
		DeltaX: deltax, 
		DeltaY: deltay, 
		Count: 0, 
		Extent: extent, 
		Bds: startbds,
		NewGeometricAttributesMap:map[string][]interface{}{},
		GeometricAttributesMap:map[string][]interface{}{},	
	}
	cur = ConvertCursor(cur)
	return &cur
}

// sets the cursors geometric attributes
func (cur *Cursor) SetCursorGeometricAttributes(attrs map[string][]interface{}) {
	cur.GeometricAttributesMap = attrs
	cur.GeometricAttributesBool = len(attrs) > 0 
}

// dumps the interface map to the proper format
func DumpInterfaceMap(data interface{}) map[string]interface{} {
	mine,boolval := data.(map[string]interface{})
	if boolval {
		return mine
	} else {
		return map[string]interface{}{}
	}
} 



func ConvertPoint(point []float64) []float64 {
	x := mercatorPole / 180.0 * point[0]

	y := math.Log(math.Tan((90.0+point[1])*math.Pi/360.0)) / math.Pi * mercatorPole
	y = math.Max(-mercatorPole, math.Min(y, mercatorPole))
	if len(point)==3 {
		return []float64{x, y,point[2]}
	} else {
		return []float64{x,y}
	}
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

// gets the elevation
func (cur *Cursor) Elevation(value float64) uint32 {
	delta_encoded_value :=  float64((value - cur.Scaling.Base) / cur.Scaling.Multiplier) - float64(cur.Scaling.Offset)
	return uint32(paramEnc(int32(delta_encoded_value)))
}

// adding an attribute to attributes list
func (cur *Cursor) AddAttribute() {
	for k,v := range cur.GeometricAttributesMap {
		cur.NewGeometricAttributesMap[k] = append(cur.NewGeometricAttributesMap[k],v[cur.AttributePosition])
	}	
}

// simple move to command
func (cur *Cursor) MovePoint(point []int32) {
	cur.Geometry = append(cur.Geometry, moveTo(1))
	cur.Geometry = append(cur.Geometry, uint32(paramEnc(point[0]-cur.LastPoint[0])))
	cur.Geometry = append(cur.Geometry, uint32(paramEnc(point[1]-cur.LastPoint[1])))
	if cur.ZBool {
		cur.Elevations = append(cur.Elevations,cur.Elevation(cur.CurrentElevation))
	}
	if cur.GeometricAttributesBool {
		cur.AddAttribute()
	}
	cur.Position++
	cur.AttributePosition++
	cur.LastPoint = point
	cur.Count = 0
}

// simple line to command
func (cur *Cursor) LinePoint(point []int32) {
	deltax := point[0] - cur.LastPoint[0]
	deltay := point[1] - cur.LastPoint[1]
	if ((deltax == 0) && (deltay == 0)) == false {
		cur.Geometry = append(cur.Geometry, uint32(paramEnc(deltax)))
		cur.Geometry = append(cur.Geometry, uint32(paramEnc(deltay)))
		if cur.ZBool {
			cur.Elevations = append(cur.Elevations,cur.Elevation(cur.CurrentElevation))
		}
		if cur.GeometricAttributesBool {
			cur.AddAttribute()
		}
		cur.Count = cur.Count + 1
		
	}
	cur.AttributePosition++
	cur.Position++
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
	//fmt.Println(lineTo(cur.Count), coords, cur.Count, len(coords), cur.Geometry)
	cur.Geometry[startpos+3] = lineTo(cur.Count)

	//return cur.Geometry

}

// reverses the coord list
func reverse(coord [][]int32) [][]int32 {
	current := len(coord) - 1
	newlist := make([][]int32,len(coord))
	i := 0
	for current != -1 {
		newlist[i] = coord[current]
		current = current - 1
		i++
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

var Power7 = math.Pow(10,-7)

// checking to see if the start / end value is trimmed
func IsTrimmed(p1,p2 []float64) bool {
	dx,dy := math.Abs(p1[0] - p2[0]),math.Abs(p1[1] - p2[1])
	return !(dx < Power7 && dy < Power7)
}


// asserts a winding order
func (cur *Cursor) AssertConvert(coord [][]float64, exp_orient string) {
	// checking to see if the end of the coord needs trimmed
	// trimming and storing last attribute
	cur.IsTrimmed = IsTrimmed(coord[0], coord[len(coord)-1])
	if !cur.IsTrimmed {
		coord = coord[:len(coord)-1]
		cur.ClosePathAttributePosition = cur.AttributePosition + len(coord) - 1
	} else {
		cur.ClosePathAttributePosition = cur.AttributePosition
	}

	count := 0
	firstpt := cur.SinglePoint(coord[0])
	weight := 0.0
	var oldpt []int32
	newlist := make([][]int32, len(coord))
	//newlist := [][]int32{firstpt}
	newlist[0] = firstpt
	// iterating through each float point
	for pos, floatpt := range coord[1:] {
		pt := cur.SinglePoint(floatpt)
		newlist[pos+1] = pt
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

	// creating cursor to make the line
	newcur := Cursor{
		LastPoint: cur.LastPoint,
		Bounds: cur.Bounds,
		DeltaX: cur.DeltaX,
		DeltaY: cur.DeltaY,
		CurrentElevation:cur.CurrentElevation,
		AttributePosition:cur.AttributePosition,
		GeometricAttributesBool:cur.GeometricAttributesBool,
		GeometricAttributesMap:cur.GeometricAttributesMap,
		NewGeometricAttributesMap:cur.NewGeometricAttributesMap,
	}
	newcur.MakeLine(newlist)
	newgeom,neweles := newcur.Geometry,newcur.Elevations
	newgeom = append(newgeom, closePath(1))

	// adding the close path interval
	if !cur.IsTrimmed {
		cur.AddAttribute()
		cur.AttributePosition++
	} else {
		cur.AddAttribute()
	}

	// cleaning up
	cur.Position++
	cur.Geometry = append(cur.Geometry, newgeom...)
	cur.Elevations = append(cur.Elevations,neweles...)
	cur.LastPoint = newlist[len(newlist)-1]
}

// makes a polygon
func (cur *Cursor) MakePolygon(coords [][][]int32) []uint32 {
	coords = TrimPolygon(coords)

	// applying the first ring
	coord := coords[0]
	coord = assert_winding_order(coord, "clockwise")
	cur.MakeLine(coord)
	//cur.Geometry = append(cur.Geometry, cur.Geometry...)
	cur.Geometry = append(cur.Geometry, closePath(1))
	//cur.GeometricAttributesIndexes = append(cur.GeometricAttributesIndexes)
	cur.Position++
	// if multiple rings exist proceed to add those also
	if len(coords) > 1 {
		for _, coord := range coords[1:] {
			coord = assert_winding_order(coord, "counter")
			newcur := Cursor{LastPoint: cur.LastPoint, Bounds: cur.Bounds, DeltaX: cur.DeltaX, DeltaY: cur.DeltaY}
			newcur.MakeLine(coord)
			newgeom := newcur.Geometry
			newgeom = append(newgeom, closePath(1))
			//cur.GeometricAttributesIndexes = append(cur.GeometricAttributesIndexes)
			cur.Position++
			cur.Geometry = append(cur.Geometry, newgeom...)
			cur.LastPoint = coord[len(coord)-1]
		}
	}

	return cur.Geometry
}

// makes a polygon
func (cur *Cursor) MakePolygonFloat(coords [][][]float64) {
	coords = TrimPolygonFloat(coords)
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
	if len(point) == 3 && cur.MovePointBool {
		cur.ZBool = true
		cur.MovePointBool = false
	}
	if cur.ZBool {
		cur.CurrentElevation = point[2]
	}
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

	if cur.ExtentBool {
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
	}

	return []int32{xval, yval,int32(cur.CurrentElevation)}
}

func (cur *Cursor) MakePointFloat(point []float64) {
	newpoint := cur.SinglePoint(point)
	var coords []int32
	if cur.ZBool {
		cur.CurrentElevation = point[2]
		coords = []int32{newpoint[0], newpoint[1]}
	} else {
		coords = []int32{newpoint[0], newpoint[1]}
	}
	cur.Geometry = []uint32{moveTo(uint32(1))}
	cur.LinePoint(coords)

}

func (cur *Cursor) MakePoint(point []int32) {
	cur.Geometry = []uint32{moveTo(uint32(1))}
	cur.LinePoint(point)
}

func (cur *Cursor) MakeMultiPointFloat(points [][]float64) {
	cur.Geometry = []uint32{moveTo(uint32(len(points)))}
	for _, point := range points {
		newpoint := cur.SinglePoint(point)
		cur.LinePoint(newpoint)
	}
}

func (cur *Cursor) MakeMultiPoint(points [][]int32) {
	cur.Geometry = []uint32{moveTo(uint32(len(points)))}
	for _, point := range points {
		cur.LinePoint(point)
	}
}

func (cur *Cursor) MakeMultiLineFloat(lines [][][]float64) {
	for _, line := range lines {
		cur.MakeLineFloat(line)
	}
}

func (cur *Cursor) MakeMultiLine(lines [][][]int32) {
	for _, line := range lines {
		cur.MakeLine(line)
	}
}

func (cur *Cursor) MakeMultiPolygonFloat(lines [][][][]float64) {
	//lines = TrimMultiPolygonFloat(lines)
	for _, line := range lines {
		cur.MakePolygonFloat(line)
	}
}

func (cur *Cursor) MakeMultiPolygon(lines [][][][]int32) {
	//lines = TrimMultiPolygon(lines)

	for _, line := range lines {
		cur.MakePolygon(line)
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
