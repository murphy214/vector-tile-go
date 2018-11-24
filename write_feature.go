package vt

import (
	"github.com/paulmach/go.geojson"
	"reflect"
	//"fmt"
	//"github.com/murphy214/geobuf/geobuf_raw"
	"github.com/murphy214/vector-tile-go/tags"
	p "github.com/murphy214/pbf"
)

// adding a geojson feature to a given layer
func (layer *LayerWrite) AddFeature(feature *geojson.Feature) {
	// creating total bytes that holds the bytes for a given layer
	var array1, array2, array3, array4, array5, array6, array7, array8, array9,array10,array11,array12,array13,array14,array15,array16,array17,array18,array19,array20 []byte
	totalsize := 0
	// refreshing cursor
	layer.RefreshCursor()
	layer.Cursor.Scaling = layer.ElevationScaling
	var stringbool bool
	var stringval string
	if feature.ID != nil {
		// do the id shit
		vv := reflect.ValueOf(feature.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			array3 = []byte{8}
			array4 = p.EncodeVarint(uint64(vv.Int()))
			totalsize+=(len(array3)+len(array4))
		case reflect.String:
			stringval = vv.String()
			stringbool = true
		}
	}

	// parsing out geometry attributes 
	// then setting the geometry attributes to the cursor object
	gattrs,boolval := feature.Properties["geometric_attributes"]
	layer.Cursor.GeometricAttributesBool = boolval
	if boolval {
		var gattrsmap map[string][]interface{}
		gattrsmap,boolval = gattrs.(map[string][]interface{})
		if boolval {
			layer.Cursor.SetCursorGeometricAttributes(gattrsmap)
			delete(feature.Properties,"geometric_attributes")
		}
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
		array5 = []byte{24, geomtype}
		totalsize+=len(array5)
	}

	var abort_bool bool
	if feature.Geometry != nil {
		switch feature.Geometry.Type {
		case "Point":
			array6 = []byte{34}
			layer.Cursor.MakePointFloat(feature.Geometry.Point)
			array7 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "LineString":
			array6 = []byte{34}
			layer.Cursor.MakeLineFloat(feature.Geometry.LineString)
			if layer.Cursor.Count == 0 {
				abort_bool = true
			}
			array7 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "Polygon":
			array6 = []byte{34}
			layer.Cursor.MakePolygonFloat(feature.Geometry.Polygon)
			array7 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPoint":
			array6 = []byte{34}
			layer.Cursor.MakeMultiPointFloat(feature.Geometry.MultiPoint)
			array7 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "MultiLineString":
			array6 = []byte{34}
			layer.Cursor.MakeMultiLineFloat(feature.Geometry.MultiLineString)
			array7 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPolygon":
			array6 = []byte{34}
			layer.Cursor.MakeMultiPolygonFloat(feature.Geometry.MultiPolygon)
			array7 = tags.WritePackedUint32(layer.Cursor.Geometry)
		}
		totalsize+=(len(array6)+len(array7))
	}

	// adding the normal attributes tags
	if len(feature.Properties) > 0 {
		// do the tag shit here
		array8 = []byte{42} // the key val
		array9 = tags.WritePackedInt(layer.TagWriter.MakeProperties(feature.Properties))
		totalsize+=(len(array8)+len(array9))
	}

	// adding the geometric attriburtes map attributes map 
	if len(layer.Cursor.NewGeometricAttributesMap) > 0 {
		array10 = []byte{50} // the key val
		array11 = tags.WritePackedInt(
			layer.TagWriter.MakeProperties(
				DumpInterfaceMap(layer.Cursor.NewGeometricAttributesMap),
			),
		)
		totalsize+=(len(array10)+len(array11))
	}

	// adding elevations
	if len(layer.Cursor.Elevations) > 0 {
		array12 = []byte{58}
		array13 = tags.WritePackedUint32(layer.Cursor.Elevations)
		totalsize+=(len(array12)+len(array13))
	}

	// adding spline knots
	if len(layer.Cursor.SplineKnots) > 0 {
		array14 = []byte{66}
		array15 = tags.WritePackedInt(layer.Cursor.SplineKnots)
		totalsize+=(len(array14)+len(array15))
	}

	// adding spline degree
	if layer.Cursor.SplineDegree > 0 {
		array16 = []byte{72}
		array17 = p.EncodeVarint(uint64(layer.Cursor.SplineDegree))
		totalsize+=(len(array16)+len(array17))
	}

	// writes a string id to a out to a feature
	if stringbool {
		array18 = []byte{82}
		array19 = p.EncodeVarint(uint64(len(stringval)))
		array20 = []byte(stringval)
		totalsize+=(len(array18)+len(array19)+len(array20))
	}



	// on the off chane one of my lines contains one point
	if !abort_bool {
		array1 = []byte{18}
		//array2 = p.EncodeVarint(uint64(len(array3) + len(array4) + len(array5) + len(array6) + len(array7) + len(array8) + len(array9)))
		array2 = p.EncodeVarint(uint64(totalsize))
		layer.Features = append(layer.Features, tags.AppendAll(
			array1,
			array2, 
			array3, 
			array4, 
			array5, 
			array6, 
			array7, 
			array8,
			array9,
			array10,
			array11, 
			array12, 
			array13, 
			array14, 
			array15, 
			array16, 
			array17,
			array18,
			array19,
			array20,
		)...)
	}
}

