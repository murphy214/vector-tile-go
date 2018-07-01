# vector-tile-go
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/murphy214/vector-tile-go)

An implementation of mapbox's vector-tile spec for reading / writing vector tiles. This is implementation is protocol buffer less meaning the serialization and deserialization is handled by me, which saves some needless allocations as well as pretty signicant performance gains.

### Why?

An implementation of mapbox's vector tile spec from with no protobuf file needed. Designed to reduce allocations and to be faster for reading / writing. When you think about how a vector tile is marshaled / unmarshaled from a regular protobuf implementation its kind of ridiculous, we transform (generally) geojson data into another entire feature struct with an entire new allocation for each, only to serialize that data structure to bytes immediately after. This repository implements API's for functionality that I find myself using for vector tiles.

### Installing

To start using vector-tile-go, install Go and run go get:

```
go get -u github.com/murphy214/vector-tile-go
```

### Features 

* Reads vector tiles into layer map geojson with a tile id given

* Can lazily read each feature within a layer so maps can be peformed based on only wanting certain properties and geometry types 

* Writes vector tile layers that can be appended to an existing vector tile byte array (this is just as useful than just supporting the full spec my api just handles one layer at a time)

* Writers vector tile layers for geojson-vt feature sets (formatted quite differently)

* Writes vector tile layers for geobuf data sets

* **Reads to geojson are 40% faster than a regular proto implementation**

* **Writes from geojson are 130% faster than regular proto implmentation**

* **Lazy reads are 70% faster than the protobuf implementation (i.e. marshaling the tile)**

### Caveats 

The api is currently still fluid and subject to change but currently everything I throw at works which I find to be pretty good, as I'm implementing several processes in one library. 

### Usage 

This repositories usage of reading and writing vector-tiles may seem a bit obtuse compared to general proto manner for reading protocol buffers. This is because for reading no vector-tile structure is every read as-is its either read completely into a geojson given the tile context OR read lazily feature / layer using each feature as needed. This provides all the structures that you would need to access anyway without having to carry around 3 representive data structures: vector tile feature, vector tile feature geometry read (integer format),vector tile feature in geojson format. 

#### Usage - Reading a Vector Tile as Geojson Features 

The example below simply reads in a vector tile as a slice of geojson features.

```golang
package main

import (
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/vector-tile-go"
)

func main() {
	// setting the byte values associated with the tile
	bytevals := []byte{0x1a, 0xc3, 0x2, 0xa, 0x4, 0x54, 0x65, 0x73, 0x74, 0x12, 0x3c, 0x12, 0x1a, 0x0, 0x0, 0x1, 0x1, 0x2, 0x2, 0x3, 0x3, 0x4, 0x4, 0x5, 0x5, 0x6, 0x6, 0x7, 0x7, 0x8, 0x8, 0x9, 0x9, 0xa, 0xa, 0xb, 0xa, 0xc, 0xb, 0x18, 0x2, 0x22, 0x1c, 0x9, 0x80, 0x41, 0xde, 0x3, 0x42, 0x75, 0x8d, 0x1, 0xab, 0x1, 0x71, 0x5d, 0x5b, 0xa9, 0x1, 0x83, 0x1, 0x8f, 0x1, 0x57, 0xdb, 0x2, 0x69, 0x43, 0x21, 0x1d, 0x19, 0x1a, 0x7, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x49, 0x44, 0x1a, 0x8, 0x53, 0x75, 0x62, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x1a, 0xa, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x1a, 0x8, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x1a, 0x6, 0x4f, 0x4e, 0x45, 0x57, 0x41, 0x59, 0x1a, 0x5, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x1a, 0x5, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x1a, 0xa, 0x53, 0x69, 0x67, 0x6e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x1a, 0x3, 0x45, 0x4d, 0x50, 0x1a, 0xa, 0x53, 0x68, 0x61, 0x70, 0x65, 0x5f, 0x4c, 0x65, 0x6e, 0x67, 0x1a, 0x3, 0x42, 0x4d, 0x50, 0x1a, 0x8, 0x53, 0x75, 0x70, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x1a, 0xa, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xf, 0xa, 0xd, 0x35, 0x30, 0x34, 0x30, 0x30, 0x35, 0x32, 0x36, 0x35, 0x30, 0x30, 0x30, 0x30, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x50, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x49, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x22, 0x2, 0x38, 0x0, 0x22, 0x7, 0xa, 0x5, 0x35, 0x32, 0x2f, 0x36, 0x35, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4a, 0x40, 0x22, 0x3, 0xa, 0x1, 0x34, 0x22, 0x9, 0x19, 0x40, 0x96, 0x4f, 0xa0, 0x99, 0x99, 0xd, 0x40, 0x22, 0x9, 0x19, 0x91, 0x80, 0xf2, 0xf3, 0x68, 0xbb, 0xb6, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x22, 0x14, 0xa, 0x12, 0x42, 0x49, 0x47, 0x20, 0x53, 0x41, 0x4e, 0x44, 0x59, 0x20, 0x52, 0x49, 0x56, 0x45, 0x52, 0x20, 0x52, 0x44, 0x78, 0x2}
	xyz := m.TileID{1107, 1578, 12} // the tileid these bytes relate to

	// reading in every feature (within every layer)
	// as can be seen below the layer field is just added as a property to each feature
	features := vt.ReadTile(bytevals, xyz)

	for _, feature := range features {
		fmt.Printf("LayerName: %s\nGeometry: %+v\nProperties: %+v\n\n", feature.Properties["layer"], feature.Geometry, feature.Properties)
	}
}

```

