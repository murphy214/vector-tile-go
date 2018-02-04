# vector-tile-go
An implementation of mapbox's vector-tile-js library for reading vector tiles lazily.


# What is it 

This repo is essentially a read version of [this] is pretty well close to being implemented in the same way. Basically it reads vector-tiles lazily allowing for faster reading period (hopefully). Benchmarks are kind of all over the place its either about the same speed or like twice as fast depending on the size of your tile etc. For example here are benchmarks with one of the faster test examples. 

The first two benchmarks are just reading the Vector_Tile into a lazy structure (what this repo does) and the other is serializing it from a raw protobuf. As you can see its quite a bit faster in that regard. The second two benchmarks pertain to convertinhsig from a vector tile represention to a geojson representation. 

```
go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/murphy214/vector-tile-go
Benchmark_New_Vector_Tile-8                 	      30	  48400647 ns/op
Benchmark_New_Vector_Tile_Proto-8           	      10	 106715383 ns/op
Benchmark_New_Vector_Tile_Geojson-8         	      10	 158165582 ns/op
Benchmark_New_Vector_Tile_Proto_Geojson-8   	       5	 230668539 ns/op
PASS
ok  	github.com/murphy214/vector-tile-go	7.673s
```

# Usage 

I added one more structure to the endpoint of the Vector_Tile being ToGeoJSON() that accepts a tileid. In the mapbox's version this can only be done at the feature level. 


```golang
package main

import (
  "io/ioutil"
  vt "github.com/murphy214/vector-tile-go"
)

func main() {
  bytevals,_ :=  ioutil.ReadFile("test_data/9-12-5.pbf")
  tileid := m.TileID{9,12,5}
  tile := vt.New_Vector_Tile(bytevals) // this is your tile structure  
  layermap := tile.ToGeoJSON(tileid) // this is your layer map map[string][]*geojson.Feature
}
```

# Caveats 

I am by no means guaranteeing this is faster for parsing vector-tiles then older methods I already use, in fact theres a good chance that will end up happening this is just an experiment I guess. 

# Analysis Against Previous Implementation

I'm doing a little statistical analysis against an mbtiles qa set from mapbox, against the times ToGeoJSON verses my previous mbtiles-util implmentation, and while I it looks like 90% of the time it is slower than previous, the times it is faster are so large, that the total amount of time spent of each is lower with the new implmentation. However is it enough to justify implementing this? 

This module has a few advantages over the others:
  * its self contained no protobuf or vector_tile pb.go file required, which can get annoying with importing from go
  * more importantly the access to lower level stuff it gives, I'm sure there will be situations where I need to feature level manipulations or statististics. (see structure above)


