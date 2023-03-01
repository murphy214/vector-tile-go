package vt

import (
	"github.com/murphy214/pbf"
	m "github.com/murphy214/mercantile"
	"fmt"
	"time"
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


var DEBUG = false
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
			if (DEBUG) {
				fmt.Println(key,val,"15")
			}
		}
		// collecting all the features
		for key == 2 && val == 2 {
			// reading for features

			layer.features = append(layer.features, tile.Buf.Pos)
			feat_size := tile.Buf.ReadVarint()

			tile.Buf.Pos += feat_size
			key, val = tile.Buf.ReadKey()
			if (DEBUG) {
				fmt.Println(key,val,"15")
			}
		}


		if key == 3 && val == 2 && !keys_bool {
			keys_bool = true 
			mypos := tile.Buf.Pos 
			layer.keys_ind[0] = mypos-1
		}

		spos := tile.Buf.Pos 

		// collecting all keys
		for key == 3 && val == 2 {
			layer.keys_ind[0] = tile.Buf.Pos 

			layer.Keys = append(layer.Keys, tile.Buf.ReadString())
			key, val = tile.Buf.ReadKey()
			if (DEBUG) {
				fmt.Println(key,val,layer.Keys,"3 2")
			}
			if !(key == 3 && val == 2) {

			}
		}
		layer.keys_ind = [2]int{spos,tile.Buf.Pos}


		if key == 4 && val == 2 && !vals_bool {
			vals_bool = true 
			mypos := tile.Buf.Pos 
			layer.values_ind[0] = mypos-1
		}
		spos = tile.Buf.Pos 

		// collecting all values
		for key == 4 && val == 2 {
			//tile.Buf.Byte()
			// spos_vals = tile.Buf.Pos 

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
			if (DEBUG) {
				fmt.Println(key,val,"4 2")
			}
			if !(key == 4 && val == 2 ) {
				layer.values_ind[1] = tile.Buf.Pos -1
			}
		}
		layer.values_ind = [2]int{spos,tile.Buf.Pos}

		if key == 5 && val == 0 {
			layer.Extent = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()
			if (DEBUG) {
				fmt.Println(key,val,"0 5")
			}
		}
		if key == 15 && val == 0 {
			layer.Version = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()
			if (DEBUG) {
				fmt.Println(key,val,"15")
			}
		}
		if (DEBUG) {
			fmt.Println(key,val,tile.Buf.Pos,endpos)
			time.Sleep(1*time.Second)
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


func (layer *Layer) FeaturePos(i int) (*Feature,error) {
	layer.feature_position = i 
	return layer.Feature()
}

func (layer *Layer) Reset() {
	layer.feature_position = 0
}

// takes a lazy reader struct and converts into a writer 
func (layer *Layer) ToLayerWrite(tileid m.TileID) (*LayerWrite,error) {
	myvals := layer.values_ind  
	myvals2 := layer.keys_ind 
	vals,vale := myvals[0],myvals[1]
	keys,keye := myvals2[0],myvals2[1]
	size_val := vale-vals
	size_key := keye-keys

	// sval,eval := layer.Buf.Pbf[myvals[0]-1:myvals[0]+1],layer.Buf.Pbf[myvals[1]-2:myvals[1]+1]
	value_bytes := []byte{}
	key_bytes := []byte{}
	keymap := map[string]uint32{}
	valuemap := map[interface{}]uint32{}
	if (size_key!=0&&size_val!=0) {
		value_bytes = layer.Buf.Pbf[myvals[0]-1:myvals[1]-1]
		key_bytes = layer.Buf.Pbf[myvals2[0]-1:myvals2[1]-1]

		// fmt.Println(value_bytes[0],value_bytes[len(value_bytes)-1])

		for pos,key := range layer.Keys {
			keymap[key] = uint32(pos) 
		}
		for pos,val := range layer.Values {
			valuemap[val] = uint32(pos)
		}
	}
	// creating cursor 
	cur := NewCursorExtent(tileid,4326)
	
	// getting the last feature 
	var feat_bytes []byte
	layer.feature_position = layer.Number_Features-1
	if layer.Number_Features>0 {
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
		if len(layer.features) > 0 {
			start_pos := layer.features[0] 
			
			layer.Buf.Pos = layer.features[len(layer.features)-1]
			//layer.Buf.Pos = layer.features[len(layer.features)-1]
			//fmt.Println(layer.Buf.Pos,layer.Buf.Pos,layer.Buf.Pbf[layer.Buf.Pos-3:layer.Buf.Pos+25],layer.Buf.Pbf[layer.Buf.Pos])
	
			end_pos := layer.Buf.Pos + int(layer.Buf.ReadVarint())
			feat_bytes = layer.Buf.Pbf[start_pos-1:end_pos]
			//fmt.Println(feat_bytes)
			//fmt.Println(start_pos,layer.Buf.Pos,layer.Buf.Pbf[start_pos-3:start_pos+3],layer.Buf.Pbf[start_pos])
		} else {
			feat_bytes = []byte{}
		}
	}

	bds := m.Bounds(tileid)

	layerwrite := &LayerWrite{
		Name:layer.Name,
		Extent:layer.Extent,
		Version:layer.Version,
		TileID:tileid,
		Features: feat_bytes,
		Cursor:cur,
		DeltaX: bds.E - bds.W,
		DeltaY: bds.N - bds.S,
		Keys_Bytes: []byte{},
		Values_Bytes: []byte{},
		Values_Map: map[interface{}]uint32{},
		Keys_Map: map[string]uint32{},
	}


	// creating the keys map
	//keymap := map[string]uint32{}

	// for _,key := range layer.Keys {
	// 	//keymap[key] = uint32(pos)
	// 	layerwrite.AddKey(key)
	// }
	
	// // creeating values map
	// //valuemap := map[interface{}]uint32{}
	// for _,value := range layer.Values {
	// 	//valuemap[value] = uint32(pos)
	// 	layerwrite.AddValue(value)
	// }
	
	if (size_key!=0&&size_val!=0) {
		layerwrite.Values_Bytes = value_bytes
		layerwrite.Keys_Bytes = key_bytes
		layerwrite.Keys_Map = keymap
		layerwrite.Values_Map = valuemap 
	}

	// fmt.Println(layerwrite.Values_Bytes[0],value_bytes[0],layerwrite.Values_Bytes[len(layerwrite.Values_Bytes)-1],value_bytes[len(value_bytes)-1])
	// for pos := range layerwrite.Values_Bytes {
	// 	fmt.Println(layerwrite.Values_Bytes[pos],value_bytes[pos],len(layerwrite.Values_Bytes),len(value_bytes))
	// }
	// fmt.Println(layerwrite.Keys_Bytes[0],key_bytes[0],layerwrite.Keys_Bytes[len(layerwrite.Keys_Bytes)-1],key_bytes[len(key_bytes)-1])
	return layerwrite,nil
}


// filters layer 
// a typical filterfunc might look something like this:
// func filterfunc(layer *Layer,tileid m.TileID) []int {
// 	poss := []int{}
// 	for layer.Next() {
// 		feat,_ := layer.Feature()
// 		fmt.Println(feat.FeaturePos)
// 		if feat.FeaturePos%8==0 {
// 			poss = append(poss,feat.FeaturePos)
// 		}
// 		fmt.Println(feat)
// 	}
// 	return poss
// }
func (layer *Layer) FilterLayer(tileid m.TileID,filterfunc func(lay *Layer,tileid m.TileID) []int) ([]byte,error) {
	layer2 := *layer
	layw,err := layer2.ToLayerWrite(tileid)
	if err != nil {
		return []byte{},err
	}

	// resetting the entire layer write context
	layw.Values_Bytes = []byte{}
	layw.Values_Map = map[interface{}]uint32{}
	layw.Keys_Bytes = []byte{}
	layw.Keys_Map = map[string]uint32{}
	layw.Features = []byte{}

	// getting layer positions from the filter func
	positions := filterfunc(layer,tileid)
	
	// writing the positions into the tile lazily 
	layw.AddFeaturesLazy(positions,layer)

	return layw.Flush(),nil
}