#### Usage - Reading a Vector Tile Lazily (As Needed)

The example below reads the vector tile by each feature, as needed parsing, the lazyfeature into what generally is a vector tile feature structure (i.e. properties & flat uint32 array for geometry), then, if desired, the geometry can be read as integers or converted entirely to a geojson feature given the tile it references.   

```golang
package main

import (
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/vector-tile-go"
)

func main() {
	// setting the byte values associated with the tile
	bytevals := []byte{0x1a, 0xc3, 0x2, 0xa, 0x4, 0x54, 0x65, 0x73, 0x74, 0x12, 0x3c, 0x12, 0x1a, 0x0, 0x0, 0x1, 0x1, 0x2, 0x2, 0x3, 0x3, 0x4, 0x4, 0x5, 0x5, 0x6, 0x6, 0x7, 0x7, 0x8, 0x8, 0x9, 0x9, 0xa, 0xa, 0xb, 0xa, 0xc, 0xb, 0x18, 0x2, 0x22, 0x1c, 0x9, 0x80, 0x41, 0xde, 0x3, 0x42, 0x75, 0x8d, 0x1, 0xab, 0x1, 0x71, 0x5d, 0x5b, 0xa9, 0x1, 0x83, 0x1, 0x8f, 0x1, 0x57, 0xdb, 0x2, 0x69, 0x43, 0x21, 0x1d, 0x19, 0x1a, 0x7, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x49, 0x44, 0x1a, 0x8, 0x53, 0x75, 0x62, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x1a, 0xa, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x1a, 0x8, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x1a, 0x6, 0x4f, 0x4e, 0x45, 0x57, 0x41, 0x59, 0x1a, 0x5, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x1a, 0x5, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x1a, 0xa, 0x53, 0x69, 0x67, 0x6e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x1a, 0x3, 0x45, 0x4d, 0x50, 0x1a, 0xa, 0x53, 0x68, 0x61, 0x70, 0x65, 0x5f, 0x4c, 0x65, 0x6e, 0x67, 0x1a, 0x3, 0x42, 0x4d, 0x50, 0x1a, 0x8, 0x53, 0x75, 0x70, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x1a, 0xa, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xf, 0xa, 0xd, 0x35, 0x30, 0x34, 0x30, 0x30, 0x35, 0x32, 0x36, 0x35, 0x30, 0x30, 0x30, 0x30, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x50, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x49, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x22, 0x2, 0x38, 0x0, 0x22, 0x7, 0xa, 0x5, 0x35, 0x32, 0x2f, 0x36, 0x35, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4a, 0x40, 0x22, 0x3, 0xa, 0x1, 0x34, 0x22, 0x9, 0x19, 0x40, 0x96, 0x4f, 0xa0, 0x99, 0x99, 0xd, 0x40, 0x22, 0x9, 0x19, 0x91, 0x80, 0xf2, 0xf3, 0x68, 0xbb, 0xb6, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x22, 0x14, 0xa, 0x12, 0x42, 0x49, 0x47, 0x20, 0x53, 0x41, 0x4e, 0x44, 0x59, 0x20, 0x52, 0x49, 0x56, 0x45, 0x52, 0x20, 0x52, 0x44, 0x78, 0x2}
	xyz := m.TileID{1107, 1578, 12} // the tileid these bytes relate to

	// loading the tile to perform lazy feature reads
	tile := vt.NewTile(bytevals)

	for layername, layer := range tile.LayerMap {
		fmt.Printf("LayerName: %s\n", layername)
		for layer.Next() {
			// getting the lazy feature
			lazyfeature := layer.Feature()
			fmt.Printf("Lazy Feature: %+v\n", lazyfeature)

			// loading the geometry as a geojson.Geometry structure
			// here the geometry is still in integer format
			// but stored a floats in the geojson.Geometry struct
			fmt.Printf("Geometry: %+v\n", lazyfeature.LoadGeometry())

			// loading the lazy feature as a geojson feature entirely
			fmt.Printf("Geojson Feature: %+v\n", lazyfeature.ToGeoJSON(xyz))

		}
	}

}
```

#### Usage - Compositing (Combining) Vector Tile Layers 

This repository intentially only implements a writer to write only layer at a time. This may seem like an odd behavior but in practice it utilizes one of the vector tile spec's most beneficial properties being that in order to combine disparate layers of the same tile, no structure in memory is needed just append one array of bytes to another and you have a vector tile byte array containing both layers.

