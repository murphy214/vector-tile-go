package tags

import (
	"github.com/murphy214/pbf"
	"reflect"
	
)

/*

uint64_t type = complex_value & 0x0F; // least significant 4 bits
uint64_t parameter = complex_value >> 4;

    Type     | Id  | Parameter
---------------------------------
string       |  0  | index into layer string_values
float        |  1  | index into layer float_values
double       |  2  | index into layer double_values
uint         |  3  | index into layer int_values
sint         |  4  | index into layer int_values (values are zigzag encoded)
inline uint  |  5  | value of unsigned integer (values between 0 to 2^60-1)
inline sint  |  6  | value of zigzag-encoded integer (values between -2^59 to 2^59-1)
bool/null    |  7  | value of 0 = null, 1 = false, 2 = true
list         |  8  | value is the number of list items to follow:
             |     |   each item in the list is a complex value
map          |  9  | value is the number of key-value pairs to follow:
             |     |   each pair is an index into layer keys
             |     |   followed by a complex_value for the value
delta-       | 10  | parameter is the number of items N in the list:
  encoded    |     |   one uint64 is an index into the Layer's attribute_scalings
  list       |     |   followed by N uint64 nullable deltas for the list items
*/

type DataType int

const (
	StringType     DataType = 0
	FloatType      DataType = 1
	DoubleType     DataType = 2
	UintType       DataType = 3
	SintType       DataType = 4
	InlineUintType DataType = 5
	InlineSintType DataType = 6
	BoolType       DataType = 7
	ListType       DataType = 8
	MapType        DataType = 9
	DeltaListType  DataType = 10
)

// gets the data type
func GetIndexType(val int) (int, DataType) {
	return val >> 4, DataType(val & 0x0F)
}

// gets the value integer
func GetValInt(index int, typeval DataType) int {
	return index<<4 + int(typeval)
}

// tags writer structuretypeval
type TagWriter struct {
	StringValues []string
	StringMap    map[string]int
	StringBytes []byte 

	FloatValues []float32
	FloatMap    map[float32]int
	FloatBytes []byte
	
	DoubleValues []float64
	DoubleMap    map[float64]int
	DoubleBytes []byte

	IntValues []int
	IntMap    map[int]int
	IntBytes []byte

	Keys    []string
	KeysMap map[string]int
	KeysBytes []byte
}

func NewTagWriter() *TagWriter {
	return &TagWriter{
		StringMap: map[string]int{},
		FloatMap:  map[float32]int{},
		DoubleMap: map[float64]int{},
		IntMap:    map[int]int{},
		KeysMap:   map[string]int{},
	}
}

// adds a key to the tag writer
func (tagwriter *TagWriter) AddKey(key string) int {
	valint, boolval := tagwriter.KeysMap[key]
	if !boolval {
		valint = len(tagwriter.Keys)
		tagwriter.KeysMap[key] = valint
		tagwriter.Keys = append(tagwriter.Keys, key)
		tagwriter.KeysBytes = append(tagwriter.KeysBytes, 26)
		tagwriter.KeysBytes = append(tagwriter.KeysBytes, pbf.EncodeVarint(uint64(len(key)))...)
		tagwriter.KeysBytes = append(tagwriter.KeysBytes, []byte(key)...)
	
	}
	return valint
}

