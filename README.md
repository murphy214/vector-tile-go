# vector-tile-go

An implementation of mapbox's vector-tile spec for reading / writing vector tiles lazily.

# Introduction

An implementation of mapbox's vector tile spec from with no protobuf file needed. Designed to reduce allocations and to be faster for reading / writing. 

# Features 

* Reads vector tiles into layer map geojson with a tile id given

* Can lazily read each feature within a layer so maps can be peformed based on only wanting certain properties and geometry types 

* Writes vector tile layers that can be appended to an existing vector tile byte array (this is just as useful than just supporting the full spec my api just handles one layer at a time)

* **Reads to geojson are 40% faster than a regular proto implementation**

* **Writes from geojson are 130% faster than regular proto implmentation**

* **Lazy reads are 70% faster than the protobuf implementation (i.e. marshaling the tile)**

# Usage 

There 3 main apis that are used to read / write vector tiles. 

### Read To Geojson 

This shows how to read a vector tile byte array to a geojson layer map ```map[string][]*geojson.Feature```.

```golang
package main 

import (
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/vector-tile-go"
)

func main() {
	// a byte array representign vt data
	bytevals := []byte{0x1a, 0xa4, 0x2, 0xa, 0x7, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x19, 0x18, 0x3, 0x22, 0x15, 0x9, 0xd0, 0x1e, 0xa0, 0x1b, 0x22, 0xa8, 0x3, 0x98, 0xb, 0x98, 0x6, 0xcf, 0x9, 0x99, 0x9, 0xd7, 0x1, 0x25, 0x10, 0xf, 0x12, 0x32, 0x18, 0x3, 0x22, 0x2e, 0x9, 0x9e, 0x13, 0xb0, 0x15, 0x2a, 0xaa, 0x6, 0xd8, 0x14, 0xf0, 0xb, 0x8f, 0xa, 0x1f, 0x87, 0xd, 0xff, 0xb, 0xf7, 0x4, 0xf9, 0x5, 0xb8, 0x7, 0xf, 0x9, 0xda, 0x12, 0xe7, 0x9, 0x22, 0xb8, 0x7, 0xc0, 0x11, 0x18, 0xa7, 0xc, 0xdf, 0x7, 0xc9, 0x5, 0x10, 0x32, 0xf, 0x12, 0x14, 0x18, 0x2, 0x22, 0x10, 0x9, 0xf8, 0x21, 0xf0, 0x13, 0x1a, 0x80, 0x1, 0xc0, 0xc, 0xd0, 0x7, 0xcf, 0x8, 0x28, 0x7, 0x12, 0x21, 0x18, 0x2, 0x22, 0x1d, 0x9, 0xb8, 0x17, 0xa8, 0x16, 0x12, 0xf8, 0x6, 0xf8, 0x8, 0xb0, 0x4, 0xb0, 0x2, 0x9, 0x67, 0xdf, 0xd, 0x2a, 0x80, 0x1, 0xc0, 0xc, 0xd0, 0x7, 0xcf, 0x8, 0x28, 0x7, 0x12, 0x9, 0x18, 0x1, 0x22, 0x5, 0x9, 0xb8, 0x17, 0xa8, 0x16, 0x12, 0xd, 0x18, 0x1, 0x22, 0x9, 0x11, 0xb8, 0x17, 0xa8, 0x16, 0xf8, 0x6, 0xf8, 0x8, 0x12, 0x17, 0x12, 0xc, 0x0, 0x0, 0x1, 0x0, 0x2, 0x1, 0x3, 0x2, 0x4, 0x1, 0x5, 0x3, 0x18, 0x1, 0x22, 0x5, 0x9, 0xf0, 0x28, 0xda, 0x15, 0x1a, 0x9, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x6b, 0x65, 0x79, 0x1a, 0x8, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x6, 0x69, 0x6e, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x9, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x6b, 0x65, 0x79, 0x1a, 0x7, 0x75, 0x69, 0x6e, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x7, 0x62, 0x6f, 0x6f, 0x6c, 0x6b, 0x65, 0x79, 0x22, 0x9, 0x19, 0x1f, 0x85, 0xeb, 0x51, 0xb8, 0xe, 0x59, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x40, 0x26, 0x50, 0x0, 0x3, 0x42, 0x22, 0x8, 0xa, 0x6, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x2, 0x38, 0x1, 0x78, 0x2}

	// setting are tileid
	tileid := m.TileID{0,0,0}

	layermap := vt.ReadTile(bytevals,tileid)

	for layername,features := range layermap {
		fmt.Println(layername,features)
	}
}
```

