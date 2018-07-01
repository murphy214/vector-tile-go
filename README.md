# vector-tile-go
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/murphy214/vector-tile-go)

An implementation of mapbox's vector-tile spec for reading / writing vector tiles. This is implementation is protocol buffer less meaning the serialization and deserialization is handled by me, which saves some needless allocations as well as pretty signicant performance gains.

#### Why?

An implementation of mapbox's vector tile spec from with no protobuf file needed. Designed to reduce allocations and to be faster for reading / writing. When you think about how a vector tile is marshaled / unmarshaled from a regular protobuf implementation its kind of ridiculous, we transform (generally) geojson data into another entire feature struct with an entire new allocation for each, only to serialize that data structure to bytes immediately after. Thats really heavy, and unneeded allocations.

#### Features 

* Reads vector tiles into layer map geojson with a tile id given

* Can lazily read each feature within a layer so maps can be peformed based on only wanting certain properties and geometry types 

* Writes vector tile layers that can be appended to an existing vector tile byte array (this is just as useful than just supporting the full spec my api just handles one layer at a time)

* Writers vector tile layers for geojson-vt feature sets (formatted quite differently)

* Writes vector tile layers for geobuf data sets

* **Reads to geojson are 40% faster than a regular proto implementation**

* **Writes from geojson are 130% faster than regular proto implmentation**

* **Lazy reads are 70% faster than the protobuf implementation (i.e. marshaling the tile)**

#### Caveats 

The api is currently still fluid and subject to change but currently everything I throw at works which I find to be pretty good, as I'm implementing several processes in one library. 

#### Usage 

This repositories usage of reading and writing vector-tiles may seem a bit obtuse compared to general proto manner for reading protocol buffers. This is because for reading no vector-tile structure is every read as-is its either read completely into a geojson given the tile context OR read lazily feature / layer using each feature as needed. This provides all the structures that you would need to access anyway without having to carry around 3 representive data structures: vector tile feature, vector tile feature geometry read (integer format),vector tile feature in geojson format. 

##### Usage Reading a vector-tile as Geojson Features 

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





