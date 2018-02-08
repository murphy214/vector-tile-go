package pbf

import (
	"fmt"
	"math"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
)

// because fuck debuggin
func c() {
    fmt.Println()
}

// a single feature
type Feature struct {
	Id uint64 // the id
	Properties map[string]interface{} // the properties associated with the layer
	Type int	
	Geometry int // a single byte indicating where a geometry is located 
	Extent int // extent of feature
	Buf *PBF 
}


// function for getting a feature at a given layer position
func (feature *Feature) LayerFeature_Raw(keys []string,values []interface{},startpos int) {
    // setting up feature

    for feature.Buf.Pos < feature.Buf.Length {
        key,val := feature.Buf.ReadKey()

        // logic for handlign id
        if key == 1 && val == 0 {
            feature.Id = feature.Buf.ReadUInt64()
        }
        // logic for handling tags
        if key == 2 && val == 2 {
            tags := feature.Buf.ReadPackedUInt32()
            //fmt.Println(len(tags),len(values),len(keys),"dasdfa")
            i := 0
            for i < len(tags) {
                //fmt.Println(tags,keys,tags[i],tags[i+1])
                var key string
                if len(keys) <= int(tags[i]) {
                    key = ""
                } else {
                    key = keys[tags[i]]
                }
                var val interface{}
                if len(values) <= int(tags[i+1]) {
                    val = ""
                } else {
                    val = values[tags[i+1]]
                }
                feature.Properties[key] = val
                i += 2
            }
        }
        // logic for handling features
        if key == 3 && val == 0 {
            feature.Type = int(feature.Buf.Varint()[0])
        }   
        // logic for handling geometry
        if key == 4 && val == 2 {
            feature.Geometry = feature.Buf.Pos
            size := feature.Buf.ReadVarint()
            feature.Buf.Pos += size + 1
        }

    }

}


