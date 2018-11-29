package vt

import (
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"math/rand"
	"testing"

)


// the test case structure
type TestCase struct {
	GeometryType string
	ZBool bool
	ClosedPolygon bool
	GeometricAttributesBool bool
	Point []float64
	MultiPoint [][]float64
	LineString [][]float64
	MultiLineString [][][]float64
	Polygon [][][]float64
	MultiPolygon [][][][]float64
	Cursor *Cursor
	CheckCursor *Cursor
}

func (testcase *TestCase) CreateTitle() string {
	var mydim string
	if testcase.ZBool {
		mydim = "3D"
	} else {
		mydim = "2D"
	}

	geoattr := ""
	if testcase.GeometricAttributesBool {
		geoattr = "With Geometric Attributes"
	} else {
		geoattr = "Without Geometric Attributes"
	}

	if testcase.GeometryType == "Polygon" || testcase.GeometryType == "MultiPolygon" {
		var closed string
		if testcase.ClosedPolygon {
			closed = "Closed Ending"
		} else {
			closed = "Open Ending"
		}
		return fmt.Sprintf("Test %s, %s, %s, %s",testcase.GeometryType,mydim,closed,geoattr)
	} else {
		return fmt.Sprintf("Test %s, %s, N/A, %s",testcase.GeometryType,mydim,geoattr)
	}
	return ""
}

// random elevation
func RandomInt() int {
	return rand.Intn(10000)
}

// creates a given cursor
func CreateCursor() *Cursor {
	cur := NewCursor(m.TileID{0,0,0})
	cur.Extent = int32(4096 * 4096)
	return cur
}

// adds z to a point
func AddZToPoint(pt []float64) []float64 {
	return append(pt,float64(RandomInt()))
}

// finishes the creation of the geometry and rolls over the 
// check cursor
func (testcase *TestCase) Convert() {
	switch testcase.GeometryType {
	case "Point":
		testcase.Cursor.MakePointFloat(testcase.Point)
	case "MultiPoint":
		testcase.Cursor.MakeMultiPointFloat(testcase.MultiPoint)
	case "LineString":
		testcase.Cursor.MakeLineFloat(testcase.LineString)
	case "MultiLineString":
		testcase.Cursor.MakeMultiLineFloat(testcase.MultiLineString)
	case "Polygon":	
		testcase.Cursor.MakePolygonFloat(testcase.Polygon)
	case "MultiPolygon":
		testcase.Cursor.MakeMultiPolygonFloat(testcase.MultiPolygon)
	}

	testcase.CheckCursor = &Cursor{
		Geometry:testcase.Cursor.Geometry,
		LastPoint:testcase.Cursor.LastPoint,
		Elevations:testcase.Cursor.Elevations,
		Bounds:testcase.Cursor.Bounds,
		DeltaX:testcase.Cursor.DeltaX,
		DeltaY:testcase.Cursor.DeltaY,
		Count:testcase.Cursor.Count,
		Extent:testcase.Cursor.Extent,
		Bds:testcase.Cursor.Bds,
		ExtentBool:testcase.Cursor.ExtentBool,
		ZBool:testcase.Cursor.ZBool,
		SplineKnots:testcase.Cursor.SplineKnots,
		SplineDegree:testcase.Cursor.SplineDegree,
		Scaling:testcase.Cursor.Scaling,
		CurrentElevation:testcase.Cursor.CurrentElevation,
		IsTrimmed:testcase.Cursor.IsTrimmed,
		GeometricAttributesBool:testcase.Cursor.GeometricAttributesBool,
		Position:testcase.Cursor.Position,
		GeometricAttributesMap:testcase.Cursor.GeometricAttributesMap,
		NewGeometricAttributesMap:testcase.Cursor.NewGeometricAttributesMap,
		MovePointBool:testcase.Cursor.MovePointBool,
		AttributePosition:testcase.Cursor.AttributePosition,
		ClosePathAttributePosition:testcase.Cursor.ClosePathAttributePosition,	
	}
	testcase.Cursor = CreateCursor()
	//testcase.Cursor.Elevations = []uint32{}
	//testcase.Cursor.Geometry = []uint32{}
	//testcase.Cursor.NewGeometricAttributesMap = map[string][]interface{}{}
}


