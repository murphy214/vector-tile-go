package pbf 


import (
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"

	//"fmt"
)

// upper vector tile structure
type Vector_Tile map[string]*Layer

// create / reads a new vector tile from a byte array 
func New_Vector_Tile(bytevals []byte) Vector_Tile {
	pbfval := &PBF{Pbf:bytevals,Length:len(bytevals)}
	//fmt.Println(pbfval.ReadVarint())
	//fmt.Println(pbfval.ReadKey())
	vt := Vector_Tile{}
	for pbfval.Pos < pbfval.Length {
		key,val := pbfval.ReadKey()
		if key == 3 && val == 2 {
			pbfval.ReadVarint()
			layer := New_Layer(pbfval)
			vt[layer.Name] = layer
		}	
	}
	return vt
}


func (vt Vector_Tile) ToGeoJSON(tileid m.TileID) map[string][]*geojson.Feature {
	totalmap := map[string][]*geojson.Feature{}
	// going through each layer
	for k,v := range vt {
		newlist := make([]*geojson.Feature,v.Number_Features)
		i := 0
		for i < v.Number_Features {
			newlist[i] = v.Feature(i).ToGeoJSON(tileid)
			i++
		}
		totalmap[k] = newlist
	}
	return totalmap
}






