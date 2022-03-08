package vt

import (
	"github.com/murphy214/pbf"
	m "github.com/murphy214/mercantile"
	"fmt"
)

// the layer structure for layer
type Layer struct {
	Name             string        // name of laeyr
	Extent           int           // size of extent
	Version          int           // the correct version of a given layer
	Keys             []string      // size of keys
	Values           []interface{} // size of values
	Number_Features  int           // the number of features
	features         []int         // size of features
	StartPos         int
	EndPos           int
	feature_position int
	Buf              *pbf.PBF
	keys_ind		[2]int
	values_ind 		[2]int
	keys_bool,vals_bool bool
}	

// creates a new layer
func (tile *Tile) NewLayer(endpos int) {
	layer := &Layer{StartPos: tile.Buf.Pos, EndPos: endpos}
	key, val := tile.Buf.ReadKey()
	keys_bool := false 
	vals_bool := false
	layer.keys_ind = [2]int{-1,-1}
	for tile.Buf.Pos < layer.EndPos {
		if key == 1 && val == 2 {
			layer.Name = tile.Buf.ReadString()
			tile.Layers = append(tile.Layers, layer.Name)
			key, val = tile.Buf.ReadKey()
		}
		// collecting all the features
		for key == 2 && val == 2 {
			// reading for features

			layer.features = append(layer.features, tile.Buf.Pos)
			feat_size := tile.Buf.ReadVarint()

			tile.Buf.Pos += feat_size
			key, val = tile.Buf.ReadKey()
		}
		if key == 3 && val == 2 && !keys_bool {
			keys_bool = true 
			mypos := tile.Buf.Pos 
			layer.keys_ind[0] = mypos-1
		}


		// collecting all keys
		for key == 3 && val == 2 {
			layer.Keys = append(layer.Keys, tile.Buf.ReadString())
			key, val = tile.Buf.ReadKey()
			if (key == 3 && val == 2) {
				layer.keys_ind[1] = tile.Buf.Pos -1

			}
		}

		if key == 4 && val == 2 && !vals_bool {
			vals_bool = true 
			mypos := tile.Buf.Pos 
			layer.values_ind[0] = mypos-1
		}

		// collecting all values
		for key == 4 && val == 2 {
			//tile.Buf.Byte()
			tile.Buf.ReadVarint()
			newkey, _ := tile.Buf.ReadKey()
			switch newkey {
			case 1:
				layer.Values = append(layer.Values, tile.Buf.ReadString())
			case 2:
				layer.Values = append(layer.Values, tile.Buf.ReadFloat())
			case 3:
				layer.Values = append(layer.Values, tile.Buf.ReadDouble())
			case 4:
				layer.Values = append(layer.Values, tile.Buf.ReadInt64())
			case 5:
				layer.Values = append(layer.Values, tile.Buf.ReadUInt64())
			case 6:
				layer.Values = append(layer.Values, tile.Buf.ReadUInt64())
			case 7:
				layer.Values = append(layer.Values, tile.Buf.ReadBool())
			}
			key, val = tile.Buf.ReadKey()
			if !(key == 4 && val == 2 ) {
				layer.values_ind[1] = tile.Buf.Pos -1

			}
		}
		if key == 5 && val == 0 {
			layer.Extent = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()
		}
		if key == 15 && val == 0 {
			layer.Version = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()

		}
	}

	if layer.Extent == 0 {
		layer.Extent = 4096
	}
	layer.Number_Features = len(layer.features)
	tile.LayerMap[layer.Name] = layer
	tile.Buf.Pos = endpos
	layer.Buf = tile.Buf
}

func (layer *Layer) Next() bool {
	return layer.feature_position < layer.Number_Features
}

func (layer *Layer) Reset() {
	layer.feature_position = 0
}


/*
func (layer *Layer) WriteLayer(tileid m.TileID) *WriteLayer {

	&LayerWrite{
		Name:layer.Name,
		Exten:layer.Extent,
		Version:layer.Version,

	}S
|

*/
func (layer *Layer) ToLayerWrite(tileid m.TileID) (*LayerWrite,error) {
	// creating cursor 
	cur := NewCursorExtent(tileid,4326)
	
	// getting the last feature 
	layer.feature_position = layer.Number_Features-1

	feat,err := layer.Feature()
	if err != nil {
		return nil,err
	}
	
	geom,err := feat.LoadGeometry()
	if err != nil {
		return nil,err
	}

	last_pt := get_last_point(geom)
	if len(last_pt) == 2 {
		cur.LastPoint = []int32{int32(last_pt[0]),int32(last_pt[1])}
	}
	//fmt.Println(cur,get_last_point(geom),"we here")

	// getting the bytes assocated with features
	var feat_bytes []byte
	if len(layer.features) > 0 {
		start_pos := layer.features[0] 
		
		layer.Buf.Pos = layer.features[len(layer.features)-1]
		//layer.Buf.Pos = layer.features[len(layer.features)-1]
		fmt.Println(layer.Buf.Pos,layer.Buf.Pos,layer.Buf.Pbf[layer.Buf.Pos-3:layer.Buf.Pos+25],layer.Buf.Pbf[layer.Buf.Pos])

		end_pos := layer.Buf.Pos + int(layer.Buf.ReadVarint())
		feat_bytes = layer.Buf.Pbf[start_pos-1:end_pos]
		fmt.Println(feat_bytes)
		fmt.Println(start_pos,layer.Buf.Pos,layer.Buf.Pbf[start_pos-3:start_pos+3],layer.Buf.Pbf[start_pos])
	} else {
		feat_bytes = []byte{}
	}

	// creating the keys map
	keymap := map[string]uint32{}
	for pos,key := range layer.Keys {
		keymap[key] = uint32(pos)
	}

	// creeating values map
	valuemap := map[interface{}]uint32{}
	for pos,value := range layer.Values {
		valuemap[value] = uint32(pos)
	}
	
	bds := m.Bounds(tileid)
	return &LayerWrite{
		Name:layer.Name,
		Extent:layer.Extent,
		Version:layer.Version,
		TileID:tileid,
		Keys_Bytes: layer.Buf.Pbf[layer.keys_ind[0]:layer.keys_ind[1]],
		Values_Bytes: layer.Buf.Pbf[layer.values_ind[0]:layer.values_ind[1]],
		Features: feat_bytes,
		Values_Map: valuemap,
		Keys_Map: keymap,
		Cursor:cur,
		DeltaX: bds.E - bds.W,
		DeltaY: bds.N - bds.S,
	},nil
}
