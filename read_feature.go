package vt

import (
	"fmt"
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"math"
)

func (tile *Tile) Feature(layername string, keys []string, values []interface{}, extent int, pos int) *geojson.Feature {
	tile.Buf.Pos = pos
	endpos := tile.Buf.Pos + tile.Buf.ReadVarint()
	//startpos := tile.Buf.Pos

	feature := &geojson.Feature{Properties: map[string]interface{}{}}
	var feature_geometry, id, geom_type int
	if extent == 0 {
		extent = 4096
	}
	for tile.Buf.Pos < endpos {
		key, val := tile.Buf.ReadKey()

		// logic for handlign id
		if key == 1 && val == 0 {
			id = int(tile.Buf.ReadUInt64())
		}
		// logic for handling tags
		if key == 2 && val == 2 {
			//fmt.Println(feature)
			tags := tile.Buf.ReadPackedUInt32()
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
			geom_type = int(tile.Buf.Varint()[0])
		}
		// logic for handling geometry
		if key == 4 && val == 2 {
			feature_geometry = tile.Buf.Pos
			size := tile.Buf.ReadVarint()
			tile.Buf.Pos += size + 1
		}
	}
	feature.Geometry = Geometry(tile, feature_geometry, geom_type, tile.TileID, extent)
	/*
			// getting geometry
			// this huge code block is to reduce allocations and shit
			if geom_type == 1 {
				size := float64(extent) * float64(math.Pow(2, float64(tile.TileID.Z)))
				x0 := float64(extent * int(tile.TileID.X))
				y0 := float64(extent * int(tile.TileID.Y))

				tile.Buf.Pos = feature_geometry

				end := tile.Buf.ReadVarint() + tile.Buf.Pos

				cmd, length, x, y := 0, 0, 0.0, 0.0
				line := [][]float64{}
				var pt []float64
				var cmdLen int
				for tile.Buf.Pos < end {
					if length == 0 {
						cmdLen = tile.Buf.ReadVarint()
						cmd = cmdLen & 0x7
						length = cmdLen >> 3
					}
					length--

					if cmd == 1 || cmd == 2 {
						x += tile.Buf.ReadSVarint()
						y += tile.Buf.ReadSVarint()
						pt = []float64{x, y}
						//if (cmd == 1) && len(line) > 0 { // moveTo
						//   line = [][]float64{}
						//}
						//if len(line)

						line = append(line, pt)
						//line.push(new Point(x, y));
					}
					if length < 0 && tile.Buf.Pos+1 == end {
						tile.Buf.Pos += 1
					}
				}

				line = Project(line, x0, y0, size)

				if len(line) == 1 {
					feature.Geometry = geojson.NewPointGeometry(line[0])
				} else {
					feature.Geometry = geojson.NewMultiPointGeometry(line...)
				}
			} else if geom_type == 2 {

				size := float64(extent) * float64(math.Pow(2, float64(tile.TileID.Z)))
				x0 := float64(extent * int(tile.TileID.X))
				y0 := float64(extent * int(tile.TileID.Y))

				tile.Buf.Pos = feature_geometry

				end := tile.Buf.ReadVarint() + tile.Buf.Pos
				cmd, length, x, y := 0, 0, 0.0, 0.0
				line := [][]float64{}
				lines := [][][]float64{}
				var pt []float64

				var cmdLen int
				for tile.Buf.Pos < end {
					if length == 0 {
						cmdLen = tile.Buf.ReadVarint()
						cmd = cmdLen & 0x7
						length = cmdLen >> 3
					}
					length--

					if cmd == 1 || cmd == 2 {
						x += tile.Buf.ReadSVarint()
						y += tile.Buf.ReadSVarint()
						pt = []float64{x, y}

						if (cmd == 1) && len(line) > 0 { // moveTo
							lines = append(lines, line)

							line = [][]float64{}
						}
						//if len(line)

						line = append(line, pt)
						//line.push(new Point(x, y));
					}
					//fmt.Println(length,tile.Buf.Pos,end)
					//if length < 0 && tile.Buf.Pos + 1 == end {
					//    tile.Buf.Pos += 1
					//}
				}
				if len(line) > 0 {
					lines = append(lines, line)
				}
				for i := range lines {
					lines[i] = Project(lines[i], x0, y0, size)
				}
				if len(lines) == 1 {
					feature.Geometry = geojson.NewLineStringGeometry(lines[0])
				} else {
					feature.Geometry = geojson.NewMultiLineStringGeometry(lines...)
				}
			} else if geom_type == 3 {

				size := float64(extent) * float64(math.Pow(2, float64(tile.TileID.Z)))
				x0 := float64(extent * int(tile.TileID.X))
				y0 := float64(extent * int(tile.TileID.Y))

				tile.Buf.Pos = feature_geometry

				end := tile.Buf.ReadVarint() + tile.Buf.Pos
				cmd, length, x, y := 0, 0, 0.0, 0.0
				line := [][]float64{}
				polygons := [][][][]float64{}
				var pt []float64

				var cmdLen int
				for tile.Buf.Pos < end {

					if length <= 0 {
						cmdLen = tile.Buf.ReadVarint()
						cmd = cmdLen & 0x7
						length = cmdLen >> 3
						//fmt.Println(cmdLen)
					}
					length--

					if cmd == 1 || cmd == 2 {
						x += tile.Buf.ReadSVarint()
						y += tile.Buf.ReadSVarint()
						pt = []float64{x, y}
						if (cmd == 1) && len(line) > 0 { // moveTo
							line = [][]float64{}
						}
						//if len(line)

						line = append(line, pt)
						//line.push(new Point(x, y));
					} else if cmd == 7 {
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

				if len(line) > 0 && len(polygons) == 0 {
					polygons = append(polygons, [][][]float64{line})
				}

				for i := range polygons {
					for j := range polygons[i] {
						polygons[i][j] = Project(polygons[i][j], x0, y0, size)
					}
				}
				if len(polygons) == 1 {
					feature.Geometry = geojson.NewPolygonGeometry(polygons[0])
				} else if len(polygons) > 0 {
					feature.Geometry = geojson.NewMultiPolygonGeometry(polygons...)
				}
		    }
	*/

	if id != 0 {
		feature.ID = id
	}
	return feature
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
		y2 := 180 - (p[1]+y0)*360.0/size
		line[j] = []float64{
			(p[0]+x0)*360.0/size - 180.0,
			360.0/math.Pi*math.Atan(math.Exp(y2*math.Pi/180.0)) - 90.0}
	}
	return line
}

