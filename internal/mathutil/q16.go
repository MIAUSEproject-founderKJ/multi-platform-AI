//MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil/q16.go

package mathutil

// Q16 represents a fixed-point confidence value (0 to 65535)
// where 0 is 0% and 65535 is 100%.
type Q16 uint16

const (
	Q16Max  = 65535
	Q16Half = 32768
)

// ToFloat64 converts a Q16 fixed-point (uint16) to a float64 (0.0 - 1.0)
func ToFloat64(q uint16) float64 {
	return float64(q) / Q16Max
}

// FromFloat64 converts a float64 (0.0 - 1.0) to a Q16 fixed-point (uint16)
func FromFloat64(f float64) uint16 {
	if f > 1.0 {
		f = 1.0
	}
	if f < 0.0 {
		f = 0.0
	}
	return uint16(f * Q16Max)
}

// Multiply performs a fixed-point multiplication: (a * b) / Q16Max
// This is essential for the Bayesian recursive updates in trust_bayesian.go
func (q Q16) Multiply(other Q16) Q16 {
	result := (uint32(q) * uint32(other)) / Q16Max
	return Q16(result)
}

// Percentage returns a rounded integer percentage for the HUD
func (q Q16) Percentage() int {
	return int((float64(q) / Q16Max) * 100)
}
