package pbf

//import "fmt"

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
func New_Layer(layer_bytes *PBF) *Layer {
	layer := &Layer{}

	key,val := layer_bytes.ReadKey()
	if key == 1 && val == 2 {
		layer.Name = layer_bytes.ReadString()
		key,val = layer_bytes.ReadKey()
	}
	// collecting all the features
	for key == 2 && val == 2 {
		// reading for features

		layer.features = append(layer.features,layer_bytes.Pos)
		feat_size := layer_bytes.ReadVarint()

		layer_bytes.Pos += feat_size
		key,val = layer_bytes.ReadKey()
	}
	// collecting all keys
	for key == 3 && val == 2 {
		layer.Keys = append(layer.Keys,layer_bytes.ReadString())
		key,val = layer_bytes.ReadKey()
	}
	// collecting all values
	for key == 4 && val == 2 {
		//layer_bytes.Byte()
		layer_bytes.ReadVarint()
		newkey,_ := layer_bytes.ReadKey()
		if newkey == 1 {
			layer.Values = append(layer.Values,layer_bytes.ReadString())			
		} else if newkey == 2 {
			layer.Values = append(layer.Values,layer_bytes.ReadFloat())
		} else if newkey == 3 {
			layer.Values = append(layer.Values,layer_bytes.ReadDouble())
		} else if newkey == 4 {
			layer.Values = append(layer.Values,layer_bytes.ReadInt64())			
		} else if newkey == 5 {
			layer.Values = append(layer.Values,layer_bytes.ReadUInt64())			
		} else if newkey == 6 {
			layer.Values = append(layer.Values,layer_bytes.ReadUInt64())					
		} else if newkey == 7 {
			layer.Values = append(layer.Values,layer_bytes.ReadBool())			
		}
		

		key,val = layer_bytes.ReadKey()
	}
	if key == 5 && val == 0 {
		layer.Extent = int(layer_bytes.ReadVarint())
		key,val = layer_bytes.ReadKey()
	}
	if key == 15 && val == 0 {
		layer.Version = int(layer_bytes.ReadVarint())
	}	

	layer.Buf = layer_bytes
	layer.Number_Features = len(layer.features)
	return layer
}

func (layer Layer) Feature(pos int) *Feature {
	layer.Buf.Pos = layer.features[pos]
	endpos := layer.Buf.Pos + layer.Buf.ReadVarint()
	feature := LayerFeature(layer.Buf,endpos,layer.Keys,layer.Values)
	feature.Extent = layer.Extent
	return feature
}










