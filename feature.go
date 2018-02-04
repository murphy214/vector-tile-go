package pbf

import (
	//"fmt"
	"math"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
)

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
func LayerFeature(feat_bytes *PBF,end int,keys []string,values []interface{}) *Feature {
	// setting up feature
	feat := &Feature{Properties:map[string]interface{}{}}

	for feat_bytes.Pos < end {
		key,val := feat_bytes.ReadKey()
		// logic for handlign id
		if key == 1 && val == 0 {
			feat.Id = feat_bytes.ReadUInt64()
		}
		// logic for handling tags
		if key == 2 && val == 2 {
			tags := feat_bytes.ReadPackedUInt32()
            //fmt.Println(len(tags),len(values),len(keys),"dasdfa")
			i := 0
			for i < len(tags) {
                //fmt.Println(i,tags[i],tags[i+1])
                key := keys[tags[i]]
                var val interface{}
                if len(values) <= int(tags[i+1]) {
                    val = ""
                } else {
                    val = values[tags[i+1]]
                }
				feat.Properties[key] = val
				i += 2
			}
		}
		// logic for handling features
		if key == 3 && val == 0 {
			feat.Type = int(feat_bytes.Varint()[0])
		}	
		// logic for handling geometry
		if key == 4 && val == 2 {
			feat.Geometry = feat_bytes.Pos
			size := feat_bytes.ReadVarint()
			feat_bytes.Pos += size
		}	
	}
	feat.Buf = feat_bytes
	return feat
}

// loads a given geometry 
func (feature *Feature) LoadGeometry() [][][]float64 {
    feature.Buf.Pos = feature.Geometry

	end := feature.Buf.ReadVarint() + feature.Buf.Pos
    cmd,length,x,y := 0,0,0.0,0.0
  	lines := [][][]float64{}
  	line := [][]float64{}

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

            if (cmd == 1) { // moveTo
          		lines = append(lines,line)
                line = [][]float64{}
            }
            //if len(line)
            
            line = append(line,[]float64{x,y})
            //line.push(new Point(x, y));

        } else if (cmd == 7) {

            // Workaround for https://github.com/mapbox/mapnik-vector-tile/issues/90
            if len(line) > 0 {
                line = append(line,line[0]) // closePolygon
            }

        }
    }

    if len(line) > 0 {
    	lines = append(lines,line)
    }
    return lines[1:]
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

// this function (hopefully) classifies the rings correctly
func classifyRings(rings [][][]float64) [][][][]float64 {
    /*
    if (len <= 1) {
    	return [rings];
    }
	*/
	polygons := [][][][]float64{}
  	polygon := [][][]float64{}

    for i := range rings {
		area := SignedArea(rings[i])
        if (area == 0) {

        }
        if (area < 0) {
            if len(polygon) > 0 {
            	polygons = append(polygons,polygon)
            }
            polygon = [][][]float64{rings[i]}

        } else {
            polygon = append(polygon,rings[i])
        }
    }
    if len(polygon) > 0 {
		polygons = append(polygons,polygon)
    } 

    return polygons
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

// function that converts a single feature to a geojson
func (feature *Feature) ToGeoJSON(tileid m.TileID) *geojson.Feature {
	size := float64(feature.Extent) * float64(math.Pow(2, float64(tileid.Z)))
   	x0 := float64(feature.Extent * int(tileid.X))
    y0 := float64(feature.Extent * int(tileid.Y))
    coords := feature.LoadGeometry()

    // getting the geometry
	var geom *geojson.Geometry
    switch (feature.Type) {
    case 1:
        points := [][]float64{}
        for i := range coords {
            points = append(points,coords[i][0])
        }
        //coords = points;
        points = Project(points,x0,y0,size)
        if len(points) == 1 {
        	geom = geojson.NewPointGeometry(points[0])
        } else {
        	geom = geojson.NewMultiPointGeometry(points[0])
        }

    case 2:
        for i := range coords {
            coords[i] = Project(coords[i],x0,y0,size)
        }
        if len(coords) == 1 {
        	geom = geojson.NewLineStringGeometry(coords[0])
        } else {
        	geom = geojson.NewMultiLineStringGeometry(coords...)
        }

    case 3:
        newcoords := classifyRings(coords);
        for i := range newcoords {
            for j := range newcoords[i] {
                newcoords[i][j] = Project(newcoords[i][j],x0,y0,size)
            }
        }
        if len(newcoords) == 1 {
        	geom = geojson.NewPolygonGeometry(newcoords[0])
        } else {
        	geom = geojson.NewMultiPolygonGeometry(newcoords...)
        }
    }

    // returnign the entire feature
    if feature.Id != 0 {
    	return &geojson.Feature{ID:feature.Id,Geometry:geom,Properties:feature.Properties}
    } else {
    	return &geojson.Feature{Geometry:geom,Properties:feature.Properties} 
    }
}




