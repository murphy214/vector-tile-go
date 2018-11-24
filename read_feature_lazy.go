package vt

import (
	//"math"
	"errors"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/pbf"
	"github.com/paulmach/go.geojson"
	"math"
)

type Feature struct {
	ID           int
	Type         string
	Properties   map[string]interface{}
	geometry_pos int
	extent       int
	geom_int     int
	
	Buf          *pbf.PBF
}

func DeltaDim(num int) float64 {
	if num%2 == 1 {
		return float64((num + 1) / -2)
	} else {
		return float64(num / 2)
	}
	return float64(0)
}

// signed area frunction
func SignedArea(ring [][]float64) float64 {
	sum := 0.0
	i := 0
	lenn := len(ring)
	j := lenn - 1
	var p1, p2 []float64

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
func Project(line [][]float64, x0 float64, y0 float64, size float64) [][]float64 {
	for j := range line {
		p := line[j]
		y2 := 180.0 - (float64(p[1])+y0)*360.0/size
		line[j] = []float64{
			(float64(p[0])+x0)*360.0/size - 180.0,
			360.0/math.Pi*math.Atan(math.Exp(y2*math.Pi/180.0)) - 90.0}
	}
	return line
}

// reads a feature lazily
func (layer *Layer) Feature() (feature *Feature, err error) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("Error in Feature()")
			layer.feature_position++
		}
	}()

	layer.Buf.Pos = layer.features[layer.feature_position]
	endpos := layer.Buf.Pos + layer.Buf.ReadVarint()
	//startpos := layer.Buf.Pos
	feature = &Feature{Properties: map[string]interface{}{}}

	for layer.Buf.Pos < endpos {
		key, val := layer.Buf.ReadKey()

		// logic for handlign id
		if key == 1 && val == 0 {
			feature.ID = int(layer.Buf.ReadUInt64())
		}
		// logic for handling tags
		if key == 2 && val == 2 {
			//fmt.Println(feature)
			tags := layer.Buf.ReadPackedUInt32()
			i := 0
			for i < len(tags) {
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
			feature.geom_int = geom_type
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
	return feature, err
}

func (feature *Feature) LoadGeometry() (geomm *geojson.Geometry, err error) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("Error in feature.LoadGeometry()")
		}
	}()

	// getting geometry
	// this huge code block is to reduce allocations and shit

	feature.Buf.Pos = feature.geometry_pos
	geom := feature.Buf.ReadPackedUInt32()
	//fmt.Println(geom)
	pos := 0
	var lines [][][]float64
	var polygons [][][][]float64
	var firstpt []float64
	geom_type := feature.geom_int

	for pos < len(geom) {
		if geom[pos] == 9 {
			pos += 1
			if pos != 1 && geom_type == 2 {
				firstpt = []float64{firstpt[0] + DeltaDim(int(geom[pos])), firstpt[1] + DeltaDim(int(geom[pos+1]))}
			} else {
				firstpt = []float64{DeltaDim(int(geom[pos])), DeltaDim(int(geom[pos+1]))}
			}
			pos += 2
			if len(geom) == 3 {
				lines = [][][]float64{{firstpt}}
			}
			if pos < len(geom) {
				//fmt.Println(geom[pos])
				cmdLen := geom[pos]
				length := int(cmdLen >> 3)
				//fmt.Println(length)
				line := make([][]float64, length+1)
				pos += 1
				endpos := pos + length*2
				line[0] = firstpt
				i := 1
				for pos < endpos && pos+1 < len(geom) {
					firstpt = []float64{firstpt[0] + DeltaDim(int(geom[pos])), firstpt[1] + DeltaDim(int(geom[pos+1]))}
					line[i] = firstpt
					i++
					pos += 2
					//fmt.Println(pos)
				}
				lines = append(lines, line[:i])
				line = [][]float64{firstpt}

			} else {
				//line := [][]float64{firstpt}
				//lines = append(lines, line)
				pos += 1
			}

		} else if pos < len(geom) {
			if geom[pos] == 15 {
				//polygons = append(polygons, lines)
				//lines = [][][]float64{}
				pos += 1
			} else {
				pos += 1
			}
		} else {
			pos += 1
		}
	}
	if geom_type == 3 {
		for pos, line := range lines {
			f, l := line[0], line[len(line)-1]
			if !(f[0] == l[0] && l[1] == f[1]) {
				line = append(line, line[0])
			}
			lines[pos] = line
		}

		if len(lines) == 1 {
			polygons = append(polygons, lines)
		} else {
			for _, line := range lines {
				if len(line) > 0 {
					val := SignedArea(line)
					if val < 0 {
						polygons = append(polygons, [][][]float64{line})
					} else {
						if len(polygons) == 0 {
							polygons = append(polygons, [][][]float64{line})

						} else {
							polygons[len(polygons)-1] = append(polygons[len(polygons)-1], line)

						}
					}
				}
			}

		}
	} else {
		polygons = append(polygons, lines)
	}

	switch geom_type {
	case 1:
		if len(polygons[0][0]) == 1 {
			geomm = geojson.NewPointGeometry(polygons[0][0][0])
		} else {
			geomm = geojson.NewMultiPointGeometry(polygons[0][0]...)

		}
	case 2:
		if len(polygons[0]) == 1 {
			geomm = geojson.NewLineStringGeometry(polygons[0][0])
		} else {
			geomm = geojson.NewMultiLineStringGeometry(polygons[0]...)

		}
	case 3:
		if len(polygons) == 1 {
			geomm = geojson.NewPolygonGeometry(polygons[0])
		} else {
			geomm = geojson.NewMultiPolygonGeometry(polygons...)

		}
	}

	return geomm, err
}