// creates the ptcases
func ptcases(point []float64) []*TestCase {
	// 2d no attributes
	case1 := &TestCase{
		GeometryType:"Point",
		ZBool:false,
		GeometricAttributesBool:false,
		Point:point,
		Cursor:CreateCursor(),
	}
	case1.Convert()
	
	// 3d no attributes
	case2 := &TestCase{
		GeometryType:"Point",
		ZBool:true,
		GeometricAttributesBool:false,
		Point:AddZToPoint(point),
		Cursor:CreateCursor(),
	}
	case2.Convert()

	attrs := map[string][]interface{}{"field1":[]interface{}{55555}}
	// 2d attributes
	case3 := &TestCase{
		GeometryType:"Point",
		ZBool:false,
		GeometricAttributesBool:true,
		Point:point,
		Cursor:CreateCursor(),
	}
	case3.Cursor.GeometricAttributesMap = attrs
	case3.Convert()
	
	// 3d attributes
	case4 := &TestCase{
		GeometryType:"Point",
		ZBool:true,
		GeometricAttributesBool:true,
		Point:AddZToPoint(point),
		Cursor:CreateCursor(),
	}
	case4.Cursor.GeometricAttributesMap = attrs
	case4.Convert()
	return []*TestCase{case1,case2,case3,case4}
}

func Addztoline(line [][]float64) [][]float64 {
	for pos,i := range line {
		line[pos] = AddZToPoint(i)
	}
	return line
}


func Addztolines(lines [][][]float64) [][][]float64 {
	for pos,i := range lines {
		lines[pos] = Addztoline(i)
	}
	return lines
}

func Addztopolygons(lines [][][][]float64) [][][][]float64 {
	for pos,i := range lines {
		lines[pos] = Addztolines(i)
	}
	return lines
}

// creates the ptcases
func multiptcases(points [][]float64) []*TestCase {
	// 2d no attributes
	case1 := &TestCase{
		GeometryType:"MultiPoint",
		ZBool:false,
		GeometricAttributesBool:false,
		MultiPoint:points,
		Cursor:CreateCursor(),
	}
	case1.Convert()
	
	// 3d no attributes
	case2 := &TestCase{
		GeometryType:"MultiPoint",
		ZBool:true,
		GeometricAttributesBool:false,
		MultiPoint:Addztoline(points),
		Cursor:CreateCursor(),
	}
	case2.Convert()

	randattrs := make([]interface{},len(points))
	for i := range points {
		randattrs[i] = RandomInt()
	}

	attrs := map[string][]interface{}{"field1":randattrs}
	
	// 2d attributes
	case3 := &TestCase{
		GeometryType:"MultiPoint",
		ZBool:false,
		GeometricAttributesBool:true,
		MultiPoint:points,
		Cursor:CreateCursor(),
	}
	case3.Cursor.GeometricAttributesMap = attrs
	case3.Convert()
	
	// 3d attributes
	case4 := &TestCase{
		GeometryType:"MultiPoint",
		ZBool:true,
		GeometricAttributesBool:true,
		MultiPoint:Addztoline(points),
		Cursor:CreateCursor(),
	}
	case4.Cursor.GeometricAttributesMap = attrs
	case4.Convert()
	return []*TestCase{case1,case2,case3,case4}
}


// creates the ptcases
func linecases(points [][]float64) []*TestCase {
	// 2d no attributes
	case1 := &TestCase{
		GeometryType:"LineString",
		ZBool:false,
		GeometricAttributesBool:false,
		LineString:points,
		Cursor:CreateCursor(),
	}
	case1.Convert()
	
	// 3d no attributes
	case2 := &TestCase{
		GeometryType:"LineString",
		ZBool:true,
		GeometricAttributesBool:false,
		LineString:Addztoline(points),
		Cursor:CreateCursor(),
	}
	case2.Convert()

	randattrs := make([]interface{},len(points))
	for i := range points {
		randattrs[i] = RandomInt()
	}

	attrs := map[string][]interface{}{"field1":randattrs}
	
	// 2d attributes
	case3 := &TestCase{
		GeometryType:"LineString",
		ZBool:false,
		GeometricAttributesBool:true,
		LineString:points,
		Cursor:CreateCursor(),
	}
	case3.Cursor.GeometricAttributesMap = attrs
	case3.Convert()
	
	// 3d attributes
	case4 := &TestCase{
		GeometryType:"LineString",
		ZBool:true,
		GeometricAttributesBool:true,
		LineString:Addztoline(points),
		Cursor:CreateCursor(),
	}
	case4.Cursor.GeometricAttributesMap = attrs
	case4.Convert()
	return []*TestCase{case1,case2,case3,case4}
}


