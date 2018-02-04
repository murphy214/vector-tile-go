package pbf

import (
	"io/ioutil"
	"testing"
	"github.com/murphy214/mbtiles-util/vector-tile/2.1"
	"github.com/golang/protobuf/proto"
	mbutil "github.com/murphy214/mbtiles-util"
	m "github.com/murphy214/mercantile"
)

var bytevals,_ =  ioutil.ReadFile("test_data/9-12-5.pbf")
var tileid = m.TileID{9,12,5}


// benchamrks every new vector tile
func Benchmark_New_Vector_Tile(b *testing.B) {
        // run the Fib function b.N times
        for n := 0; n < b.N; n++ {
         	New_Vector_Tile(bytevals) 	
        }
}


// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto(b *testing.B) {
		tile := &vector_tile.Tile{}
        // run the Fib function b.N times
        for n := 0; n < b.N; n++ {
         	proto.Unmarshal(bytevals,tile) 	
        }
}


// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Geojson(b *testing.B) {
        // run the Fib function b.N times
        for n := 0; n < b.N; n++ {
         	New_Vector_Tile(bytevals).ToGeoJSON(tileid)
        }
}



// benchamrks every new vector tile
func Benchmark_New_Vector_Tile_Proto_Geojson(b *testing.B) {
        // run the Fib function b.N times
        for n := 0; n < b.N; n++ {
         	mbutil.Convert_Vt_Bytes(bytevals,tileid)
        }
}