// loads a geojson feature from a lazy feature
func (feature *Feature) ToGeoJSON(tile m.TileID) (*geojson.Feature, error) {
	var err error
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("Error in feature.ToGeoJSON()")
		}
	}()
	// this values will be used to preproject the coordinates
	extent := feature.extent
	size := float64(extent) * float64(math.Pow(2, float64(tile.Z)))
	x0 := float64(extent) * float64(tile.X)
	y0 := float64(extent) * float64(tile.Y)
	geometry, err := feature.LoadGeometry()
	if err != nil {
		return &geojson.Feature{}, err
	}
	switch geometry.Type {
	case "Point":
		geometry.Point = Project([][]float64{geometry.Point}, x0, y0, size)[0]
	case "MultiPoint":
		geometry.MultiPoint = Project(geometry.MultiPoint, x0, y0, size)
	case "LineString":
		geometry.LineString = Project(geometry.LineString, x0, y0, size)
	case "MultiLineString":
		for i := range geometry.MultiLineString {
			geometry.MultiLineString[i] = Project(geometry.MultiLineString[i], x0, y0, size)
		}
	case "Polygon":
		for i := range geometry.Polygon {
			geometry.Polygon[i] = Project(geometry.Polygon[i], x0, y0, size)
		}
	case "MultiPolygon":
		for i := range geometry.MultiPolygon {
			for j := range geometry.MultiPolygon[i] {
				geometry.MultiPolygon[i][j] = Project(geometry.MultiPolygon[i][j], x0, y0, size)
			}
		}
	}

	new_feature := geojson.NewFeature(geometry)
	new_feature.Properties = feature.Properties
	new_feature.ID = feature.ID

	return new_feature, err
}

// converts a single geometry
func convertpt(pt []float64, dim float64) []float64 {
	if pt[0] < 0 {
		//pt[0] = 0
	}
	if pt[1] < 0 {
		//pt[1] = 0
	}
	return []float64{pbf.Round(pt[0]/dim, .5, 0), pbf.Round(pt[1]/dim, .5, 0)}
}

// converts the line
func convertln(ln [][]float64, dim float64) [][]float64 {
	for i := range ln {
		ln[i] = convertpt(ln[i], dim)
	}
	return ln
}

// convert lines
func convertlns(lns [][][]float64, dim float64) [][][]float64 {
	for i := range lns {
		lns[i] = convertln(lns[i], dim)
	}
	return lns
}

func ConvertGeometry(geom *geojson.Geometry, dimf float64) *geojson.Geometry {
	if geom == nil {
		return &geojson.Geometry{}
	}
	switch geom.Type {
	case "Point":
		geom.Point = convertpt(geom.Point, dimf)
	case "MultiPoint":
		geom.MultiPoint = convertln(geom.MultiPoint, dimf)
	case "LineString":
		geom.LineString = convertln(geom.LineString, dimf)
	case "MultiLineString":
		geom.MultiLineString = convertlns(geom.MultiLineString, dimf)
	case "Polygon":
		geom.Polygon = convertlns(geom.Polygon, dimf)
	case "MultiPolygon":
		for i := range geom.MultiPolygon {
			geom.MultiPolygon[i] = convertlns(geom.MultiPolygon[i], dimf)
		}
	}
	return geom
}

// loads the given geometry but scaled to the given dimennsioin
func (feature *Feature) LoadGeometryScaled(dim float64) (geomm *geojson.Geometry, err error) {
	geom, err := feature.LoadGeometry()
	geom2 := ConvertGeometry(geom, dim)
	return geom2, err
}
