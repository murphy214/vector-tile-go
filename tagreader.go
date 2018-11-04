package vt

// delta encoding 
func DeltaDim(num int) float64 {
	if num%2 == 1 {
		return float64((num + 1) / -2)
	} else {
		return float64(num / 2)
	}
	return float64(0)
}

// the tag reader struct
type TagReader struct {
	StringValues []string
	FloatValues  []float32
	DoubleValues []float64
	IntValues    []int
	Keys         []string
	Tags         []int
	Pos          int
}

// checks to see if we can get the next tag integer
func (tagreader *TagReader) Next() bool {
	return len(tagreader.Tags) > tagreader.Pos
}

// gets the next tag integer
func (tagreader *TagReader) Tag() int {
	val := tagreader.Tags[tagreader.Pos]
	tagreader.Pos++
	return val
}

// gets the next two tags one representing a key
// the other representing either the first complex tag or a simple tag
func (tagreader *TagReader) NextKeyValue() bool {
	return len(tagreader.Tags) > tagreader.Pos+1
}

// gets the tags in nextkeyvalue
func (tagreader *TagReader) TagVal() (int,int) {
	val, val2 := tagreader.Tags[tagreader.Pos], tagreader.Tags[tagreader.Pos+1]
	tagreader.Pos += 2
	return val, val2
}

// Reads a tag
func (tagsreader *TagReader) ReadTag(tag int) interface{} {
	indexval, t := GetIndexType(tag)
	switch t {
	case StringType:
		return tagsreader.StringValues[indexval]
	case FloatType:
		return tagsreader.FloatValues[indexval]
	case DoubleType:
		return tagsreader.DoubleValues[indexval]
	case UintType:
		return tagsreader.IntValues[indexval]
	case SintType:
		return DeltaDim(tagsreader.IntValues[indexval])
	case InlineSintType:
		return DeltaDim(indexval)
	case InlineUintType:
		return indexval
	case BoolType:
		if indexval == 2 {
			return true
		} else if indexval == 1 {
			return false
		} else if indexval == 3 {
			return nil
		}
	case ListType:
		i := 0
		total := []interface{}{}
		for tagsreader.Next() && i < indexval {
			tag := tagsreader.Tag()
			total = append(total, tagsreader.ReadTag(tag))
			i++
		}
		return total
	case MapType:
		i := 0
		total := map[string]interface{}{}
		for tagsreader.NextKeyValue() && i < indexval {
			k, v := tagsreader.TagVal()
			key := tagsreader.Keys[k]
			total[key] = tagsreader.ReadTag(v)
			i++
		}
		return total
	}
	return 0

}

// reads all the tags 
func (tagsreader *TagReader) ReadTags(tags []int) map[string]interface{} {
	tagsreader.Tags = tags
	total := map[string]interface{}{}
	for tagsreader.NextKeyValue() {
		k, v := tagsreader.TagVal()
		key := tagsreader.Keys[k]		
		total[key] = tagsreader.ReadTag(v)
	}
	return total
}
