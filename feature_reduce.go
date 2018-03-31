package vt

import (
	//"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
)

var default_steps = 8
var default_percentage = .005

type Reduce_Config struct {
	Point_Steps int     // should be a binary number
	Percent     float64 // percentage of occupied dimmension
	TileID      m.TileID
	TileMap     map[m.TileID]string // tilemap
	DeltaX      float64
	DeltaY      float64
	Zoom        int
	OperateBool bool
}

func NewReduceConfig(tileid m.TileID) *Reduce_Config {
	bds := m.Bounds(tileid)
	deltax := bds.E - bds.W
	deltay := bds.N - bds.S
	zoom := int(tileid.Z) + default_steps
	return &Reduce_Config{
		Point_Steps: default_steps,
		Percent:     default_percentage,
		TileID:      tileid,
		TileMap:     map[m.TileID]string{},
		DeltaX:      deltax,
		DeltaY:      deltay,
		Zoom:        zoom,
		OperateBool: true,
	}
}

// BoundingBox implementation as per https://tools.ietf.org/html/rfc7946
// BoundingBox syntax: "bbox": [west, south, east, north]
// BoundingBox defaults "bbox": [-180.0, -90.0, 180.0, 90.0]
func BoundingBox_Points(pts [][]float64) []float64 {
	// setting opposite default values
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	for _, pt := range pts {
		x, y := pt[0], pt[1]
		// can only be one condition
		// using else if reduces one comparison
		if x < west {
			west = x
		} else if x > east {
			east = x
		}

		if y < south {
			south = y
		} else if y > north {
			north = y
		}
	}
	return []float64{west, south, east, north}
}

func Push_Two_BoundingBoxs(bb1 []float64, bb2 []float64) []float64 {
	// setting opposite default values
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	// setting bb1 and bb2
	west1, south1, east1, north1 := bb1[0], bb1[1], bb1[2], bb1[3]
	west2, south2, east2, north2 := bb2[0], bb2[1], bb2[2], bb2[3]

	// handling west values: min
	if west1 < west2 {
		west = west1
	} else {
		west = west2
	}

	// handling south values: min
	if south1 < south2 {
		south = south1
	} else {
		south = south2
	}

	// handling east values: max
	if east1 > east2 {
		east = east1
	} else {
		east = east2
	}

	// handling north values: max
	if north1 > north2 {
		north = north1
	} else {
		north = north2
	}

	return []float64{west, south, east, north}
}

// this functions takes an array of bounding box objects and
// pushses them all out
func Expand_BoundingBoxs(bboxs [][]float64) []float64 {
	bbox := bboxs[0]
	for _, temp_bbox := range bboxs[1:] {
		bbox = Push_Two_BoundingBoxs(bbox, temp_bbox)
	}
	return bbox
}

// boudning box on a normal point geometry
// relatively useless
func BoundingBox_PointGeometry(pt []float64) []float64 {
	return []float64{pt[0], pt[1], pt[0], pt[1]}
}

// Returns BoundingBox for a MultiPoint
func BoundingBox_MultiPointGeometry(pts [][]float64) []float64 {
	return BoundingBox_Points(pts)
}

// Returns BoundingBox for a LineString
func BoundingBox_LineStringGeometry(line [][]float64) []float64 {
	return BoundingBox_Points(line)
}

// Returns BoundingBox for a MultiLineString
func BoundingBox_MultiLineStringGeometry(multiline [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, line := range multiline {
		bboxs = append(bboxs, BoundingBox_Points(line))
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns BoundingBox for a Polygon
func BoundingBox_PolygonGeometry(polygon [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, cont := range polygon {
		bboxs = append(bboxs, BoundingBox_Points(cont))
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns BoundingBox for a Polygon
func BoundingBox_MultiPolygonGeometry(multipolygon [][][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, polygon := range multipolygon {
		for _, cont := range polygon {
			bboxs = append(bboxs, BoundingBox_Points(cont))
		}
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns a BoundingBox for a geometry collection
func BoundingBox_GeometryCollection(gs []*geojson.Geometry) []float64 {
	bboxs := [][]float64{}
	for _, g := range gs {
		bboxs = append(bboxs, g.Get_BoundingBox())
	}
	return Expand_BoundingBoxs(bboxs)
}

// retrieves a boundingbox for a given geometry
func Get_BoundingBox(g *geojson.Geometry) []float64 {
	switch g.Type {
	case "Point":
		return BoundingBox_PointGeometry(g.Point)
	case "MultiPoint":
		return BoundingBox_MultiPointGeometry(g.MultiPoint)
	case "LineString":
		return BoundingBox_LineStringGeometry(g.LineString)
	case "MultiLineString":
		return BoundingBox_MultiLineStringGeometry(g.MultiLineString)
	case "Polygon":
		return BoundingBox_PolygonGeometry(g.Polygon)
	case "MultiPolygon":
		return BoundingBox_MultiPolygonGeometry(g.MultiPolygon)

	}
	return []float64{}
}

func Filter(feature *geojson.Feature, config *Reduce_Config) bool {
	if !config.OperateBool {
		return true
	}
	switch feature.Geometry.Type {
	case "Point":

		tile := m.Tile(feature.Geometry.Point[0], feature.Geometry.Point[1], config.Zoom)
		_, boolval := config.TileMap[tile]
		if !boolval {
			config.TileMap[tile] = ""
			return true
		} else {
			return false
		}
	case "MultiPoint":
		total_x, total_y := 0.0, 0.0
		for _, point := range feature.Geometry.MultiPoint {
			//tile := m.Tile(point[0], point[1], config.Zoom)
			total_x += point[0]
			total_y += point[1]
		}
		size := float64(len(feature.Geometry.MultiPoint))
		avg_x, avg_y := total_x/size, total_y/size
		tile := m.Tile(avg_x, avg_y, config.Zoom)
		_, boolval := config.TileMap[tile]
		if !boolval {
			config.TileMap[tile] = ""
			feature.Geometry.Type = "Point"
			feature.Geometry.MultiPoint = [][]float64{}
			feature.Geometry.Point = []float64{avg_x, avg_y}
			return true
		} else {
			return false
		}
	case "LineString", "MultiLineString", "MultiPolygon", "Polygon":
		//west, south, east, north
		//m.Extrema{W:feature.BoundingBox[0],S:feature.BoundingBox[1],EA}
		bbox := Get_BoundingBox(feature.Geometry)
		percentx := (bbox[2] - bbox[0]) / config.DeltaX
		percenty := (bbox[3] - bbox[1]) / config.DeltaY
		//fmt.Println(percentx, percenty)

		return percentx > config.Percent || percenty > config.Percent
	}
	return false
}
