package pbf

import (
	"fmt"
	//"github.com/paulmach/go.geojson"
	// m "github.com/murphy214/mercantile"
)

// fuck debugging
func b() {
	fmt.Println()
}

// the layer structure for layer 
type Layer struct {
	Name string // name of laeyr
	Extent int // size of extent
	Version int // the correct version of a given layer
	Keys []string // size of keys
	Values []interface{} // size of values
	Number_Features int // the number of features
	features []int // size of features
	Buf *PBF // the buffer associated with a layer 
}


// creates a new layer
func (layer *Layer) New_Layer() {
	key,val := layer.Buf.ReadKey()

	for layer.Buf.Pos < layer.Buf.Length {
		if key == 1 && val == 2 {
			layer.Name = layer.Buf.ReadString()
			key,val = layer.Buf.ReadKey()
		}
		// collecting all the features
		for key == 2 && val == 2 {
			// reading for features

			layer.features = append(layer.features,layer.Buf.Pos)
			feat_size := layer.Buf.ReadVarint()

			layer.Buf.Pos += feat_size
			key,val = layer.Buf.ReadKey()
		}
		// collecting all keys
		for key == 3 && val == 2 {
			layer.Keys = append(layer.Keys,layer.Buf.ReadString())
			key,val = layer.Buf.ReadKey()
		}
		// collecting all values
		for key == 4 && val == 2 {
			//layer.Buf.Byte()
			layer.Buf.ReadVarint()
			newkey,_ := layer.Buf.ReadKey()
			if newkey == 1 {
				layer.Values = append(layer.Values,layer.Buf.ReadString())			
			} else if newkey == 2 {
				layer.Values = append(layer.Values,layer.Buf.ReadFloat())
			} else if newkey == 3 {
				layer.Values = append(layer.Values,layer.Buf.ReadDouble())
			} else if newkey == 4 {
				layer.Values = append(layer.Values,layer.Buf.ReadInt64())			
			} else if newkey == 5 {
				layer.Values = append(layer.Values,layer.Buf.ReadUInt64())			
			} else if newkey == 6 {
				layer.Values = append(layer.Values,layer.Buf.ReadUInt64())					
			} else if newkey == 7 {
				layer.Values = append(layer.Values,layer.Buf.ReadBool())			
			}
			

			key,val = layer.Buf.ReadKey()
		}
		if key == 5 && val == 0 {
			layer.Extent = int(layer.Buf.ReadVarint())
			key,val = layer.Buf.ReadKey()
		}
		if key == 15 && val == 0 {
			layer.Version = int(layer.Buf.ReadVarint())
			key,val = layer.Buf.ReadKey()

		}
		//fmt.Println(layer.Buf.Pos,layer.Buf.Pbf,layer.Buf.Length)	
	}
	layer.Buf = layer.Buf
	layer.Number_Features = len(layer.features)
}

func (layer Layer) Feature_Raw(pos int) *Feature {
	layer.Buf.Pos = layer.features[pos]
	endpos := layer.Buf.Pos + layer.Buf.ReadVarint()
	startpos := layer.Buf.Pos
	//&PBF{PBF:layer.Buf.Pbf[startpos:endpos],Length:endpos-startpos}
	feature := &Feature{
		Buf:
			&PBF{
				Pbf:layer.Buf.Pbf[startpos:endpos],
				Length:(endpos-startpos)+1,
			},
		Properties:map[string]interface{}{},
	}
	feature.LayerFeature_Raw(layer.Keys,layer.Values,0)
	feature.Extent = layer.Extent
	return feature
}













