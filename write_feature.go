package vt

import (
	"github.com/paulmach/go.geojson"
	"reflect"
	//"fmt"
)

// adding a geojson feature to a given layer
func (layer *LayerWrite) AddFeature(feature *geojson.Feature) {
	// creating total bytes that holds the bytes for a given layer
	var array1,array2,array3,array4,array5,array6,array7,array8,array9 []byte
	// refreshing cursor 
	layer.RefreshCursor()
	

	if feature.ID != nil {
		// do the id shit
		vv := reflect.ValueOf(feature.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
			array3 = []byte{8}
			array4 = EncodeVarint(uint64(vv.Int()))
		}
	}
	
	if len(feature.Properties) > 0 {
		// do the tag shit here 
		array5 = []byte{18} // the key val
		array6 = WritePackedUint32(layer.GetTags(feature.Properties))
	}
	if feature.Geometry != nil {
		// do the geometry type shit here
		var geomtype byte
		switch feature.Geometry.Type {
		case "Point","MultiPoint":
			geomtype = 1
		case "LineString","MultiLineString":
			geomtype = 2
		case "Polygon","MultiPolygon":
			geomtype = 3
		}
		array7 = []byte{24,geomtype}
	}
	if feature.Geometry != nil {
		switch feature.Geometry.Type {
		case "Point":
			array8 = []byte{34}
			layer.Cursor.MakePointFloat(feature.Geometry.Point)
			array9 = WritePackedUint32(layer.Cursor.Geometry)		
		case "LineString":
			array8 = []byte{34}
			layer.Cursor.MakeLineFloat(feature.Geometry.LineString)
			array9 = WritePackedUint32(layer.Cursor.Geometry)		
		case "Polygon":
			array8 = []byte{34}
			layer.Cursor.MakePolygonFloat(feature.Geometry.Polygon)
			array9 = WritePackedUint32(layer.Cursor.Geometry)		
		case "MultiPoint":
			array8 = []byte{34}
			layer.Cursor.MakeMultiPointFloat(feature.Geometry.MultiPoint)
			array9 = WritePackedUint32(layer.Cursor.Geometry)		
		case "MultiLineString":
			array8 = []byte{34}
			layer.Cursor.MakeMultiLineFloat(feature.Geometry.MultiLineString)
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPolygon":
			array8 = []byte{34}
			layer.Cursor.MakeMultiPolygonFloat(feature.Geometry.MultiPolygon)
			array9 = WritePackedUint32(layer.Cursor.Geometry)

		}
	}
	array1 = []byte{18}
	array2 = EncodeVarint(uint64(len(array3)+len(array4)+len(array5)+len(array6)+len(array7)+len(array8)+len(array9)))
	layer.Features = append(layer.Features,AppendAll(array1,array2,array3,array4,array5,array6,array7,array8,array9)...)
}
