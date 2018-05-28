# vector-tile-go
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/murphy214/vector-tile-go)

An implementation of mapbox's vector-tile spec for reading / writing vector tiles. This is implementation is protocol buffer less meaning the serialization and deserialization is handled by me, which saves some needless allocations as well as pretty signicant performance gains.

## Why?

An implementation of mapbox's vector tile spec from with no protobuf file needed. Designed to reduce allocations and to be faster for reading / writing. When you think about how a vector tile is marshaled / unmarshaled from a regular protobuf implementation its kind of ridiculous, we transform (generally) geojson data into another entire feature struct with an entire new allocation for each, only to serialize that data structure to bytes immediately after. Thats really heavy, and unneeded allocations.

## Features 

* Reads vector tiles into layer map geojson with a tile id given

* Can lazily read each feature within a layer so maps can be peformed based on only wanting certain properties and geometry types 

* Writes vector tile layers that can be appended to an existing vector tile byte array (this is just as useful than just supporting the full spec my api just handles one layer at a time)

* Writers vector tile layers for geojson-vt feature sets (formatted quite differently)

* Writes vector tile layers for geobuf data sets

* **Reads to geojson are 40% faster than a regular proto implementation**

* **Writes from geojson are 130% faster than regular proto implmentation**

* **Lazy reads are 70% faster than the protobuf implementation (i.e. marshaling the tile)**

# Caveats 

The api is currently still fluid and subject to change but currently everything I throw at works which I find to be pretty good, as I'm implementing several processes in one library. 
