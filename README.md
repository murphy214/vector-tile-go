# vector-tile-go
An implementation of mapbox's vector-tile-js library for reading vector tiles lazily.

# Updates

Recently a went down a rabbit hole to track down lack of performance in tile-reduce one of my other repos and I noticed this being the case no matter what vector-tile reader I used. Anyway much performance profiling with pprof / escape analysis of code later I've reduced allocations dramatically. The problem was I used a lot of methods on structs to create my vector tile and when that occurs it leaves things to be held in memory much longer than needed. 

# What is it 

This repo is essentially a read version of [this](https://github.com/mapbox/vector-tile-js) is pretty well close to being implemented in the same way. Basically it reads vector-tiles lazily allowing for faster reading period (hopefully). Benchmarks are kind of all over the place its either about the same speed or like twice as fast depending on the size of your tile etc. For example here are benchmarks with one of the faster test examples. 

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

# Structure / API 

To use the API like mapbox's implementation is a bit janky to be honest. However this jankiness, is what makes it worth using. (hopefully) 

### Vector_Tile 

A vector tile is just a map[string]Layer of where each key is a layer name. It has the top level method .ToGeoJSON to convert the whole tile.

### Layer 

The layer level is where things get intersting here you have access to all the keys / values, the layer name, the finally version & extent, the lazy structure is found in the .features private field that stores integer positions of starting features. Thoses positions are then used to get features using ```layer.Feature(pos int)``` to get the feature at a given position in the layer.

### Feature

The feature layer is where most of the compute is done at is the methods ```ToGeoJSON(tileid)``` and ```LoadGeometry(tileid)``` although there may be reasons to change or add a method later not requiring a tileid, for now it does. The cool part here is you have access to feature properties,ids, and their respective types, without loading the geometry with a very sparse structure. So one could lazily read in a feature perform a mapping on it (like an osm mapping with a bunch of field / geometry type filters) and then read in the geojson feature if it met the criteria. 


# Analysis Against Previous Implementation

I'm doing a little statistical analysis against an mbtiles qa set from mapbox, against the times ToGeoJSON verses my previous mbtiles-util implmentation, and while I it looks like 90% of the time it is slower than previous, the times it is faster are so large, that the total amount of time spent of each is lower with the new implmentation. However is it enough to justify implementing this? 

This module has a few advantages over the others:
  * its self contained no protobuf or vector_tile pb.go file required, which can get annoying with importing from go
  * more importantly the access to lower level stuff it gives, I'm sure there will be situations where I need to feature level manipulations or statististics. (see structure above)


