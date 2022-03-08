package vt

import (
	"github.com/murphy214/pbf"
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
	keys_bool true
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
			mypos := tile.layer.Buf.Pos 
			layer.keys_ind[0] = mypos
		}


		// collecting all keys
		for key == 3 && val == 2 {
			layer.Keys = append(layer.Keys, tile.Buf.ReadString())
			key, val = tile.Buf.ReadKey()
			if (key == 3 && val == 2) {
				layer.keys_ind[1] = tile.Buf.Pos 

			}
		}

		if key == 4 && val == 2 && !vals_bool {
			vals_bool = true 
			mypos := tile.layer.Buf.Pos 
			layer.keys_ind[0] = mypos
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
				layer.values_ind[1] = tile.Buf.Pos 

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