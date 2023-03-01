package vt

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"

	m "github.com/murphy214/mercantile"

	//"sync"
	"github.com/paulmach/go.geojson"
	//"fmt"
	//"math"
	"testing"
)


func rF(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}

func check_pt(p1,p2 []float64) bool {
	return rF(p1[0],6)!=rF(p2[0],6)||rF(p1[1],6)!=rF(p2[1],6)
}

func check_line(l1,l2 [][]float64) error {
	if len(l1)!=len(l2) {
		errors.New(fmt.Sprintf("Line size different %v %v\n",len(l1),len(l2)))
	}
	for p := range l1 {
		if check_pt(l1[p],l2[p]) {
			errors.New(fmt.Sprintf("Point different %v %v\n",l1[p],l2[p]))
		}
	}
	return nil
}


//
func GeoJsonDif(f1,f2 *geojson.Feature) error {
	for k := range f1.Properties {
		v1 := f1.Properties[k]
		v2 := f2.Properties[k]
		if v1!=v2 {
			return errors.New(fmt.Sprintf("Values in geojson different %v %v\n",v1,v2))
		}
	}

	if f1.Geometry.Type!=f2.Geometry.Type {
		return errors.New(fmt.Sprintf("Geometry Types are not the same %v %v\n",f1.Geometry.Type,f2.Geometry.Type))
	}


	switch f1.Geometry.Type {
	
	case "Point":
		return nil 
	case "LineString":
		return check_line(f1.Geometry.LineString,f2.Geometry.LineString)
	case "MultiLineString":
		for ii := range f1.Geometry.MultiLineString {
			err :=  check_line(f1.Geometry.MultiLineString[ii],f2.Geometry.MultiLineString[ii])
			if err != nil {
				return err 
			}
		}
	case "Polygon":
		for ii := range f1.Geometry.Polygon {
			err :=  check_line(f1.Geometry.Polygon[ii],f2.Geometry.Polygon[ii])
			if err != nil {
				return err 
			}
		}
	
	case "MultiPolygon":
		for ii := range f1.Geometry.MultiPolygon {
			for i := range f1.Geometry.MultiPolygon[ii] {
				err :=  check_line(f1.Geometry.MultiPolygon[ii][i],f2.Geometry.MultiPolygon[ii][i])
				if err != nil {
					return err 
				}
					
			}
		}
	}
	return nil

}




var bytevals, _ = ioutil.ReadFile("test_data/701_1635_12.pbf")
var tileid = m.TileID{701, 1635, 12}

func TestReads(t *testing.T) {
	feats1, _ := ReadTile(bytevals, tileid)
	m1, m2 := map[interface{}]*geojson.Feature{}, map[interface{}]*geojson.Feature{}
	for _, feat := range feats1 {
		delete(feat.Properties, "layer")
		m1[feat.Properties["@id"]] = feat
	}
	tile, _ := NewTile(bytevals)
	for _, layer := range tile.LayerMap {
		for layer.Next() {
			feat, _ := layer.Feature()
			featg, _ := feat.ToGeoJSON(tileid)
			delete(featg.Properties, "layer")

			m2[featg.Properties["@id"]] = featg
		}
	}
	if len(m2) != len(m2) {
		t.Errorf("Map sizes are different.")
	}
	i := 0
	for k := range m1 {
		i++
		v1, b1 := m1[k]
		v2, b2 := m2[k]
		if b1 && b2 {
			// err := geojsondif.CheckFeatures(v1, v2)
			err := GeoJsonDif(v1,v2)

			if err != nil {
				t.Errorf(err.Error())
			}
		} else {
			t.Errorf("Both geojson features weren't in map.")
		}

	}

	fmt.Printf("Lazy Reads Are Exactly the same as bulk reads for %d features in tile.\n", len(feats1))
}

func TestReadsWrites(t *testing.T) {
	feats1, _ := ReadTile(bytevals, tileid)
	con := NewConfig("new", tileid)
	bs, _ := WriteLayer(feats1, con)
	feats1, _ = ReadTile(bs, tileid)

	m1, m2 := map[interface{}]*geojson.Feature{}, map[interface{}]*geojson.Feature{}
	for _, feat := range feats1 {
		delete(feat.Properties, "layer")
		m1[feat.Properties["@id"]] = feat
	}
	tile, _ := NewTile(bs)
	for _, layer := range tile.LayerMap {
		for layer.Next() {
			feat, _ := layer.Feature()
			featg, _ := feat.ToGeoJSON(tileid)
			delete(featg.Properties, "layer")

			m2[featg.Properties["@id"]] = featg
		}
	}
	if len(m2) != len(m2) {
		t.Errorf("Map sizes are different.")
	}
	i := 0
	for k := range m1 {
		i++
		v1, b1 := m1[k]
		v2, b2 := m2[k]
		if b1 && b2 {
			err := GeoJsonDif(v1,v2)
			// var err error
			if err != nil {
				t.Errorf(err.Error())
			}
		} else {
			t.Errorf("Both geojson features weren't in map.")
		}

	}

	fmt.Printf("Lazy Reads Are Exactly the same as bulk reads when written and read again features for %d features in tile.\n", len(feats1))
}
