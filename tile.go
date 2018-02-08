package pbf 


import (
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"

	"fmt"
)

// fuck debugging
func a() {
	fmt.Println()
}

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
			size := pbfval.ReadVarint()
			//fmt.Println(size)
			if size != 0 {
				//fmt.Println(bytevals[pbfval.Pos:pbfval.Pos+20])
				//fmt.Println(bytevals[:20])
				layer := &Layer{
					Buf:&PBF{
						Pbf:bytevals[pbfval.Pos:pbfval.Pos+size],
						Length:size,
						},
				}
				layer.New_Layer()
				//layer.Buf.Pos = 
				//fmt.Println(layer.Buf.Pos,pbfval.Pos,"here")
				vt[layer.Name] = layer
				pbfval.Pos += size
			}

		}	
	}
	return vt
}


func (vt Vector_Tile) ToGeoJSON(tileid m.TileID) map[string][]*geojson.Feature {
	totalmap := map[string][]*geojson.Feature{}
	// going through each layer
	for k,v := range vt {
		totalmap[k] = make([]*geojson.Feature,v.Number_Features)
		if v.Number_Features > 10000000 {
			c := make(chan *geojson.Feature)
			guard := make(chan struct{}, 1000)
			i := 0
			for i < v.Number_Features {
				guard <- struct{}{}
				go func(i int,c chan *geojson.Feature) {
					<-guard
					c <- v.Feature(i,tileid)
				}(i,c)

				//totalmap[k][i] = v.Feature(i).ToGeoJSON(tileid)
				//fmt.Println(i,v.Number_Features)
				i++
			}
		
			ii := 0 
			for ii < v.Number_Features {
				totalmap[k][ii] = <-c
				ii++
			}
		} else {
			i := 0 
			for i < v.Number_Features {
				totalmap[k][i] = v.Feature(i,tileid)

				i++
			}			
		}

	}
	return totalmap
}

func ToGeoJSON(bytevals []byte,tileid m.TileID) map[string][]*geojson.Feature {
	var bytevals2 = bytevals
	var vt = New_Vector_Tile(bytevals2)

	totalmap := map[string][]*geojson.Feature{}
	// going through each layer
	for k,v := range vt {
		totalmap[k] = make([]*geojson.Feature,v.Number_Features)

		if v.Number_Features > 100000000 {
			c := make(chan *geojson.Feature)
			guard := make(chan struct{}, 1000)
			i := 0
			for i < v.Number_Features {
				guard <- struct{}{}
				go func(i int,c chan *geojson.Feature) {
					<-guard
					c <- v.Feature(i,tileid)
				}(i,c)

				//totalmap[k][i] = v.Feature(i).ToGeoJSON(tileid)
				//fmt.Println(i,v.Number_Features)
				i++
			}


			ii := 0 
			for ii < v.Number_Features {
				totalmap[k][ii] = <-c
				ii++
			}
		} else {
			i := 0 
			for i < v.Number_Features {
				totalmap[k][i] = v.Feature(i,tileid)		
				i++
			}			
		}

	}
	return totalmap
}






