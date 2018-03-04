package vt

import (
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	//"fmt"
)

// the defualt layer structure from a layer
type LayerWrite struct {
	TileID m.TileID // the tile id associated with the layer
	Name string // the name associated with the layer
	Extent int // extent will assume 4096 if 0 
	Version int // version number will assume 15 if 0
	Keys_Map map[string]uint32
	Keys_Bytes []byte // the byte value of keys
	Values_Map map[interface{}]uint32 // the values map
	Values_Bytes []byte // the byte values of values
	Features []byte // the byte values of features
	Cursor *Cursor // the cursor for adding geometries
}

// the configuration struct
type Config struct {
	TileID m.TileID // the tile id associated with the layer
	Name string // the name associated with the layer
	Extent int32 // extent will assume 4096 if 0 
	Version int // version number will assume 15 if 0
}

// creates a new layer
func NewLayer(tileid m.TileID,name string) LayerWrite {
	keys_map := map[string]uint32{}
	values_map := map[interface{}]uint32{}
	cur := NewCursor(tileid)
	return LayerWrite{TileID:tileid,Keys_Map:keys_map,Values_Map:values_map,Name:name,Cursor:cur}
}

// a function for creatnig a new conifguratoin
func NewConfig(filename string,tileid m.TileID) Config {
	return Config{Name:filename,TileID:tileid}
}

// creates a layer from a configuration 
func NewLayerConfig(config Config) LayerWrite {
	keys_map := map[string]uint32{}
	values_map := map[interface{}]uint32{}
	if config.Extent == int32(0) {
		config.Extent = int32(4096)
	}
	if config.Version == 0 {
		config.Version = 2
	}
	cur := NewCursorExtent(config.TileID,config.Extent)
	return LayerWrite{TileID:config.TileID,
					Keys_Map:keys_map,
					Values_Map:values_map,
					Name:config.Name,
					Cursor:cur,
					Version:config.Version,
					Extent:int(config.Extent),}
}

// adds a single key to a given layer 
func (layer *LayerWrite) AddKey(key string) uint32 {
	layer.Keys_Bytes = append(layer.Keys_Bytes,26)
	layer.Keys_Bytes = append(layer.Keys_Bytes,EncodeVarint(uint64(len(key)))...)
	layer.Keys_Bytes = append(layer.Keys_Bytes,[]byte(key)...)
	myint := uint32(len(layer.Keys_Map))
	layer.Keys_Map[key] = myint
	return myint
}

// adds a single value to a given
func (layer *LayerWrite) AddValue(value interface{}) uint32 {
	layer.Values_Bytes = append(layer.Values_Bytes,WriteValue(value)...)
	myint := uint32(len(layer.Values_Map))
	layer.Values_Map[value] = myint
	return myint
}

// gets the tags for a given set of properties
func (layer *LayerWrite) GetTags(properties map[string]interface{}) []uint32 {
	tags := make([]uint32,len(properties)*2)
	i := 0
	for k,v := range properties {
		keytag,keybool := layer.Keys_Map[k]
		if keybool == false {
			keytag = layer.AddKey(k)
		}
		valuetag,valuebool := layer.Values_Map[v]
		if valuebool == false {
			valuetag = layer.AddValue(v)
		}
		tags[i] = keytag
		tags[i+1] = valuetag
		i += 2
	}
	return tags
}

// refreshes the cursor
func (layer *LayerWrite) RefreshCursor() {
	layer.Cursor.Count = 0
	layer.Cursor.LastPoint = []int32{0,0}
	layer.Cursor.Geometry = []uint32{}
}

// creates a layer outright using a configuration and a set of features
// this is the outermost function that should be used
// the outer functions is wrapped like this to reduce allocations
// if it was used as method it could cause leaks which I'll have to check 
// later via escape analysis
func WriteLayer(features []*geojson.Feature,config Config) []byte {
	// creating layer
	mylayer := NewLayerConfig(config)

	// adding features
	for _,feat := range features {
		mylayer.AddFeature(feat)
	}

	// creating total_bytes
	total_bytes := []byte{}

	// writing name
	if len(mylayer.Name) > 0 {
		total_bytes = append(total_bytes,10)
		total_bytes = append(total_bytes,EncodeVarint(uint64(len(mylayer.Name)))...)
		total_bytes = append(total_bytes,[]byte(mylayer.Name)...)
	}
	
	// appending features
	total_bytes = append(total_bytes,mylayer.Features...)

	// appending keys 
	total_bytes = append(total_bytes,mylayer.Keys_Bytes...)

	// appending values 
	total_bytes = append(total_bytes,mylayer.Values_Bytes...)

	// appending extra config values
	if mylayer.Extent != 4096 {
		total_bytes = append(total_bytes,5)
		total_bytes = append(total_bytes,EncodeVarint(uint64(mylayer.Extent))...)
	}
	
	//if mylayer.Version != 0 {
	total_bytes = append(total_bytes,120)
	total_bytes = append(total_bytes,byte(mylayer.Version))
	//}
	beg := append([]byte{26},EncodeVarint(uint64(len(total_bytes)))...)
	return append(beg,total_bytes...)
}


