package vt

import (
    //"math"
    "github.com/paulmach/go.geojson"
     m "github.com/murphy214/mercantile"
     "math"
    "github.com/murphy214/pbf"

)

type Feature struct {
    ID int
    Type string
    Properties map[string]interface{}
    geometry_pos int
    extent int
    Buf *pbf.PBF
}

func (layer *Layer) Feature() *Feature {

    layer.Buf.Pos = layer.features[layer.feature_position]
    endpos := layer.Buf.Pos + layer.Buf.ReadVarint()
    //startpos := layer.Buf.Pos
    feature := &Feature{Properties:map[string]interface{}{}}    

    for layer.Buf.Pos < endpos {
        key,val := layer.Buf.ReadKey()

        // logic for handlign id
        if key == 1 && val == 0 {
            feature.ID = int(layer.Buf.ReadUInt64())
        }
        // logic for handling tags
        if key == 2 && val == 2 {
            //fmt.Println(feature)
            tags := layer.Buf.ReadPackedUInt32()
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
            geom_type := int(layer.Buf.Varint()[0])
            switch geom_type {
            case 1:
                feature.Type = "Point"
            case 2:
                feature.Type = "LineString"
            case 3:
                feature.Type = "Polygon"
            }
        }   
        // logic for handling geometry
        if key == 4 && val == 2 {
            feature.geometry_pos = layer.Buf.Pos
            size := layer.Buf.ReadVarint()
            layer.Buf.Pos += size + 1
        }
    }
    feature.extent = layer.Extent
    feature.Buf = layer.Buf
    layer.feature_position += 1
    return feature 
}

func (feature *Feature) LoadGeometry() *geojson.Geometry {
    // getting geometry
    // this huge code block is to reduce allocations and shit
     switch feature.Type {
     case "Point":

        feature.Buf.Pos = feature.geometry_pos
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
                
                line = append(line,pt)
                //line.push(new Point(x, y));
            }
            if length < 0 && feature.Buf.Pos + 1 == end {
                feature.Buf.Pos += 1
            }
        }


        if len(line) == 1 {
            return geojson.NewPointGeometry(line[0])
        } else {
            return geojson.NewMultiPointGeometry(line...)
        }
    case "LineString":
        feature.Buf.Pos = feature.geometry_pos

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
        }
        if len(line) > 0 {
            lines = append(lines,line)
        }

        if len(lines) == 1 {
            return geojson.NewLineStringGeometry(lines[0])
        } else {
            return geojson.NewMultiLineStringGeometry(lines...)
        }
    case "Polygon":
        feature.Buf.Pos = feature.geometry_pos

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
                //fmt.Println(SignedArea(line))
                //newline = append(newline,newline[0])
                if SignedArea(line) < 0 {
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
        }
        if len(line) > 0 {
            polygons = append(polygons,[][][]float64{line})
        }

        if len(polygons) == 1 {
            return geojson.NewPolygonGeometry(polygons[0])
        } else if len(polygons) > 0 {
            return geojson.NewMultiPolygonGeometry(polygons...)
        }    
    }   
    return &geojson.Geometry{}
}

// loads a geojson feature from a lazy feature
func (feature *Feature) ToGeoJSON(tile m.TileID) *geojson.Feature {
    // this values will be used to preproject the coordinates
    extent := feature.extent
    size := float64(extent) * float64(math.Pow(2, float64(tile.Z)))
    x0 := float64(extent * int(tile.X))
    y0 := float64(extent * int(tile.Y))   
    geometry := feature.LoadGeometry()

    switch geometry.Type {
    case "Point":
        geometry.Point = Project([][]float64{geometry.Point},x0,y0,size)[0]
    case "MultiPoint":
        geometry.MultiPoint = Project(geometry.MultiPoint,x0,y0,size)
    case "LineString":
        geometry.LineString = Project(geometry.LineString,x0,y0,size)
    case "MultiLineString":
        for i := range geometry.MultiLineString {
            geometry.MultiLineString[i] = Project(geometry.MultiLineString[i],x0,y0,size)
        }
    case "Polygon":
        for i := range geometry.Polygon {
            geometry.Polygon[i] = Project(geometry.Polygon[i],x0,y0,size)
        }
    case "MultiPolygon":
         for i := range geometry.MultiPolygon {
            for j := range geometry.MultiPolygon[i] {
                geometry.MultiPolygon[i][j] = Project(geometry.MultiPolygon[i][j],x0,y0,size)
            }
        }       
    }

    new_feature := geojson.NewFeature(geometry)
    new_feature.Properties = feature.Properties
    new_feature.ID = feature.ID

    return new_feature
}
