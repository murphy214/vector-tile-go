package vt

import (
	"github.com/paulmach/go.geojson"
	"testing"
	m "github.com/murphy214/mercantile"
)


var polygon_s = `{"geometry": {"type": "Polygon", "coordinates": [[[-7.734374999999999, 25.799891182088334], [10.8984375, -34.016241889667015], [45.703125, 17.644022027872726], [-5.9765625, 26.43122806450644], [-7.734374999999999, 25.799891182088334]]]}, "type": "Feature", "properties": {}}`
var multipolygon_s = `{"type":"Feature","properties":{},"geometry":{"type":"MultiPolygon","coordinates":[[[[-71.71875,51.17934297928927],[-36.2109375,-49.15296965617039],[30.585937499999996,0.3515602939922709],[29.179687499999996,59.17592824927136],[-38.3203125,70.72897946208789],[-71.71875,51.17934297928927]]],[[[33.3984375,74.68325030051861],[75.234375,16.29905101458183],[76.2890625,64.77412531292873],[32.6953125,75.23066741281573],[33.3984375,74.68325030051861]]]]}}`
var linestring_s = `{"geometry": {"type": "LineString", "coordinates": [[10.8984375, 56.17002298293205], [16.5234375, -2.108898659243126], [59.4140625, 42.032974332441405], [61.17187499999999, 42.293564192170095]]}, "type": "Feature", "properties": {}}`	
var multilinestring_s = `{"geometry": {"type": "MultiLineString", "coordinates": [[[-48.1640625, 47.754097979680026], [-9.140625, 4.214943141390651], [15.468749999999998, -9.102096738726443]], [[10.8984375, 56.17002298293205], [16.5234375, -2.108898659243126], [59.4140625, 42.032974332441405], [61.17187499999999, 42.293564192170095]]]}, "type": "Feature", "properties": {}}`
var point_s = `{"geometry": {"type": "Point", "coordinates": [-48.1640625, 47.754097979680026]}, "type": "Feature", "properties": {}}`
var multipoint_s = `{"geometry": {"type": "MultiPoint", "coordinates": [[-48.1640625, 47.754097979680026], [-9.140625, 4.214943141390651]]}, "type": "Feature", "properties": {}}`



var polygon,err = geojson.UnmarshalFeature([]byte(polygon_s))
var multipolygon,_ = geojson.UnmarshalFeature([]byte(multipolygon_s))
var linestring,_ = geojson.UnmarshalFeature([]byte(linestring_s))
var multilinestring,_ = geojson.UnmarshalFeature([]byte(multilinestring_s))
var point,_ = geojson.UnmarshalFeature([]byte(point_s))
var multipoint,_ = geojson.UnmarshalFeature([]byte(multipoint_s))

var polygon_geometry = []uint32{0x9, 0xf50, 0xda0, 0x22, 0x1a8, 0x598, 0x318, 0x4cf, 0x499, 0xd7, 0x25, 0x10, 0xf}
var multipolygon_geometry = []uint32{0x9, 0x99e, 0xab0, 0x2a, 0x32a, 0xa58, 0x5f0, 0x50f, 0x1f, 0x687, 0x5ff, 0x277, 0x2f9, 0x3b8, 0xf, 0x9, 0x95a, 0x4e7, 0x22, 0x3b8, 0x8c0, 0x18, 0x627, 0x3df, 0x2c9, 0x10, 0x32, 0xf}
var linestring_geometry = []uint32{0x9, 0x10f8, 0x9f0, 0x1a, 0x80, 0x640, 0x3d0, 0x44f, 0x28, 0x7}
var multilinestring_geometry = []uint32{0x9, 0xbb8, 0xb28, 0x12, 0x378, 0x478, 0x230, 0x130, 0x9, 0x67, 0x6df, 0x2a, 0x80, 0x640, 0x3d0, 0x44f, 0x28, 0x7}
var point_geometry = []uint32{0x9, 0xbb8, 0xb28}
var multipoint_geometry = []uint32{0x11, 0xbb8, 0xb28, 0x378, 0x478}

func TestPolygonFloat(t *testing.T) {
	cur := NewCursorExtent(m.TileID{0,0,0},4096)
	cur.MakePolygonFloat(polygon.Geometry.Polygon)

	for i := range cur.Geometry {
		if cur.Geometry[i] != polygon_geometry[i] {
			t.Errorf("Polygon Geometry Error")
		}
	}
}

func TestMultiPolygonFloat(t *testing.T) {
	cur := NewCursorExtent(m.TileID{0,0,0},4096)
	cur.MakeMultiPolygonFloat(multipolygon.Geometry.MultiPolygon)

	for i := range cur.Geometry {
		if cur.Geometry[i] != multipolygon_geometry[i] {
			t.Errorf("multipolygon Geometry Error")
		}
	}
}

func TestLineStringFloat(t *testing.T) {
	cur := NewCursorExtent(m.TileID{0,0,0},4096)
	cur.MakeLineFloat(linestring.Geometry.LineString)

	for i := range cur.Geometry {
		if cur.Geometry[i] != linestring_geometry[i] {
			t.Errorf("linestring Geometry Error")
		}
	}
}

func TestMultiLineStringFloat(t *testing.T) {
	cur := NewCursorExtent(m.TileID{0,0,0},4096)
	cur.MakeMultiLineFloat(multilinestring.Geometry.MultiLineString)

	for i := range cur.Geometry {
		if cur.Geometry[i] != multilinestring_geometry[i] {
			t.Errorf("multilinestring Geometry Error")
		}
	}
}


func TestPointFloat(t *testing.T) {
	cur := NewCursorExtent(m.TileID{0,0,0},4096)
	cur.MakePointFloat(point.Geometry.Point)

	for i := range cur.Geometry {
		if cur.Geometry[i] != point_geometry[i] {
			t.Errorf("point Geometry Error")
		}
	}
}


func TestMultiPointFloat(t *testing.T) {
	cur := NewCursorExtent(m.TileID{0,0,0},4096)
	cur.MakeMultiPointFloat(multipoint.Geometry.MultiPoint)

	for i := range cur.Geometry {
		if cur.Geometry[i] != multipoint_geometry[i] {
			t.Errorf("point Geometry Error")
		}
	}
}