This process MapBox calls compositing and allows you to only have to construct one layer at a time, which makes things much easier for generating multiple layer vector tiles. The example below shows compositing of two layers "layer1" & "layer2" respectively.

```golang
package main

import (
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/vector-tile-go"
	"github.com/paulmach/go.geojson"
)

func main() {

	// the features that will be written to both layers
	// these of course would be different for different layers normally
	features := []*geojson.Feature{&geojson.Feature{ID: interface{}(nil), Type: "", BoundingBox: []float64(nil), Geometry: (*geojson.Geometry)(&geojson.Geometry{Type: "LineString", BoundingBox: []float64(nil), Point: []float64(nil), MultiPoint: [][]float64(nil), LineString: [][]float64{[]float64{-82.61581420898438, 38.1305226701526}, []float64{-82.6170802116394, 38.13172105071115}, []float64{-82.61892557144165, 38.132683116639384}, []float64{-82.61993408203125, 38.13345951147531}, []float64{-82.62175798416138, 38.1345734548573}, []float64{-82.62330293655396, 38.135316074332906}, []float64{-82.62703657150269, 38.136210583213455}, []float64{-82.62776613235474, 38.13649749883382}, []float64{-82.62808799743652, 38.13671690413537}}, MultiLineString: [][][]float64(nil), Polygon: [][][]float64(nil), MultiPolygon: [][][][]float64(nil), Geometries: []*geojson.Geometry(nil), CRS: map[string]interface{}(nil)}), Properties: map[string]interface{}{"ROUTEID": "5040052650000", "SubRoute": 65, "Route": 52, "BMP": 0, "SuppCode": 0, "StreetName": "BIG SANDY RIVER RD", "CountyCode": 50, "District": 2, "ONEWAY": false, "Label": "52/65", "SignSystem": "4", "EMP": 3.70000005, "Shape_Leng": 5819.40997234}, CRS: map[string]interface{}(nil)}}

	xyz := m.TileID{1107, 1578, 12} // the tileid these bytes relate to

	// creating a 2 layers being "layer1" & "layer2" respectively
	config1 := vt.NewConfig("layer1", xyz)
	layer1bytes := vt.WriteLayer(features, config1)
	config2 := vt.NewConfig("layer2", xyz)
	layer2bytes := vt.WriteLayer(features, config2)

	// combining layers
	combined_layers_bytes := append(layer1bytes, layer2bytes...)

	// reading vector tile
	tile := vt.NewTile(combined_layers_bytes)

	// printing each layer in the tile
	for layername := range tile.LayerMap {
		fmt.Println(layername)
	}

}
```

#### Usage - Writing A Single Vector Tile Layer From GeoJSON Features

This example writes a single vector tile layer from geojson features from a given tile. 

```golang
package main

import (
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/vector-tile-go"
	"github.com/paulmach/go.geojson"
)

func main() {

	// the features that will be written to both layers
	// these of course would be different for different layers normally
	feature := &geojson.Feature{ID: interface{}(nil), Type: "", BoundingBox: []float64(nil), Geometry: (*geojson.Geometry)(&geojson.Geometry{Type: "LineString", BoundingBox: []float64(nil), Point: []float64(nil), MultiPoint: [][]float64(nil), LineString: [][]float64{[]float64{-82.61581420898438, 38.1305226701526}, []float64{-82.6170802116394, 38.13172105071115}, []float64{-82.61892557144165, 38.132683116639384}, []float64{-82.61993408203125, 38.13345951147531}, []float64{-82.62175798416138, 38.1345734548573}, []float64{-82.62330293655396, 38.135316074332906}, []float64{-82.62703657150269, 38.136210583213455}, []float64{-82.62776613235474, 38.13649749883382}, []float64{-82.62808799743652, 38.13671690413537}}, MultiLineString: [][][]float64(nil), Polygon: [][][]float64(nil), MultiPolygon: [][][][]float64(nil), Geometries: []*geojson.Geometry(nil), CRS: map[string]interface{}(nil)}), Properties: map[string]interface{}{"ROUTEID": "5040052650000", "SubRoute": 65, "Route": 52, "BMP": 0, "SuppCode": 0, "StreetName": "BIG SANDY RIVER RD", "CountyCode": 50, "District": 2, "ONEWAY": false, "Label": "52/65", "SignSystem": "4", "EMP": 3.70000005, "Shape_Leng": 5819.40997234}, CRS: map[string]interface{}(nil)}
	features := []*geojson.Feature{feature}

	xyz := m.TileID{1107, 1578, 12} // the tileid these bytes relate to

	// creating a 2 layers being "layer1" & "layer2" respectively
	config := vt.NewConfig("mylayer", xyz)
	layerbytes := vt.WriteLayer(features, config)

	// reading the feature I just wrote!
	fmt.Println(vt.ReadTile(layerbytes, xyz)[0])
}
```




