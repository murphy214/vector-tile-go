package vt 


import (
	m "github.com/murphy214/mercantile"
   "github.com/paulmach/go.geojson"
)

// upper vector tile structure
type Tile struct{
	LayerMap map[string]*Layer
	Buf *PBF
	TileID m.TileID
}

// create / reads a new vector tile from a byte array 
func NewTile(bytevals []byte) *Tile {
	// creating vector tile
	tile := &Tile{
		LayerMap:map[string]*Layer{},
		Buf:&PBF{Pbf:bytevals,Length:len(bytevals)},
	}
	for tile.Buf.Pos < tile.Buf.Length {
		key,val := tile.Buf.ReadKey()
		if key == 3 && val == 2 {
			size := tile.Buf.ReadVarint()
			if size != 0 {
				tile.NewLayer(tile.Buf.Pos+size)
			}

		}	
	}
	return tile
}

// reads a tile as lazily as possible
func ReadTile(bytevals []byte,tileid m.TileID) map[string][]*geojson.Feature {
	// getting tile
	tile := NewTile(bytevals)
	tile.TileID = tileid

	// creating layermap 
	layermap := map[string][]*geojson.Feature{}

	// iterating through each layer
	for layername,v := range tile.LayerMap {
		// creating each layer in the map
		layermap[layername] = make([]*geojson.Feature,v.Number_Features)
		
		// iterating through each feature
		for i,pos := range v.features {
			layermap[layername][i] = tile.Feature(layername,v.Keys,v.Values,v.Extent,pos)
		}
	}

	return layermap
}