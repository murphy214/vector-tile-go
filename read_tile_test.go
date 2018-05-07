package vt

import (
	"github.com/golang/protobuf/proto"
	mbutil "github.com/murphy214/mbtiles-util"
	"github.com/murphy214/mbtiles-util/vector-tile/2.1"
	m "github.com/murphy214/mercantile"
	"io/ioutil"
	"strings"
	//"sync"
	//"fmt"
	"testing"
)

var bytevals, _ = ioutil.ReadFile("test_data/701_1635_12.pbf")
var tileid = m.TileID{701, 1635, 12}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile(b *testing.B) {
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		NewTile(bytevals)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto(b *testing.B) {
	b.ReportAllocs()

	tile := &vector_tile.Tile{}
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		proto.Unmarshal(bytevals, tile)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Geojson(b *testing.B) {
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		ReadTile(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto_Geojson(b *testing.B) {
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		mbutil.Convert_Vt_Bytes(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Geojson_1(b *testing.B) {

	bytevals, _ = ioutil.ReadFile("test_data/701_1635_12.pbf")
	tileid = m.TileID{701, 1635, 12}
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		ReadTile(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto_Geojson_1(b *testing.B) {
	bytevals, _ = ioutil.ReadFile("test_data/701_1635_12.pbf")
	tileid = m.TileID{701, 1635, 12}

	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		mbutil.Convert_Vt_Bytes(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Geojson_2(b *testing.B) {

	bytevals, _ = ioutil.ReadFile("test_data/701_1637_12.pbf")
	tileid = m.TileID{701, 1637, 12}
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		ReadTile(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto_Geojson_2(b *testing.B) {

	bytevals, _ = ioutil.ReadFile("test_data/701_1637_12.pbf")
	tileid = m.TileID{701, 1637, 12}

	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		mbutil.Convert_Vt_Bytes(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Geojson_3(b *testing.B) {
	bytevals, _ = ioutil.ReadFile("test_data/702_1636_12.pbf")
	tileid = m.TileID{702, 1636, 12}
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		ReadTile(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto_Geojson_3(b *testing.B) {

	bytevals, _ = ioutil.ReadFile("test_data/703_1635_12.pbf")
	tileid = m.TileID{703, 1635, 12}

	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		mbutil.Convert_Vt_Bytes(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Geojson_4(b *testing.B) {

	bytevals, _ = ioutil.ReadFile("test_data/703_1635_12.pbf")
	tileid = m.TileID{703, 1635, 12}
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		ReadTile(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto_Geojson_4(b *testing.B) {

	bytevals, _ = ioutil.ReadFile("test_data/703_1635_12.pbf")
	tileid = m.TileID{703, 1635, 12}

	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		mbutil.Convert_Vt_Bytes(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Geojson_5(b *testing.B) {

	bytevals, _ = ioutil.ReadFile("test_data/703_1637_12.pbf")
	tileid = m.TileID{703, 1637, 12}
	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		ReadTile(bytevals, tileid)
	}
}

// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto_Geojson_5(b *testing.B) {

	bytevals, _ = ioutil.ReadFile("test_data/703_1637_12.pbf")
	tileid = m.TileID{703, 1637, 12}

	b.ReportAllocs()

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		mbutil.Convert_Vt_Bytes(bytevals, tileid)
	}
}

/*
func Benchmark_ReadPacked_Old(b *testing.B) {
	b.ReportAllocs()

	pbf := pbf.PBF{Pbf: []byte{27, 0x0, 0xc, 0x1, 0xeb, 0x2, 0x2, 0xec, 0x2, 0x3, 0xec, 0x2, 0x4, 0xec, 0x2, 0x5, 0xec, 0x2, 0x6, 0xec, 0x2, 0x7, 0xed, 0x2, 0x8, 0x4, 0x9, 0x5}, Length: 27}
	for n := 0; n < b.N; n++ {
		pbf.ReadPackedUInt32()
		pbf.Pos = 0
	}
}

func Benchmark_ReadPacked_New(b *testing.B) {
	b.ReportAllocs()

	pbf := pbf.PBF{Pbf: []byte{27, 0x0, 0xc, 0x1, 0xeb, 0x2, 0x2, 0xec, 0x2, 0x3, 0xec, 0x2, 0x4, 0xec, 0x2, 0x5, 0xec, 0x2, 0x6, 0xec, 0x2, 0x7, 0xed, 0x2, 0x8, 0x4, 0x9, 0x5}, Length: 27}
	for n := 0; n < b.N; n++ {
		pbf.ReadPacked()
		pbf.Pos = 0
	}
}
*/

func Benchmark_All_Non_Proto(b *testing.B) {
	b.ReportAllocs()

	filenames := []string{"./test_data/1171_1566_12.pbf", "./test_data/1206_1540_12.pbf", "./test_data/1206_1541_12.pbf", "./test_data/8801_5371_14.pbf", "./test_data/654_1583_12.pbf", "./test_data/701_1635_12.pbf", "./test_data/701_1636_12.pbf", "./test_data/701_1637_12.pbf", "./test_data/702_1636_12.pbf", "./test_data/703_1635_12.pbf", "./test_data/703_1637_12.pbf"}

	byte_array := map[m.TileID][]byte{}
	for _, filename := range filenames {
		vals := strings.Split(filename, "/")
		tileid := vals[len(vals)-1]
		tileid = tileid[:len(tileid)-4]
		tileid = strings.Replace(tileid, "_", "/", -1)
		newtileid := m.Strtile(tileid)
		bytevals, _ := ioutil.ReadFile(filename)
		byte_array[newtileid] = bytevals
	}

	for n := 0; n < b.N; n++ {

		for k, bytevals := range byte_array {
			mbutil.Convert_Vt_Bytes(bytevals, k)
		}
	}
}

func Benchmark_All_Non(b *testing.B) {
	b.ReportAllocs()

	filenames := []string{"./test_data/1171_1566_12.pbf", "./test_data/1206_1540_12.pbf", "./test_data/1206_1541_12.pbf", "./test_data/8801_5371_14.pbf", "./test_data/654_1583_12.pbf", "./test_data/701_1635_12.pbf", "./test_data/701_1636_12.pbf", "./test_data/701_1637_12.pbf", "./test_data/702_1636_12.pbf", "./test_data/703_1635_12.pbf", "./test_data/703_1637_12.pbf"}

	byte_array := map[m.TileID][]byte{}
	for _, filename := range filenames {
		vals := strings.Split(filename, "/")
		tileid := vals[len(vals)-1]
		tileid = tileid[:len(tileid)-4]
		tileid = strings.Replace(tileid, "_", "/", -1)
		newtileid := m.Strtile(tileid)
		bytevals, _ := ioutil.ReadFile(filename)
		byte_array[newtileid] = bytevals
	}

	for n := 0; n < b.N; n++ {
		for k, bytevals := range byte_array {
			ReadTile(bytevals, k)
		}
	}
}
