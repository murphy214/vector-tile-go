package vt

import (
	"github.com/golang/protobuf/proto"
	"github.com/murphy214/mbtiles-util/vector-tile/2.1"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"strings"
	"testing"
)

var fc_s = `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-7.734375,25.799891182088302],[10.8984375,-34.016241889667015],[45.703125,17.644022027872722],[-6.064453125,26.431228064506442],[-7.734375,25.799891182088302]]]},"properties":null},{"type":"Feature","geometry":{"type":"MultiPolygon","coordinates":[[[[-71.806640625,51.17934297928926],[-36.2109375,-49.1529696561704],[30.5859375,0.3515602939922644],[29.1796875,59.17592824927135],[-38.3203125,70.72897946208789],[-71.806640625,51.17934297928926]]],[[[33.3984375,74.68325030051861],[75.234375,16.299051014581835],[76.2890625,64.7741253129287],[32.6953125,75.25305660483545],[33.3984375,74.68325030051861]]]]},"properties":null},{"type":"Feature","geometry":{"type":"LineString","coordinates":[[10.8984375,56.17002298293204],[16.5234375,-2.108898659243124],[59.4140625,42.03297433244137],[61.171875,42.293564192170095]]},"properties":null},{"type":"Feature","geometry":{"type":"MultiLineString","coordinates":[[[-48.1640625,47.75409797968001],[-9.140625,4.214943141390634],[15.46875,-9.102096738726445]],[[10.8984375,56.17002298293204],[16.5234375,-2.108898659243124],[59.4140625,42.03297433244137],[61.171875,42.293564192170095]]]},"properties":null},{"type":"Feature","geometry":{"type":"Point","coordinates":[-48.1640625,47.75409797968001]},"properties":null},{"type":"Feature","geometry":{"type":"MultiPoint","coordinates":[[-48.1640625,47.75409797968001],[-9.140625,4.214943141390634]]},"properties":null},{"type":"Feature","geometry":{"type":"Point","coordinates":[49.921875,50.007739014636854]},"properties":{"boolkey":true,"doublekey":100.23,"floatkey":100.23,"intkey":10201203912,"stringkey":"string","uintkey":10201203912}}]}`
var fc, _ = geojson.UnmarshalFeatureCollection([]byte(fc_s))

var tile_expected = []byte{0x1a, 0xa4, 0x2, 0xa, 0x7, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x19, 0x18, 0x3, 0x22, 0x15, 0x9, 0xd0, 0x1e, 0xa0, 0x1b, 0x22, 0xa8, 0x3, 0x98, 0xb, 0x98, 0x6, 0xcf, 0x9, 0x99, 0x9, 0xd7, 0x1, 0x25, 0x10, 0xf, 0x12, 0x32, 0x18, 0x3, 0x22, 0x2e, 0x9, 0x9e, 0x13, 0xb0, 0x15, 0x2a, 0xaa, 0x6, 0xd8, 0x14, 0xf0, 0xb, 0x8f, 0xa, 0x1f, 0x87, 0xd, 0xff, 0xb, 0xf7, 0x4, 0xf9, 0x5, 0xb8, 0x7, 0xf, 0x9, 0xda, 0x12, 0xe7, 0x9, 0x22, 0xb8, 0x7, 0xc0, 0x11, 0x18, 0xa7, 0xc, 0xdf, 0x7, 0xc9, 0x5, 0x10, 0x32, 0xf, 0x12, 0x14, 0x18, 0x2, 0x22, 0x10, 0x9, 0xf8, 0x21, 0xf0, 0x13, 0x1a, 0x80, 0x1, 0xc0, 0xc, 0xd0, 0x7, 0xcf, 0x8, 0x28, 0x7, 0x12, 0x21, 0x18, 0x2, 0x22, 0x1d, 0x9, 0xb8, 0x17, 0xa8, 0x16, 0x12, 0xf8, 0x6, 0xf8, 0x8, 0xb0, 0x4, 0xb0, 0x2, 0x9, 0x67, 0xdf, 0xd, 0x2a, 0x80, 0x1, 0xc0, 0xc, 0xd0, 0x7, 0xcf, 0x8, 0x28, 0x7, 0x12, 0x9, 0x18, 0x1, 0x22, 0x5, 0x9, 0xb8, 0x17, 0xa8, 0x16, 0x12, 0xd, 0x18, 0x1, 0x22, 0x9, 0x11, 0xb8, 0x17, 0xa8, 0x16, 0xf8, 0x6, 0xf8, 0x8, 0x12, 0x17, 0x12, 0xc, 0x0, 0x0, 0x1, 0x0, 0x2, 0x1, 0x3, 0x2, 0x4, 0x1, 0x5, 0x3, 0x18, 0x1, 0x22, 0x5, 0x9, 0xf0, 0x28, 0xda, 0x15, 0x1a, 0x9, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x6b, 0x65, 0x79, 0x1a, 0x8, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x6, 0x69, 0x6e, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x9, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x6b, 0x65, 0x79, 0x1a, 0x7, 0x75, 0x69, 0x6e, 0x74, 0x6b, 0x65, 0x79, 0x1a, 0x7, 0x62, 0x6f, 0x6f, 0x6c, 0x6b, 0x65, 0x79, 0x22, 0x9, 0x19, 0x1f, 0x85, 0xeb, 0x51, 0xb8, 0xe, 0x59, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x40, 0x26, 0x50, 0x0, 0x3, 0x42, 0x22, 0x8, 0xa, 0x6, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x2, 0x38, 0x1, 0x78, 0x2}
var bytevals2, _ = ioutil.ReadFile("test_data/1206_1541_12.pbf")
var tileidval = m.TileID{1206, 1541, 12}
var features = ReadTile(bytevals2, tileidval)["osm"]
var config = NewConfig("mylayer", tileidval)

