package vt

import (
	"errors"
	m "github.com/murphy214/mercantile"
	"github.com/murphy214/pbf"
	"github.com/paulmach/go.geojson"
	"math"
)

// upper vector tile structure
type Tile struct {
	LayerMap map[string]*Layer // generalized layer map with all the methods hanging off each layer
	Buf      *pbf.PBF
	TileID   m.TileID // the tileid of the given tile
	Layers   []string // the layer names in the order in which they occur within the vt
}

// create / reads a new vector tile from a byte array
func NewTile(bytevals []byte) (tile *Tile, err error) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("Error in NewTile.")
		}
	}()

	// creating vector tile
	tile = &Tile{
		LayerMap: map[string]*Layer{},
		Buf:      &pbf.PBF{Pbf: bytevals, Length: len(bytevals)},
	}
	for tile.Buf.Pos < tile.Buf.Length {
		key, val := tile.Buf.ReadKey()
		if key == 3 && val == 2 {
			size := tile.Buf.ReadVarint()
			if size != 0 {
				tile.NewLayer(tile.Buf.Pos + size)
			}

		}
	}
	return tile, err
}

// renders what is currently in the amp for me 
func (tile *Tile) Render() []byte {
	totalbs := []byte{}
	for _,v := range tile.LayerMap {
		bs := v.Buf.Pbf[v.StartPos:v.EndPos]
		prefix := pbf.EncodeVarint(uint64(len(bs)))
		layerbs := append(append([]byte{26},prefix...),bs...)
		totalbs = append(totalbs,layerbs...)
	}
	return totalbs
}

// deletes a given layer
func (tile *Tile) DeleteLayer(layername string) []byte {
	delete(tile.LayerMap,layername)
	bs := tile.Render()
	return bs
}

// deletes a given layer
func (tile *Tile) DeleteLayers(layernames []string) []byte {
	for _,layername := range layernames {
		delete(tile.LayerMap,layername)
	}
	bs := tile.Render()
	return bs
}



// create / reads a new vector tile from a byte array
func ReadTile(bytevals []byte, tileid m.TileID) (totalfeautures []*geojson.Feature, err error) {

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("Error in ReadTile")
		}
	}()

	// creating vector tile
	tile := &Tile{
		Buf:    pbf.NewPBF(bytevals),
		TileID: tileid,
	}
	totalfeautures = []*geojson.Feature{}
	for tile.Buf.Pos < tile.Buf.Length {
		key, val := tile.Buf.ReadKey()
		if key == 3 && val == 2 {
			sizex := tile.Buf.ReadVarint()
			endpos := tile.Buf.Pos + sizex
			//var layer *Layer
			var extent, number_features int
			var layername string
			var features []int
			var keys []string
			var values []interface{}
			if sizex != 0 {
				//layer = &Layer{StartPos: tile.Buf.Pos, EndPos: endpos}
				key, val := tile.Buf.ReadKey()
				for tile.Buf.Pos < endpos {
					if key == 1 && val == 2 {
						layername = tile.Buf.ReadString()
						key, val = tile.Buf.ReadKey()
					}
					// collecting all the features
					for key == 2 && val == 2 {
						// reading for features

						features = append(features, tile.Buf.Pos)
						feat_size := tile.Buf.ReadVarint()

						tile.Buf.Pos += feat_size
						key, val = tile.Buf.ReadKey()
					}
					// collecting all keys
					for key == 3 && val == 2 {
						keys = append(keys, tile.Buf.ReadString())
						key, val = tile.Buf.ReadKey()
					}
					// collecting all values
					for key == 4 && val == 2 {
						//tile.Buf.Byte()
						tile.Buf.ReadVarint()
						newkey, _ := tile.Buf.ReadKey()
						switch newkey {
						case 1:
							values = append(values, tile.Buf.ReadString())
						case 2:
							values = append(values, tile.Buf.ReadFloat())
						case 3:
							values = append(values, tile.Buf.ReadDouble())
						case 4:
							values = append(values, tile.Buf.ReadInt64())
						case 5:
							values = append(values, tile.Buf.ReadUInt64())
						case 6:
							values = append(values, tile.Buf.ReadUInt64())
						case 7:
							values = append(values, tile.Buf.ReadBool())
						}
						key, val = tile.Buf.ReadKey()
					}
					if key == 5 && val == 0 {
						extent = int(tile.Buf.ReadVarint())
						key, val = tile.Buf.ReadKey()
					}
					if key == 15 && val == 0 {
						_ = int(tile.Buf.ReadVarint())
						key, val = tile.Buf.ReadKey()

					}
				}
				if extent == 0 {
					extent = 4096
				}
				number_features = len(features)
				tile.Buf.Pos = endpos
			}
			feats := make([]*geojson.Feature, number_features)
			size := float64(extent) * float64(math.Pow(2, float64(tile.TileID.Z)))
			x0 := float64(extent) * float64(tile.TileID.X)
			y0 := float64(extent) * float64(tile.TileID.Y)
			var feature_geometry, id, geom_type int
			if extent == 0 {
				extent = 4096
			}
			for i, pos := range features {
				tile.Buf.Pos = pos
				endpos := tile.Buf.Pos + tile.Buf.ReadVarint()
				//startpos := tile.Buf.Pos

				feature := &geojson.Feature{Properties: map[string]interface{}{}}

				for tile.Buf.Pos < endpos {
					key, val := tile.Buf.ReadKey()

					// logic for handlign id
					if key == 1 && val == 0 {
						id = int(tile.Buf.ReadUInt64())
					}
					// logic for handling tags
					if key == 2 && val == 2 {
						tags := tile.Buf.ReadPackedUInt32()
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

				tile.Buf.Pos = feature_geometry
				geom := tile.Buf.ReadPackedUInt32()
				//fmt.Println(geom)
				pos := 0
				var lines [][][]float64
				var polygons [][][][]float64
				var firstpt []float64
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

				for i := range polygons {
					for j := range polygons[i] {
						polygons[i][j] = Project(polygons[i][j], x0, y0, size)
					}
				}

				switch geom_type {
				case 1:
					if len(polygons[0][0]) == 1 {
						feature.Geometry = geojson.NewPointGeometry(polygons[0][0][0])
					} else {
						feature.Geometry = geojson.NewMultiPointGeometry(polygons[0][0]...)

					}
				case 2:
					if len(polygons[0]) == 1 {
						feature.Geometry = geojson.NewLineStringGeometry(polygons[0][0])
					} else {
						feature.Geometry = geojson.NewMultiLineStringGeometry(polygons[0]...)

					}
				case 3:
					if len(polygons) == 1 {
						feature.Geometry = geojson.NewPolygonGeometry(polygons[0])
					} else {
						feature.Geometry = geojson.NewMultiPolygonGeometry(polygons...)

					}
				}

				if id != 0 {
					feature.ID = id
				}
				feature.Properties[`layer`] = layername
				feats[i] = feature
			}

			totalfeautures = append(totalfeautures, feats...)
			tile.Buf.Pos = endpos

		}
	}
	if len(totalfeautures) == 0 {
		err = errors.New("No features read from given tile.")
	}
	return totalfeautures, err

}
