package pbf

import (
	//"io/ioutil"
	//"fmt"
	"bytes"
	"encoding/binary"
	"fmt"
	//"vector-tile/2.1"
	//"github.com/golang/protobuf/proto"
)

// 
type PBF struct {
	Pbf []byte
	Pos int
	Length int
}

const maxVarintBytes = 10 // maximum Length of a varint

// EncodeVarint returns the varint encoding of x.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
// Not used by the package itself, but helpful to clients
// wishing to use the same encoding.
func EncodeVarint(x uint64) []byte {
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

// DecodeVarint reads a varint-encoded integer from the slice.
// It returns the integer and the number of bytes consumed, or
// zero if there is not enough.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func DecodeVarint(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}

// a much faster key integration (microseconds to nanoseconds)
// returns the value number and key number for a given byte
func Key(x byte) (byte, byte) {
	//fmt.Printf("%08b\n",x)
	val := x >> 3

	// if the x value has a value in the 8 place
	if int(x) >= 8 {
		x = x & 0x07

	} else {
		return val, x
	}
	// if the x value has a value in the 16 place
	if int(x) >= 16 {
		x = x & 0x0f

	} else {
		return val, x
	}

	if int(x) >= 32 {
		x = x & 0x1f

	} else {
		return val, x
	}

	if int(x) >= 64 {
		x = x & 0x3f

	} else {
		return val, x
	}

	if int(x) >= 128 {
		x = x & 0x7f

	} else {
		return val, x
	}

	return val, x

}


func ReadInt32(buf []byte) int32 {
	if len(buf) == 4 {
    	return int32(((int(buf[0])) | (int(buf[1]) << 8) | (int(buf[2]) << 16)) + (int(buf[3]) << 24))
	} else if len(buf) == 3 {
    	return int32(((int(buf[0])) | (int(buf[1]) << 8) | (int(buf[2]) << 16)))
	} else if len(buf) == 2 {
    	return int32(((int(buf[0])) | (int(buf[1]) << 8)))
	} else if len(buf) == 1 {
    	return int32(buf[0])
	}
	return int32(0)
}

func ReadUInt32(buf []byte) uint32 {
	if len(buf) == 4 {
    	return uint32(((int(buf[0])) | (int(buf[1]) << 8) | (int(buf[2]) << 16)) + (int(buf[3]) * 0x1000000))
	} else if len(buf) == 3 {
    	return uint32(((int(buf[0])) | (int(buf[1]) << 8) | (int(buf[2]) << 16)))
	} else if len(buf) == 2 {
    	return uint32(((int(buf[0])) | (int(buf[1]) << 8)))
	} else if len(buf) == 1 {
    	return uint32(buf[0])
	}
	return uint32(0)
}

// reads a uint64 from a list of bytes
func ReadUint64(bytes []byte) uint64 {
	v, _ := DecodeVarint(bytes)
	return v
}

// reads a uint64 from a list of bytes
func ReadInt64(bytes []byte) int64 {
	v, _ := DecodeVarint(bytes)
	return int64(v)
}



func (pbf *PBF) ReadKey() (byte,byte) {
	var key,val byte
	if pbf.Pos > pbf.Length - 1 {
		key,val = 100,100
	} else {
		key,val = Key(pbf.Pbf[pbf.Pos])
		pbf.Pos += 1

	}

	return key,val
}


func (pbf *PBF) ReadVarint() int {
	if pbf.Pos + 1 >= pbf.Length {
		return 0
	}
	startPos := pbf.Pos 
	for pbf.Pbf[pbf.Pos] > 127 {
		pbf.Pos += 1
	}
	pbf.Pos += 1
	val, _ := DecodeVarint(pbf.Pbf[startPos:pbf.Pos])
	return int(val)
}

func (pbf *PBF) ReadSVarint() float64 {
	num := pbf.ReadVarint()
	if num%2 == 1 {
		return float64((num + 1) / -2)
	} else {
		return float64(num / 2)
	}
	return float64(0)
}

// var int bytes
func (pbf *PBF) Varint() []byte {
	startPos := pbf.Pos 
	for pbf.Pbf[pbf.Pos] > 127 {
		pbf.Pos += 1
	}
	pbf.Pos += 1
	return pbf.Pbf[startPos:pbf.Pos]
}



func (pbf *PBF) ReadFixed32() uint32 {
	val := ReadUInt32(pbf.Pbf[pbf.Pos:pbf.Pos+4])

	pbf.Pos += 4
	return val
}

func (pbf *PBF) ReadUInt32() uint32 {
	return ReadUInt32(pbf.Varint())
}




func (pbf *PBF) ReadSFixed32() int32 {
	val := ReadInt32(pbf.Pbf[pbf.Pos:pbf.Pos+4])
	pbf.Pos += 4
	return val
}

func (pbf *PBF) ReadInt32() int32 {
	return ReadInt32(pbf.Varint())
}

// reads a uint64 from a list of bytes
func (pbf *PBF) ReadFixed64() uint64 {
	v, _ := DecodeVarint(pbf.Pbf[pbf.Pos:pbf.Pos+8])
	pbf.Pos += 8
	return v
}

func (pbf *PBF) ReadUInt64() uint64 {
	return ReadUint64(pbf.Varint())
}

// reads a uint64 from a list of bytes
func (pbf *PBF) ReadSFixed64() int64 {
	v, _ := DecodeVarint(pbf.Pbf[pbf.Pos:pbf.Pos+8])
	pbf.Pos += 8
	return int64(v)
}


func (pbf *PBF) ReadInt64() int64 {
	return ReadInt64(pbf.Varint())
}


func (pbf *PBF) ReadFloat() float32 {
	var pi32 float32
	buf := bytes.NewReader(pbf.Pbf[pbf.Pos:pbf.Pos+4])
	err := binary.Read(buf, binary.LittleEndian, &pi32)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pbf.Pos += 4
	return pi32
}

// reading a double
func (pbf *PBF) ReadDouble() float64 {
	var pi32 float64
	buf := bytes.NewReader(pbf.Pbf[pbf.Pos:pbf.Pos+8])
	err := binary.Read(buf, binary.LittleEndian, &pi32)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	pbf.Pos += 8
	return pi32
}

func (pbf *PBF) ReadString() string {
	size := pbf.ReadVarint()
	stringval := string(pbf.Pbf[pbf.Pos:pbf.Pos+size])
	pbf.Pos += size
	return stringval
}

func (pbf *PBF) ReadBool() bool {
	pbf.Byte()

	size := pbf.ReadVarint()
	buf := pbf.Pbf[pbf.Pos:pbf.Pos+size]
	if buf[0] == 1 {
		return true
	} else if buf[0] == 0 {
		return false
	}
	pbf.Pos += size
	return false
}




func (pbf *PBF) ReadPackedUInt32() []uint32 {
	size := pbf.ReadVarint()
	arr := []uint32{}
	endpos := pbf.Pos + size
	for pbf.Pos < endpos {
		arr = append(arr,pbf.ReadUInt32())
	}
	return arr
}

func (pbf *PBF) Byte() {
	fmt.Println(pbf.Pbf[pbf.Pos],"current")
	fmt.Println(pbf.Pbf[pbf.Pos-5:pbf.Pos+5],"next5")
}