### Lazy Reads 

This code example shows you the lazy reading api and the methods used to call them.
```golang
package main 

import (
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/vector-tile-go"
)

func main() {
	// a byte array representign vt data
	bytevals := []byte{0x1a, 0xa4, 0x2, 0xa, 0x7, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x19, 0x18, 0x3, 0x22, 0x15, 0x9, 0xd0, 0x1e, 0xa0, 0x1b, 0x22, 0xa8, 0x3, 0x98, 0xb, 0x98, 0x6, 0xcf, 0x9, 0x99, 0x9, 0xd7, 0x1, 0x25, 0x10, 0xf, 0x12, 0x32, 0x18, 0x3, 0x22, 0x2e, 0x9, 0x9e, 0x13, 0xb0, 0x15, 0x2a, 0xaa, 0x6, 0xd8, 0x14, 0xf0, 0xb, 0x8f, 0xa, 0x1f, 0x87, 0xd, 0xff, 0xb, 0xf7, 0x4, 0xf9, 0x5, 0xb8, 0x7, 0xf, 0x9, 0xda, 0x12, 0xe7, 0x9, 0x22, 0xb8, 0x7, 0xc0, 0x11, 0x18, 0xa7, 0xc, 0xdf, 0x7, 0xc9, 0x5, 0x10, 0x32, 0xf, 0x12, 0x14, 0x18, 0x2, 0x22, 0x10, 0x9, 0xf8, 0x21, 0xf0, 0x13, 0x1a, 0x80, 0x1, 0xc0, 0xc, 0xd0, 0x7, 0xcf, 0x8, 0x28, 0x7, 0x12, 0x21, 0x18, 0x2, 0x22, 0x1d, 0x9, 0xb8, 0x17, 0xa8, 0x16, 0x12, 0xf8, 0x6, 0xf8, 0x8, 0xb0, 0x4, 0xb0, 0x2, 0x9, 0x67, 0xdf, 0xd, 0x2a, 0x80, 0x1, 0xc0, 0xc, 0xd0, 0x7, 0xcf, 0x8, 0x28, 0x7, 0x12, 0x9, 0x18, 0x1, 0x22, 0x5, 0x9, 0xb8, 0x17, 0xa8, 0x16, 0x12, 0xd, 0x18, 0x1, 0x22, 0x9, 0x11, 0xb8, 0x17, 0xa8, 0x16, 0xf8, 0x6, 0xf8, 0x8, 0x12, 0x17, 0x12, 0xc, 0x0, 0x0, 0x1, 0x0, 0x2, 0x1, 0x3, 0x2, 0x4, 0x1, 0x5, 0x3, 0x18, 0x1, 0x22, 0x5, 0x9, 0xf0, 0x28, 0xda, 0x15, 0x1a, 0x9, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x6b, 0x65, 0x79, 0x1a, 0x8, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x6, 0x69, 0x6e, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x9, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x6b, 0x65, 0x79, 0x1a, 0x7, 0x75, 0x69, 0x6e, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x7, 0x62, 0x6f, 0x6f, 0x6c, 0x6b, 0x65, 0x79, 0x22, 0x9, 0x19, 0x1f, 0x85, 0xeb, 0x51, 0xb8, 0xe, 0x59, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x40, 0x26, 0x50, 0x0, 0x3, 0x42, 0x22, 0x8, 0xa, 0x6, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x2, 0x38, 0x1, 0x78, 0x2}

	// setting are tileid
	tileid := m.TileID{0,0,0}

	tile := vt.NewTile(bytevals)	

	for layername,layer := range tile.LayerMap {
		fmt.Println(layername)
		// iterating through each feature in a layer
		for layer.Next() {
			// the lazy feature 
			lazy_feature := layer.Feature()
			fmt.Printf("lazy feature: %v\n",lazy_feature)

			// loading the geometry if needed 
			fmt.Printf("geometry: %v\n",lazy_feature.LoadGeometry())

			// finally loading the entire geojson feature 
			fmt.Printf("geojson feature: %v\n",lazy_feature.ToGeoJSON(tileid))
			fmt.Println("\n\n\n")
		}
	}
}
```