// creates the ptcases
func multilinecases(points [][][]float64) []*TestCase {
	// 2d no attributes
	case1 := &TestCase{
		GeometryType:"MultiLineString",
		ZBool:false,
		GeometricAttributesBool:false,
		MultiLineString:points,
		Cursor:CreateCursor(),
	}
	case1.Convert()
	
	// 3d no attributes
	case2 := &TestCase{
		GeometryType:"MultiLineString",
		ZBool:true,
		GeometricAttributesBool:false,
		MultiLineString:Addztolines(points),
		Cursor:CreateCursor(),
	}
	case2.Convert()


	randattrs := []interface{}{}
	for _,line := range points {
		for range line {
			randattrs = append(randattrs,RandomInt())
		}
	}

	attrs := map[string][]interface{}{"field1":randattrs}
	
	// 2d attributes
	case3 := &TestCase{
		GeometryType:"MultiLineString",
		ZBool:false,
		GeometricAttributesBool:true,
		MultiLineString:points,
		Cursor:CreateCursor(),
	}
	case3.Cursor.GeometricAttributesMap = attrs
	case3.Convert()
	
	// 3d attributes
	case4 := &TestCase{
		GeometryType:"MultiLineString",
		ZBool:true,
		GeometricAttributesBool:true,
		MultiLineString:Addztolines(points),
		Cursor:CreateCursor(),
	}
	case4.Cursor.GeometricAttributesMap = attrs
	case4.Convert()
	return []*TestCase{case1,case2,case3,case4}
}


// creates the ptcases
func polygoncases(points [][][]float64) []*TestCase {
	// 2d no attributes
	case1 := &TestCase{
		GeometryType:"Polygon",
		ZBool:false,
		GeometricAttributesBool:false,
		Polygon:points,
		Cursor:CreateCursor(),
		ClosedPolygon:true,
	}
	case1.Convert()
	
	// 3d no attributes
	case2 := &TestCase{
		GeometryType:"Polygon",
		ZBool:true,
		GeometricAttributesBool:false,
		Polygon:Addztolines(points),
		Cursor:CreateCursor(),
		ClosedPolygon:true,
	}
	case2.Convert()

	randattrs := []interface{}{}
	for _,line := range points {
		for range line {
			randattrs = append(randattrs,RandomInt())
		}
	}

	attrs := map[string][]interface{}{"field1":randattrs}
	
	// 2d attributes
	case3 := &TestCase{
		GeometryType:"Polygon",
		ZBool:false,
		GeometricAttributesBool:true,
		Polygon:points,
		Cursor:CreateCursor(),
		ClosedPolygon:true,
	}
	case3.Cursor.GeometricAttributesMap = attrs
	case3.Convert()
	
	// 3d attributes
	case4 := &TestCase{
		GeometryType:"Polygon",
		ZBool:true,
		GeometricAttributesBool:true,
		Polygon:Addztolines(points),
		Cursor:CreateCursor(),
		ClosedPolygon:true,
	}
	case4.Cursor.GeometricAttributesMap = attrs
	case4.Convert()
	for pos,i := range points {
		points[pos] = i[:len(i)-1]
	}

	// 2d no attributes
	case5 := &TestCase{
		GeometryType:"Polygon",
		ZBool:false,
		GeometricAttributesBool:false,
		Polygon:points,
		Cursor:CreateCursor(),
	}
	case5.Convert()
	
	// 3d no attributes
	case6 := &TestCase{
		GeometryType:"Polygon",
		ZBool:true,
		GeometricAttributesBool:false,
		Polygon:Addztolines(points),
		Cursor:CreateCursor(),
	}
	case6.Convert()

	// 2d attributes
	case7 := &TestCase{
		GeometryType:"Polygon",
		ZBool:false,
		GeometricAttributesBool:true,
		Polygon:points,
		Cursor:CreateCursor(),
	}
	case7.Cursor.GeometricAttributesMap = attrs
	case7.Convert()
	
	// 3d attributes
	case8 := &TestCase{
		GeometryType:"Polygon",
		ZBool:true,
		GeometricAttributesBool:true,
		Polygon:Addztolines(points),
		Cursor:CreateCursor(),
	}
	case8.Cursor.GeometricAttributesMap = attrs
	case8.Convert()

	return []*TestCase{case1,case2,case3,case4,case5,case6,case7,case8}
}


