// internal/encoding/convert_int_byte.go
package encoding

import (
	"encoding/binary"
	"math"
)

// Int16ToBytes converts int16 slice into deterministic little-endian bytes
func Int16ToBytes(data []int16) []byte {
	out := make([]byte, len(data)*2)
	for i, v := range data {
		binary.LittleEndian.PutUint16(out[i*2:], uint16(v))
	}
	return out
}

// BytesToInt16 restores int16 slice safely
func BytesToInt16(b []byte) []int16 {
	out := make([]int16, len(b)/2)
	for i := range out {
		out[i] = int16(binary.LittleEndian.Uint16(b[i*2:]))
	}
	return out
}

// Float64ToBytes safely encodes float64 deterministically
func Float64ToBytes(data []float64) []byte {
	out := make([]byte, len(data)*8)
	for i, v := range data {
		binary.LittleEndian.PutUint64(out[i*8:], math.Float64bits(v))
	}
	return out
}

// BytesToFloat64 restores float64 slice correctly
func BytesToFloat64(b []byte) []float64 {
	out := make([]float64, len(b)/8)
	for i := range out {
		bits := binary.LittleEndian.Uint64(b[i*8:])
		out[i] = math.Float64frombits(bits)
	}
	return out
}