### Writes from Geojson

This example shows how to write layers of vector tiles from an array of geojson features. **It also shows how only being able to write one layer at a time really isn't an issue.**
```golang
package main

import (
	"github.com/paulmach/go.geojson"
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/vector-tile-go"
)

func main() {
	fc_s := `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-7.734375,25.799891182088302],[10.8984375,-34.016241889667015],[45.703125,17.644022027872722],[-6.064453125,26.431228064506442],[-7.734375,25.799891182088302]]]},"properties":{}},{"type":"Feature","geometry":{"type":"MultiPolygon","coordinates":[[[[-71.806640625,51.17934297928926],[-36.2109375,-49.1529696561704],[30.5859375,0.3515602939922644],[29.1796875,59.17592824927135],[-38.3203125,70.72897946208789],[-71.806640625,51.17934297928926]]],[[[33.3984375,74.68325030051861],[75.234375,16.299051014581835],[76.2890625,64.7741253129287],[32.6953125,75.25305660483545],[33.3984375,74.68325030051861]]]]},"properties":{}},{"type":"Feature","geometry":{"type":"LineString","coordinates":[[10.8984375,56.17002298293204],[16.5234375,-2.108898659243124],[59.4140625,42.03297433244137],[61.171875,42.293564192170095]]},"properties":{}},{"type":"Feature","geometry":{"type":"MultiLineString","coordinates":[[[-48.1640625,47.75409797968001],[-9.140625,4.214943141390634],[15.46875,-9.102096738726445]],[[10.8984375,56.17002298293204],[16.5234375,-2.108898659243124],[59.4140625,42.03297433244137],[61.171875,42.293564192170095]]]},"properties":{}},{"type":"Feature","geometry":{"type":"Point","coordinates":[-48.1640625,47.75409797968001]},"properties":{}},{"type":"Feature","geometry":{"type":"MultiPoint","coordinates":[[-48.1640625,47.75409797968001],[-9.140625,4.214943141390634]]},"properties":{}},{"type":"Feature","geometry":{"type":"Point","coordinates":[49.921875,50.007739014636854]},"properties":{"boolkey":true,"doublekey":100.23,"floatkey":100.23,"intkey":10201203912,"stringkey":"string","uintkey":10201203912}}]}`
	fc,_ := geojson.UnmarshalFeatureCollection([]byte(fc_s))


	// writing layer1
	config_layer1 := vt.NewConfig("layer1",m.TileID{0,0,0})
	bytevals_layer1 := vt.WriteLayer(fc.Features,config_layer1)	

	// writing layer2 	
	config_layer2 := vt.NewConfig("layer2",m.TileID{0,0,0})
	bytevals_layer2 := vt.WriteLayer(fc.Features,config_layer2)	

	// appending each byte array 
	total_bytes := append(bytevals_layer1,bytevals_layer2...)

	// reading the new vector tile to show how layers work
	fmt.Println(vt.ReadTile(total_bytes,m.TileID{0,0,0}))
}
```

#### Output:
```
map[layer1:[0xc42008c410 0xc42008c460 0xc42008c4b0 0xc42008c500 0xc42008c550 0xc42008c5a0 0xc42008c5f0] layer2:[0xc42008c640 0xc42008c690 0xc42008c6e0 0xc42008c730 0xc42008c780 0xc42008c7d0 0xc42008c820]]
```