// creates the ptcases
func multipolygoncases(points [][][][]float64) []*TestCase {
	// 2d no attributes
	case1 := &TestCase{
		GeometryType:"MultiPolygon",
		ZBool:false,
		GeometricAttributesBool:false,
		MultiPolygon:points,
		Cursor:CreateCursor(),
		ClosedPolygon:true,
	}
	case1.Convert()
	
	// 3d no attributes
	case2 := &TestCase{
		GeometryType:"MultiPolygon",
		ZBool:true,
		GeometricAttributesBool:false,
		MultiPolygon:Addztopolygons(points),
		Cursor:CreateCursor(),
		ClosedPolygon:true,
	}
	case2.Convert()

	randattrs := []interface{}{}
	for _,line := range points {
		for range line {
			randattrs = append(randattrs,RandomInt())
		}
	}

	attrs := map[string][]interface{}{"field1":randattrs}
	
	// 2d attributes
	case3 := &TestCase{
		GeometryType:"MultiPolygon",
		ZBool:false,
		GeometricAttributesBool:true,
		MultiPolygon:points,
		Cursor:CreateCursor(),
		ClosedPolygon:true,
	}
	case3.Cursor.GeometricAttributesMap = attrs
	case3.Convert()
	
	// 3d attributes
	case4 := &TestCase{
		GeometryType:"MultiPolygon",
		ZBool:true,
		GeometricAttributesBool:true,
		MultiPolygon:Addztopolygons(points),
		Cursor:CreateCursor(),
		ClosedPolygon:true,
	}
	case4.Cursor.GeometricAttributesMap = attrs
	case4.Convert()
	for pos,i := range points {
		for pos2,ii := range i {
			i[pos2] = ii[:len(ii)-1]
		}
		points[pos] = i
		//points[pos] = i[:len(i)-1]
	}

	// 2d no attributes
	case5 := &TestCase{
		GeometryType:"MultiPolygon",
		ZBool:false,
		GeometricAttributesBool:false,
		MultiPolygon:points,
		Cursor:CreateCursor(),
	}
	case5.Convert()
	
	// 3d no attributes
	case6 := &TestCase{
		GeometryType:"MultiPolygon",
		ZBool:true,
		GeometricAttributesBool:false,
		MultiPolygon:Addztopolygons(points),
		Cursor:CreateCursor(),
	}
	case6.Convert()

	// 2d attributes
	case7 := &TestCase{
		GeometryType:"MultiPolygon",
		ZBool:false,
		GeometricAttributesBool:true,
		MultiPolygon:points,
		Cursor:CreateCursor(),
	}
	case7.Cursor.GeometricAttributesMap = attrs
	case7.Convert()
	
	// 3d attributes
	case8 := &TestCase{
		GeometryType:"MultiPolygon",
		ZBool:true,
		GeometricAttributesBool:true,
		MultiPolygon:Addztopolygons(points),
		Cursor:CreateCursor(),
	}
	case8.Cursor.GeometricAttributesMap = attrs
	case8.Convert()

	return []*TestCase{case1,case2,case3,case4,case5,case6,case7,case8}
}