func TestWriteLayer(t *testing.T) {
	config := NewConfig("testing", m.TileID{0, 0, 0})
	bytevals := WriteLayer(fc.Features, config)
	tile1 := &vector_tile.Tile{}
	tile2 := &vector_tile.Tile{}

	err := proto.Unmarshal(bytevals, tile1)
	if err != nil {
		t.Errorf("", err)
	}

	err = proto.Unmarshal(tile_expected, tile2)
	if err != nil {
		t.Errorf("", err)
	}
	// checking to see if byte values are the same size
	if len(tile_expected) != len(bytevals) {
		t.Errorf("writelayer test failed %d %d", len(tile_expected), len(bytevals))
	} else {
		if len(tile1.Layers) != len(tile2.Layers) {
			t.Errorf("writelayer test failed %d %d", len(tile1.Layers), len(tile2.Layers))
		} else {
			for i := range tile1.Layers {
				layer1, layer2 := tile1.Layers[i], tile2.Layers[i]
				if *layer1.Name != *layer2.Name {
					t.Errorf("writelayer test failed %s %s", *layer1.Name, *layer2.Name)
				}
				if *layer1.Version != *layer2.Version {
					t.Errorf("writelayer test failed %d %d", *layer1.Version, *layer2.Version)
				}
				//t.Errorf("%v %v", layer1.Extent, layer2.Extent)
				if layer1.Extent != layer2.Extent {
					t.Errorf("writelayer test failed %d %d", layer1.Extent, layer2.Extent)
				}

				if len(layer1.Keys) != len(layer2.Keys) {
					t.Errorf("writelayer test failed %d %d", len(layer1.Keys), len(layer2.Keys))
				}

				if len(layer1.Values) != len(layer2.Values) {
					t.Errorf("writelayer test failed %d %d", len(layer1.Values), len(layer2.Values))
				}

				if len(layer1.Features) != len(layer2.Features) {
					t.Errorf("writelayer test failed %d %d", len(layer1.Features), len(layer2.Features))
				}

			}
		}
	}
}

func BenchmarkWriteLayerSmall(b *testing.B) {
	b.ReportAllocs()
	config := NewConfig("testing", m.TileID{0, 0, 0})

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		WriteLayer(fc.Features, config)
	}
}

func BenchmarkWriteFeatures(b *testing.B) {
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		WriteLayer(features, config)
	}
}

func BenchmarkWriteAll(b *testing.B) {
	b.ReportAllocs()

	filenames := []string{"./test_data/1171_1566_12.pbf", "./test_data/1206_1540_12.pbf", "./test_data/1206_1541_12.pbf", "./test_data/8801_5371_14.pbf", "./test_data/654_1583_12.pbf", "./test_data/701_1635_12.pbf", "./test_data/701_1636_12.pbf", "./test_data/701_1637_12.pbf", "./test_data/702_1636_12.pbf", "./test_data/703_1635_12.pbf", "./test_data/703_1637_12.pbf", "./test_data/9_12_5.pbf"}

	byte_array := map[m.TileID][]*geojson.Feature{}
	for _, filename := range filenames {
		vals := strings.Split(filename, "/")
		tileid := vals[len(vals)-1]
		tileid = tileid[:len(tileid)-4]
		tileid = strings.Replace(tileid, "_", "/", -1)
		newtileid := m.Strtile(tileid)
		bytevals, _ := ioutil.ReadFile(filename)
		byte_array[newtileid] = ReadTile(bytevals, newtileid)[`osm`]
	}

	for n := 0; n < b.N; n++ {
		for k, feats := range byte_array {
			WriteLayer(feats, NewConfig("testing", k))
		}
	}
}