func DeltaDim(num int) float64 {
	if num%2 == 1 {
		return float64((num + 1) / -2)
	} else {
		return float64(num / 2)
	}
	return float64(0)
}

func Geometry2(tile *Tile, feature_geometry int, geom_type int) *geojson.Geometry {
	tile.Buf.Pos = feature_geometry
	geom := tile.Buf.ReadPackedUInt32()

	switch geom_type {
	case 1:
		if len(geom) == 3 {
			// single point
			return geojson.NewPointGeometry([]float64{DeltaDim(int(geom[1])), DeltaDim(int(geom[2]))})
		} else {
			pos := 0
			cmdLen := geom[3]
			length := int(cmdLen >> 3)
			line := make([][]float64, length/2+1)
			line[0] = []float64{DeltaDim(int(geom[1])), DeltaDim(int(geom[2]))}
			ii := 1

			for i := pos + 3; i < pos+3+length*2; {
				line[ii] = []float64{DeltaDim(int(geom[i])), DeltaDim(int(geom[i+1]))}
				i += 2
				ii++
			}
			return geojson.NewMultiPointGeometry(line...)
		}
	case 2:
		// linestring
		fmt.Println(geom)
		pos := 0
		lines := [][][]float64{}
		for pos < len(geom) {
			cmdLen := geom[pos+3]
			length := int(cmdLen >> 3)
			line := make([][]float64, length/2+1)
			line[0] = []float64{DeltaDim(int(geom[1])), DeltaDim(int(geom[2]))}
			ii := 1
			for i := pos + 3; ii <= len(line) && i+1 < len(geom); {
				fmt.Println(i, ii, length, []float64{DeltaDim(int(geom[i])), DeltaDim(int(geom[i+1]))})
				line[ii] = []float64{DeltaDim(int(geom[i])), DeltaDim(int(geom[i+1]))}
				i += 2
				ii++
			}
			pos += (3 + length*2 + 1)
			lines = append(lines, line)
		}
		if len(lines) == 1 {
			return geojson.NewLineStringGeometry(lines[0])
		} else {
			return geojson.NewMultiLineStringGeometry(lines...)
		}

	case 3:
		// polygon

	}

	return &geojson.Geometry{}
}

func Geometry(tile *Tile, feature_geometry int, geom_type int, tileid m.TileID, extent int) *geojson.Geometry {
	size := float64(extent) * float64(math.Pow(2, float64(tile.TileID.Z)))
	x0 := float64(extent * int(tile.TileID.X))
	y0 := float64(extent * int(tile.TileID.Y))
	tile.Buf.Pos = feature_geometry
	geom := tile.Buf.ReadPackedUInt32()

	pos := 0
	var lines [][][]float64
	var polygons [][][][]float64

	for pos < len(geom) {
		if geom[pos] == 9 {
			pos += 1
			firstpt := []float64{DeltaDim(int(geom[pos])), DeltaDim(int(geom[pos+1]))}
			pos += 2
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
				}
				lines = append(lines, line[:i])

			} else {
				line := [][]float64{firstpt}
				lines = append(lines, line)

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

	if len(lines) == 1 {
		polygons = append(polygons, lines)
	} else {
		for _, line := range lines {
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

	for i := range polygons {
		for j := range polygons[i] {
			polygons[i][j] = Project(polygons[i][j], x0, y0, size)
		}
	}

	switch geom_type {
	case 1:
		if len(polygons[0][0]) == 1 {
			return geojson.NewPointGeometry(polygons[0][0][0])
		} else {
			return geojson.NewMultiPointGeometry(polygons[0][0]...)

		}
	case 2:
		if len(polygons[0]) == 1 {
			return geojson.NewLineStringGeometry(polygons[0][0])
		} else {
			return geojson.NewMultiLineStringGeometry(polygons[0]...)

		}
	case 3:
		if len(polygons) == 1 {
			return geojson.NewPolygonGeometry(polygons[0])
		} else {
			return geojson.NewMultiPolygonGeometry(polygons...)

		}
	}

	return &geojson.Geometry{}
}
