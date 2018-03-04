# vector-tile-go
An implementation of mapbox's vector-tile-js library for reading vector tiles lazily.

# What is it 

An implementation of mapbox's vector tile spec from with no protobuf file needed. Designed to reduce allocations and to be faster for reading / writing. 

# Features 

* Reads vector tiles into layer map geojson with a tile id given

* Can lazily read each feature within a layer so maps can be peformed based on only wanting certain properties and geometry types 

* Writes vector tile layers that can be appended to an existing vector tile byte array (this is just as useful than just supporting the full spec my api just handles one layer at a time)

* Reads to geojson are 30% faster than a regular proto implementation\

* Writes from geojson are 130% faster than regular proto implmentation

* Lazy reads that can be used to map each from tile structure are 70% faster than the protobuf implementation (i.e. marshaling the tile) 

