package vt

import (
	"github.com/murphy214/pbf"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/vector-tile-go/tags"
)

// Formula for values in dimension:
// value = base + multiplier * (delta_encoded_value + offset)
type Scaling struct {
	Offset int // default = 0,0
	Multiplier float64 // default = 1
	Base float64 // default = 0.0
}

// creates a new scaling object
func NewScaling() *Scaling {
	return &Scaling{Multiplier:0}
}
 
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

	// stuff added in v3
	TagReader *tags.TagReader // the tag reader for the v3 spec
	ElevationScaling *Scaling // the scaling dimension fro version 3
	AttributeScalings []*Scaling // attribute scalingas
	TileID m.TileID
}

// creates a new layer
func (tile *Tile) NewLayer(endpos int) {
	layer := &Layer{StartPos: tile.Buf.Pos, EndPos: endpos,TagReader:&tags.TagReader{},ElevationScaling:NewScaling()}
	key, val := tile.Buf.ReadKey()
	for tile.Buf.Pos < layer.EndPos {
		if key == 1 && val == 2 {
			layer.Name = tile.Buf.ReadString()
			tile.Layers = append(tile.Layers, layer.Name)
			key, val = tile.Buf.ReadKey()
		}
		// collecting all the features
		for key == 2 && val == 2 {
			// reading through features for layer
			layer.features = append(layer.features, tile.Buf.Pos)
			feat_size := tile.Buf.ReadVarint()			
			tile.Buf.Pos += feat_size
			key, val = tile.Buf.ReadKey()
		}
		// collecting all keys
		for key == 3 && val == 2 {
			layer.Keys = append(layer.Keys, tile.Buf.ReadString())
			key, val = tile.Buf.ReadKey()
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
		}
		
		// reading extent
		if key == 5 && val == 0 {
			layer.Extent = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()
		}

		// reading all the string values 
		for key == 6 && val == 2 {
			layer.TagReader.StringValues = append(layer.TagReader.StringValues, tile.Buf.ReadString())
			key, val = tile.Buf.ReadKey()
		}

		// reading all the float32 values 
		if key == 7 && val == 2 {
			size := tile.Buf.ReadVarint()
			endpos := tile.Buf.Pos + size
			for tile.Buf.Pos < endpos {
				layer.TagReader.FloatValues = append(layer.TagReader.FloatValues, tile.Buf.ReadFloat())
			}
			key, val = tile.Buf.ReadKey()
		}

		// reading all the double values 
		if key == 8 && val == 2 {
			size := tile.Buf.ReadVarint()
			endpos := tile.Buf.Pos + size
			for tile.Buf.Pos < endpos {
				layer.TagReader.DoubleValues = append(layer.TagReader.DoubleValues, tile.Buf.ReadDouble())
			}
			key, val = tile.Buf.ReadKey()
		}

		// reading the fixed ujint64 values 
		if key == 9 && val == 2 {
			size := tile.Buf.ReadVarint()
			endpos := tile.Buf.Pos + size
			for tile.Buf.Pos < endpos {
				layer.TagReader.IntValues = append(layer.TagReader.IntValues, int(tile.Buf.ReadFixed64()))
			}
			key, val = tile.Buf.ReadKey()
		}

		// reading the scaling for this layer
		if key ==  10 && val == 2 {
			// getting size and end position
			tile.Buf.ReadVarint()
			//endpos := tile.Buf.Pos + size 
			key,val = tile.Buf.ReadKey()
			
			// reading off set
			if key == 1 && val == 0 {
				layer.ElevationScaling.Offset = int(tile.Buf.ReadSVarint())			
				key,val = tile.Buf.ReadKey()
			}
			// reading multiplier
			if key == 2 && val == 1 {
				layer.ElevationScaling.Multiplier = tile.Buf.ReadDouble()			
				key,val = tile.Buf.ReadKey()
			}
			// reading multiplier
			if key == 3 && val == 1 {
				layer.ElevationScaling.Base = tile.Buf.ReadDouble()			
				key,val = tile.Buf.ReadKey()
			}

		}	
		
		// reading the dense attribute scalings
		for key ==  11 && val == 2 {
			size := tile.Buf.ReadVarint()
			endpos := tile.Buf.Pos + size 
			
			for tile.Buf.Pos < endpos {
				//endpos := tile.Buf.Pos + size 
				key,val = tile.Buf.ReadKey()
				tempscaling := NewScaling()
				
				// reading off set
				if key == 1 && val == 0 {
					tempscaling.Offset = int(tile.Buf.ReadSVarint())			
					key,val = tile.Buf.ReadKey()
				}
				// reading multiplier
				if key == 2 && val == 1 {
					tempscaling.Multiplier = tile.Buf.ReadDouble()			
					key,val = tile.Buf.ReadKey()
				}
				// reading multiplier
				if key == 3 && val == 1 {
					tempscaling.Base = tile.Buf.ReadDouble()			
					key,val = tile.Buf.ReadKey()				
				}
				layer.AttributeScalings = append(layer.AttributeScalings,tempscaling)
			}
			//key,val = tile.Buf.ReadKey()				
			tile.Buf.Pos = endpos
			key,val = tile.Buf.ReadKey()				
		}
 

		// operation to get tileid out if applicable
		var tilex,tiley,tilez int
		// reading tilex 
		if key == 12 && val == 0 {
			tilex = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()
		}
		// reading tiley 
		if key == 13 && val == 0 {
			tiley = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()
		}

		// reading tilex 
		if key == 14 && val == 0 {
			tilez = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()
		}

		// adding tileid to layer
		layer.TileID = m.TileID{int64(tilex),int64(tiley),uint64(tilez)}




		// reading version
		if key == 15 && val == 0 {
			layer.Version = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()

		}
	}
	if layer.Extent == 0 {
		layer.Extent = 4096
	}
	layer.Number_Features = len(layer.features)
	layer.TagReader.Keys = layer.Keys
	tile.LayerMap[layer.Name] = layer
	tile.Buf.Pos = endpos
	layer.Buf = tile.Buf
}

// increments to the next feature
func (layer *Layer) Next() bool {
	return layer.feature_position < layer.Number_Features
}

// resets a layer
func (layer *Layer) Reset() {
	layer.feature_position = 0
}
