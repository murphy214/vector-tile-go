package vt

import (
	"github.com/paulmach/go.geojson"
	"reflect"
	//"fmt"
	//"github.com/murphy214/geobuf/geobuf_raw"

	p "github.com/murphy214/pbf"
)

// adding a geojson feature to a given layer
func (layer *LayerWrite) AddFeature(feature *geojson.Feature) {
	// creating total bytes that holds the bytes for a given layer
	var array1, array2, array3, array4, array5, array6, array7, array8, array9 []byte
	// refreshing cursor
	layer.RefreshCursor()

	if feature.ID != nil {
		// do the id shit
		vv := reflect.ValueOf(feature.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			array3 = []byte{8}
			array4 = p.EncodeVarint(uint64(vv.Int()))
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
		case "Point", "MultiPoint":
			geomtype = 1
		case "LineString", "MultiLineString":
			geomtype = 2
		case "Polygon", "MultiPolygon":
			geomtype = 3
		}
		array7 = []byte{24, geomtype}
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
	array2 = p.EncodeVarint(uint64(len(array3) + len(array4) + len(array5) + len(array6) + len(array7) + len(array8) + len(array9)))
	layer.Features = append(layer.Features, AppendAll(array1, array2, array3, array4, array5, array6, array7, array8, array9)...)
}

// adding a geobuf byte array to a given layer
// this function house's both the ingestion and output to vector tiles
// hopefully to reduce allocations
func (layer *LayerWrite) AddFeatureGeobuf(bytevals []byte) {
	boolval := false
	// the pbf representing a feauture
	pbf := p.PBF{Pbf: bytevals, Length: len(bytevals)}

	// creating total bytes that holds the bytes for a given layer
	var array1, array2, array3, array4, array5, array6, array7, array8, array9 []byte
	// refreshing cursor
	layer.RefreshCursor()

	key, val := pbf.ReadKey()

	if key == 1 && val == 0 {
		array3 = []byte{8}
		startpos := pbf.Pos
		pbf.ReadVarint()
		array4 = pbf.Pbf[startpos:pbf.Pos]
	}
	tags := []uint32{}
	for key == 2 && val == 2 {
		// starting properties shit here

		size := pbf.ReadVarint()
		endpos := pbf.Pos + size
		pbf.Pos += 1
		keyvalue := pbf.ReadString()
		keytag, keybool := layer.Keys_Map[keyvalue]
		if keybool == false {
			keytag = layer.AddKey(keyvalue)
		}
		tags = append(tags, keytag)

		pbf.Pos += 1
		pbf.ReadVarint()
		newkey, _ := pbf.ReadKey()
		var value interface{}
		switch newkey {
		case 1:
			value = pbf.ReadString()
		case 2:
			value = pbf.ReadFloat()
		case 3:
			value = pbf.ReadDouble()
		case 4:
			value = pbf.ReadInt64()
		case 5:
			value = pbf.ReadUInt64()
		case 6:
			value = pbf.ReadUInt64()
		case 7:
			value = pbf.ReadBool()
		}
		valuetag, valuebool := layer.Values_Map[value]
		if valuebool == false {
			valuetag = layer.AddValue(value)
		}
		tags = append(tags, valuetag)

		pbf.Pos = endpos
		key, val = pbf.ReadKey()
	}

	array5 = []byte{18} // the key val
	array6 = WritePackedUint32(tags)
	var geomtype string
	if key == 3 && val == 0 {
		switch int(pbf.Pbf[pbf.Pos]) {
		case 1:
			geomtype = "Point"
		case 2:
			geomtype = "LineString"
		case 3:
			geomtype = "Polygon"
		case 4:
			geomtype = "MultiPoint"
		case 5:
			geomtype = "MultiLineString"
		case 6:
			geomtype = "MultiPolygon"
		}
		pbf.Pos += 1
		key, val = pbf.ReadKey()
	}
	if len(geomtype) > 0 {
		// do the geometry type shit here
		var geomtypeb byte
		switch geomtype {
		case "Point", "MultiPoint":
			geomtypeb = 1
		case "LineString", "MultiLineString":
			geomtypeb = 2
		case "Polygon", "MultiPolygon":
			geomtypeb = 3
		}
		array7 = []byte{24, geomtypeb}
	}
	if key == 4 && val == 2 {
		boolval = true
		size := pbf.ReadVarint()
		endpos := pbf.Pos + size

		switch geomtype {
		case "Point":
			array8 = []byte{34}
			layer.Cursor.MakePointFloat(pbf.ReadPoint(endpos))
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "LineString":
			array8 = []byte{34}
			layer.Cursor.MakeLineFloat(pbf.ReadLine(0, endpos))
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "Polygon":
			array8 = []byte{34}
			layer.Cursor.MakePolygonFloat(pbf.ReadPolygon(endpos))
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPoint":
			array8 = []byte{34}
			layer.Cursor.MakeMultiPointFloat(pbf.ReadLine(0, endpos))
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "MultiLineString":
			array8 = []byte{34}
			layer.Cursor.MakeMultiLineFloat(pbf.ReadPolygon(endpos))
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPolygon":
			array8 = []byte{34}
			layer.Cursor.MakeMultiPolygonFloat(pbf.ReadMultiPolygon(endpos))
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		}
		key, val = pbf.ReadKey()
	}
	if boolval {
		array1 = []byte{18}
		array2 = p.EncodeVarint(uint64(len(array3) + len(array4) + len(array5) + len(array6) + len(array7) + len(array8) + len(array9)))
		layer.Features = append(layer.Features, AppendAll(array1, array2, array3, array4, array5, array6, array7, array8, array9)...)
	}
}