// function for adding the feature for a raw implementation
func (layer *LayerWrite) AddFeatureRaw(id int, geomtype int, geometry []uint32, properties map[string]interface{}) {
	var array1, array2, array3, array4, array5, array6, array7, array8, array9 []byte
	// refreshing cursor
	layer.RefreshCursor()

	if id > 0 {
		// do the id shit
		array3 = []byte{8}
		array4 = p.EncodeVarint(uint64(id))
	}

	if len(properties) > 0 {
		// do the tag shit here
		array5 = []byte{18} // the key val
		array6 = tags.WritePackedUint32(layer.GetTags(properties))
	}
	if geomtype != 0 {
		// do the geometry type shit here
		array7 = []byte{24, byte(geomtype)}
	}
	// adding geometry
	if len(geometry) > 0 {
		array8 = []byte{34}
		array9 = tags.WritePackedUint32(geometry)
	}

	// on the off chane one of my lines contains one point
	array1 = []byte{18}
	array2 = p.EncodeVarint(uint64(len(array3) + len(array4) + len(array5) + len(array6) + len(array7) + len(array8) + len(array9)))
	layer.Features = append(layer.Features, tags.AppendAll(array1, array2, array3, array4, array5, array6, array7, array8, array9)...)
}

// adding a geobuf byte array to a given layer
// this function house's both the ingestion and output to vector tiles
// hopefully to reduce allocations
func (layer *LayerWrite) AddFeatureGeobuf(bytevals []byte) {
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
	tagss := []uint32{}
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
		tagss = append(tagss, keytag)

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
		tagss = append(tagss, valuetag)

		pbf.Pos = endpos
		key, val = pbf.ReadKey()
	}

	array5 = []byte{18} // the key val
	array6 = tags.WritePackedUint32(tagss)
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
	var abort_bool bool
	if key == 4 && val == 2 {
		size := pbf.ReadVarint()
		endpos := pbf.Pos + size

		switch geomtype {
		case "Point":
			array8 = []byte{34}
			layer.Cursor.MakePointFloat(pbf.ReadPoint(endpos))
			array9 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "LineString":
			array8 = []byte{34}
			layer.Cursor.MakeLineFloat(pbf.ReadLine(0, endpos))
			if layer.Cursor.Count == 0 {
				abort_bool = true
			}
			array9 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "Polygon":
			array8 = []byte{34}
			layer.Cursor.MakePolygonFloat(pbf.ReadPolygon(endpos))
			array9 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPoint":
			array8 = []byte{34}
			layer.Cursor.MakeMultiPointFloat(pbf.ReadLine(0, endpos))
			array9 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "MultiLineString":
			array8 = []byte{34}
			layer.Cursor.MakeMultiLineFloat(pbf.ReadPolygon(endpos))
			array9 = tags.WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPolygon":
			array8 = []byte{34}
			layer.Cursor.MakeMultiPolygonFloat(pbf.ReadMultiPolygon(endpos))
			array9 = tags.WritePackedUint32(layer.Cursor.Geometry)
		}
		key, val = pbf.ReadKey()
	}
	if !abort_bool {
		array1 = []byte{18}
		array2 = p.EncodeVarint(uint64(len(array3) + len(array4) + len(array5) + len(array6) + len(array7) + len(array8) + len(array9)))
		layer.Features = append(layer.Features, tags.AppendAll(array1, array2, array3, array4, array5, array6, array7, array8, array9)...)
	}

}
