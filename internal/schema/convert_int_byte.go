//internal\schema\convert_int_byte.go

package schema

func int16ToBytes(data []int16) []byte {
	out := make([]byte, len(data)*2)
	for i, v := range data {
		out[i*2] = byte(v)
		out[i*2+1] = byte(v >> 8)
	}
	return out
}

func bytesToFloat64(b []byte) []float64 {
	out := make([]float64, len(b))
	for i := range b {
		out[i] = float64(b[i])
	}
	return out
}