// adds a single value (althoguh could be complex) to the tag writer struture
func (tagwriter *TagWriter) AddValue(val interface{}) (int, bool) {
	vv := reflect.ValueOf(val)
	kd := vv.Kind()

	// switching for each type
	switch kd {
	case reflect.String:
		myval := vv.String()
		valint, boolval := tagwriter.StringMap[myval]
		if !boolval {
			valint = GetValInt(len(tagwriter.StringMap), StringType)
			tagwriter.StringMap[myval] = valint
			tagwriter.StringValues = append(tagwriter.StringValues, myval)
			tagwriter.StringBytes = append(tagwriter.StringBytes, 50)
			tagwriter.StringBytes = append(tagwriter.StringBytes, pbf.EncodeVarint(uint64(len(myval)))...)
			tagwriter.StringBytes = append(tagwriter.StringBytes, []byte(myval)...)
		}
		return valint, true
	case reflect.Float32:
		myval := float32(vv.Float())
		valint, boolval := tagwriter.FloatMap[myval]
		if !boolval {
			valint = GetValInt(len(tagwriter.FloatMap), FloatType)
			tagwriter.FloatMap[myval] = valint
			tagwriter.FloatValues = append(tagwriter.FloatValues, myval)
			tagwriter.FloatBytes = append(tagwriter.FloatBytes, FloatVal32Raw(myval)...)
		}
		return valint, true
	case reflect.Float64:
		myval := float64(vv.Float())
		valint, boolval := tagwriter.DoubleMap[myval]
		if !boolval {
			valint = GetValInt(len(tagwriter.DoubleMap), DoubleType)
			tagwriter.DoubleMap[myval] = valint
			tagwriter.DoubleValues = append(tagwriter.DoubleValues, myval)
			tagwriter.DoubleBytes = append(tagwriter.DoubleBytes,FloatVal64Raw(myval)...)
		}
		return valint, true

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		myval := int(vv.Int())
		valint, boolval := tagwriter.IntMap[myval]
		if !boolval {
			valint = GetValInt(len(tagwriter.IntMap), UintType)
			tagwriter.IntMap[myval] = valint
			tagwriter.IntValues = append(tagwriter.IntValues, myval)
			tagwriter.IntBytes = append(tagwriter.IntBytes,IntVal64Raw(myval)...)
		}
		return valint, true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		myval := int(vv.Uint())
		valint, boolval := tagwriter.IntMap[myval]
		if !boolval {
			valint = GetValInt(len(tagwriter.IntMap), UintType)
			tagwriter.IntMap[myval] = valint
			tagwriter.IntValues = append(tagwriter.IntValues, myval)
			tagwriter.IntBytes = append(tagwriter.IntBytes,IntVal64Raw(myval)...)
		}
		return valint, true
	case reflect.Bool:
		if vv.Bool() == true {
			return 39, true
		} else if vv.Bool() == false {
			return 23, true
		}
	}

	return 0, false
}

// needed to make a complex type
func (tagwriter *TagWriter) MakeComplex(val interface{}) []int {
	val_list, listbool := val.([]interface{})
	if listbool {
		total := []int{GetValInt(len(val_list), ListType)}
		for _, i := range val_list {
			simpletag, complete := tagwriter.AddValue(i)
			if complete {
				total = append(total, simpletag)
			} else {
				total = append(total, tagwriter.MakeComplex(i)...)
			}
		}
		return total
	}

	val_map, listbool := val.(map[string]interface{})
	if listbool {
		total := []int{GetValInt(len(val_map), MapType)}
		for k, v := range val_map {
			total = append(total,tagwriter.AddKeyValue(k, v)...)
		}
		return total
	}
	return []int{}
}

// adds a set of key values
func (tagwriter *TagWriter) AddKeyValue(key string, v interface{}) []int {
	simpletag, complete := tagwriter.AddValue(v)
	if complete {
		return []int{tagwriter.AddKey(key), simpletag}
	} else {
		return append([]int{tagwriter.AddKey(key)}, tagwriter.MakeComplex(v)...)
	}
}

// given a set of key value property features
// returns there given attribute tags
func (tagwriter *TagWriter) MakeProperties(props map[string]interface{}) []int {
	total := []int{}
	for k, v := range props {
		total = append(tagwriter.AddKeyValue(k, v), total...)
	}
	return total
}

// creates the tagsreader object from the writer
func (tagwriter *TagWriter) Reader() *TagReader {
	return &TagReader{
		StringValues:tagwriter.StringValues,
		IntValues:tagwriter.IntValues,
		FloatValues:tagwriter.FloatValues,
		DoubleValues:tagwriter.DoubleValues,
		Keys:tagwriter.Keys,
	}
}