func (layer Layer) Feature(pos int,tileid m.TileID) *geojson.Feature {
    layer.Buf.Pos = layer.features[pos]
    endpos := layer.Buf.Pos + layer.Buf.ReadVarint()
    startpos := layer.Buf.Pos
    //&PBF{PBF:layer.Buf.Pbf[startpos:endpos],Length:endpos-startpos}
    
    feature := Feature{
        Buf:
            &PBF{
                Pbf:layer.Buf.Pbf[startpos:endpos],
                Length:(endpos-startpos)+1,
            },
        Properties:map[string]interface{}{},
    }
    // setting up feature

    for feature.Buf.Pos < feature.Buf.Length {
        //fmt.Println(len(feature.Buf.Pbf))
        key,val := feature.Buf.ReadKey()

        // logic for handlign id
        if key == 1 && val == 0 {
            feature.Id = feature.Buf.ReadUInt64()
        }
        // logic for handling tags
        if key == 2 && val == 2 {
            //fmt.Println(feature)
            tags := feature.Buf.ReadPackedUInt32()
            //fmt.Println(len(tags),len(values),len(keys),"dasdfa")
            i := 0
            for i < len(tags) {
                //fmt.Println(tags,keys,tags[i],tags[i+1])
                var key string
                if len(layer.Keys) <= int(tags[i]) {
                    key = ""
                } else {
                    key = layer.Keys[tags[i]]
                }
                var val interface{}
                if len(layer.Values) <= int(tags[i+1]) {
                    val = ""
                } else {
                    val = layer.Values[tags[i+1]]
                }
                feature.Properties[key] = val
                i += 2
            }
        }
        // logic for handling features
        if key == 3 && val == 0 {
            feature.Type = int(feature.Buf.Varint()[0])
        }   
        // logic for handling geometry
        if key == 4 && val == 2 {
            feature.Geometry = feature.Buf.Pos
            size := feature.Buf.ReadVarint()
            feature.Buf.Pos += size + 1
        }
        //fmt.Println(key,val,feat_bytes.Pos,end) 

    }
    var geometry *geojson.Geometry
    // getting geometry
    // this huge code block is to reduce allocations and shit
     if feature.Type == 1 {
        size := float64(feature.Extent) * float64(math.Pow(2, float64(tileid.Z)))
        x0 := float64(feature.Extent * int(tileid.X))
        y0 := float64(feature.Extent * int(tileid.Y))   

        feature.Buf.Pos = feature.Geometry

        end := feature.Buf.ReadVarint() + feature.Buf.Pos
        cmd,length,x,y := 0,0,0.0,0.0
        line := [][]float64{}
        var pt []float64
        var cmdLen int
        for feature.Buf.Pos < end {
            if length == 0 {
                cmdLen = feature.Buf.ReadVarint();
                cmd = cmdLen & 0x7
                length = cmdLen >> 3
            }
            length--

            if (cmd == 1 || cmd == 2) {
                x += feature.Buf.ReadSVarint()
                y += feature.Buf.ReadSVarint()
                pt = []float64{x,y}
                //if (cmd == 1) && len(line) > 0 { // moveTo
                 //   line = [][]float64{}
                //}
                //if len(line)
                
                line = append(line,pt)
                //line.push(new Point(x, y));
            }
            if length < 0 && feature.Buf.Pos + 1 == end {
                feature.Buf.Pos += 1
            }
        }


        line = Project(line,x0,y0,size)

        if len(line) == 1 {
            geometry =  geojson.NewPointGeometry(line[0])
        } else {
            geometry = geojson.NewMultiPointGeometry(line...)
        }
    } else if feature.Type == 2 {
        size := float64(feature.Extent) * float64(math.Pow(2, float64(tileid.Z)))
        x0 := float64(feature.Extent * int(tileid.X))
        y0 := float64(feature.Extent * int(tileid.Y))   

        feature.Buf.Pos = feature.Geometry

        end := feature.Buf.ReadVarint() + feature.Buf.Pos
        cmd,length,x,y := 0,0,0.0,0.0
        line := [][]float64{}
        lines := [][][]float64{}
        var pt []float64

        var cmdLen int
        for feature.Buf.Pos < end {
            if length == 0 {
                cmdLen = feature.Buf.ReadVarint();
                cmd = cmdLen & 0x7
                length = cmdLen >> 3
            }
            length--

            if (cmd == 1 || cmd == 2) {
                x += feature.Buf.ReadSVarint()
                y += feature.Buf.ReadSVarint()
                pt = []float64{x,y}

                if (cmd == 1) && len(line) > 0 { // moveTo
                   lines = append(lines,line)

                   line = [][]float64{}
                }
                //if len(line)
                
                line = append(line,pt)
                //line.push(new Point(x, y));
            }
            //fmt.Println(length,feature.Buf.Pos,end)
            if length < 0 && feature.Buf.Pos + 1 == end {
                feature.Buf.Pos += 1
            }
        }
        if len(line) > 0 {
            lines = append(lines,line)
        }

        for i := range lines {
            lines[i] = Project(lines[i],x0,y0,size)
        }
        if len(lines) == 1 {
            geometry = geojson.NewLineStringGeometry(lines[0])
        } else {
            geometry = geojson.NewMultiLineStringGeometry(lines...)
        }
    } else if feature.Type == 3 {
        size := float64(feature.Extent) * float64(math.Pow(2, float64(tileid.Z)))
        x0 := float64(feature.Extent * int(tileid.X))
        y0 := float64(feature.Extent * int(tileid.Y))   

        feature.Buf.Pos = feature.Geometry

        end := feature.Buf.ReadVarint() + feature.Buf.Pos
        cmd,length,x,y := 0,0,0.0,0.0
        line := [][]float64{}
        polygons := [][][][]float64{} 
        var pt []float64

        var cmdLen int
        for feature.Buf.Pos < end {
            if length <= 0 {
                cmdLen = feature.Buf.ReadVarint();
                cmd = cmdLen & 0x7
                length = cmdLen >> 3
                //fmt.Println(cmdLen)
            }
            length--

            if (cmd == 1 || cmd == 2) {
                x += feature.Buf.ReadSVarint()
                y += feature.Buf.ReadSVarint()
                pt = []float64{x,y}
                if (cmd == 1) && len(line) > 0 { // moveTo
                   line = [][]float64{}
                }
                //if len(line)
                
                line = append(line,pt)
                //line.push(new Point(x, y));
            } else if (cmd == 7) {

                //newline = append(newline,newline[0])
                if SignedArea(line) > 0 {
                    polygons = append(polygons, [][][]float64{line})
                    //newline = [][]int{}
                } else {
                    if len(polygons) == 0 {
                        polygons = append(polygons, [][][]float64{line})

                    } else {
                        polygons[len(polygons)-1] = append(polygons[len(polygons)-1], line)

                    }
                    line = [][]float64{}
                }

            }
            //fmt.Println(length,polygons,feature.Buf.Pbf[feature.Buf.Pos:],end)
            //if length < -100 && feature.Buf.Pos + 1 == end {
              //  feature.Buf.Pos += 1
            //}    
        }
        //fmt.Println(polygons)



        for i := range polygons {
            for j := range polygons[i] {
                polygons[i][j] = Project(polygons[i][j],x0,y0,size)
            }
        }
        if len(polygons) == 1 {
            geometry = geojson.NewPolygonGeometry(polygons[0])
        } else {
            geometry = geojson.NewMultiPolygonGeometry(polygons...)
        }    
    }   





    if feature.Id != 0 {
        return &geojson.Feature{ID:feature.Id,Geometry:geometry,Properties:feature.Properties}
    } else {
        return &geojson.Feature{Geometry:geometry,Properties:feature.Properties}
    }    
}




// signed area frunction
func SignedArea(ring [][]float64) float64 {
    sum := 0.0
    i := 0
    lenn := len(ring)
    j := lenn - 1 
    var p1,p2 []float64
   
    for i < lenn {
        if i != 0 {
            j = i - 1
        }
        p1 = ring[i]
        p2 = ring[j]
        sum += (p2[0] - p1[0]) * (p1[1] + p2[1])
        i++
    }
    return sum
}

// this function projects a single line
func Project(line [][]float64,x0 float64,y0 float64,size float64) [][]float64 {
    for j := range line {
        p := line[j]
        y2 := 180 - (p[1] + y0) * 360.0 / size 
        line[j] = []float64{
            (p[0] + x0) * 360.0 / size - 180.0,
            360.0 / math.Pi * math.Atan(math.Exp(y2 * math.Pi / 180.0)) - 90.0}
    }
    return line
}