//MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil/q16.go

package mathutil

// Q16EnvConfd represents a fixed-point confidence value (0 to 65535)
type Q16EnvConfd uint16

const Q16Max = 65535

func Q16FromFloat(f float64) Q16EnvConfd {
	if f <= 0 { return 0 }
	if f >= 1 { return Q16Max }
	return Q16EnvConfd(f * Q16Max)
}

func (q Q16EnvConfd) Float64() float64 {
	return float64(q) / Q16Max
}