// creates all the relevant test cases
func create() []*TestCase  {
	string1 := `{"id":1,"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-75.773786,39.7222],[-75.753228,39.757988999999995],[-75.71705899999999,39.792325],[-75.662846,39.821425],[-75.5943169052201,39.8345949730913],[-75.570433,39.839185],[-75.481207,39.829191],[-75.41506199999999,39.801919],[-75.459439,39.765813],[-75.47764,39.715013],[-75.509742,39.686113],[-75.535144,39.647211999999996],[-75.559446,39.629812],[-75.543965,39.596000000000004],[-75.512732,39.578],[-75.527676,39.535278],[-75.528088,39.498114],[-75.593068,39.479186],[-75.57182999999999,39.438897],[-75.521682,39.387871],[-75.50564316735289,39.370394560079596],[-75.58476499999999,39.308644],[-75.619631,39.310058],[-75.65115899999999,39.291593999999996],[-75.714901,39.299366],[-75.7604414164505,39.296789621100096],[-75.76689499999999,39.377499],[-75.76690460670919,39.3776515935512],[-75.788596,39.722198999999996],[-75.773786,39.7222]]]},"properties":{"AREA":"10003","COLORKEY":"#9DDE06","area":"10003","index":948}}`
	string2 := `{"id":0,"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-75.7604414164505,39.296789621100096],[-75.714901,39.299366],[-75.65115899999999,39.291593999999996],[-75.619631,39.310058],[-75.58476499999999,39.308644],[-75.50564316735289,39.370394560079596],[-75.469324,39.330819999999996],[-75.40837599999999,39.264697999999996],[-75.39479,39.188354],[-75.407473,39.133706],[-75.396277,39.057884],[-75.34089,39.01996],[-75.3066521095097,38.9476601633284],[-75.38016999999999,38.961892999999996],[-75.410463,38.916418],[-75.48412499999999,38.904447999999995],[-75.555013,38.835649],[-75.72310269333269,38.8298265565277],[-75.7481548142541,39.143131730959695],[-75.7564352155685,39.246687537125304],[-75.7604414164505,39.296789621100096]]]},"properties":{"AREA":"10001","COLORKEY":"#DDCD07","area":"10001","index":3143}}`
	feat1,_ := geojson.UnmarshalFeature([]byte(string1))
	feat2,_ := geojson.UnmarshalFeature([]byte(string2))
	pt := feat1.Geometry.Polygon[0][0]
	case1 := ptcases(pt)
	case2 := multiptcases(feat1.Geometry.Polygon[0])
	case3 := linecases(feat1.Geometry.Polygon[0])
	case4 := multilinecases(feat1.Geometry.Polygon)
	case5 := polygoncases(feat1.Geometry.Polygon)
	case6 := multipolygoncases([][][][]float64{feat1.Geometry.Polygon,feat2.Geometry.Polygon})
	total := []*TestCase{}
	total = append(total,case1...)
	total = append(total,case2...)
	total = append(total,case3...)
	total = append(total,case4...)
	total = append(total,case5...)
	total = append(total,case6...)
	return total
}

func FindError(testcase *TestCase) string {
	switch testcase.GeometryType {
	case "Point":
		testcase.Cursor.MakePointFloat(testcase.Point)
	case "MultiPoint":
		testcase.Cursor.MakeMultiPointFloat(testcase.MultiPoint)
	case "LineString":
		testcase.Cursor.MakeLineFloat(testcase.LineString)
	case "MultiLineString":
		testcase.Cursor.MakeMultiLineFloat(testcase.MultiLineString)
	case "Polygon":	
		testcase.Cursor.MakePolygonFloat(testcase.Polygon)
	case "MultiPolygon":
		testcase.Cursor.MakeMultiPolygonFloat(testcase.MultiPolygon)
	}

	val1,val2 := testcase.Cursor.Elevations,testcase.CheckCursor.Elevations
	if len(val1) != len(val2) {
		return fmt.Sprintf("Elevations not the same size,%v %v %v %s",testcase,val1,val2,testcase.CreateTitle())
	}
	
	for i := range val1 {
		if val1[i] != val2[i] {
			return fmt.Sprintf("Elevation values not the same check: %d mine: %d %s",val2[i],val1[i],testcase.CreateTitle())
		}
	}
	val1,val2 = testcase.Cursor.Geometry,testcase.CheckCursor.Geometry
	if len(val1) == len(val2) {
		return fmt.Sprintf("Geometry not the same size,%s",testcase.CreateTitle())
	}
	
	for i := range val1 {
		if val1[i] != val2[i] {
			return fmt.Sprintf("Geometry values not the same check: %d mine: %d %s",val2[i],val1[i],testcase.CreateTitle())
		}
	}
	return ""
}

func TestCases(t *testing.T) {
	cases := create()
	for _,ca := range cases {
		val := FindError(ca)
		if len(val) > 0 {
			t.Errorf(val)
		}
	}
}




