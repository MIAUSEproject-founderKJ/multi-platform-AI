//MIAUSEproject-founderKJ/multi-platform-AI/internal/mathutil/q16.go

package mathutil

import "math"

// Q16 represents a fixed-point confidence value (0 to 65535)
// where 0 is 0% and 65535 is 100%.
type Q16 uint16

const (
	Q16Max  = 65535
	Q16Half = 32768
)

// Q16FromFloat converts a float (0.0 to 1.0) to Q16
func Q16FromFloat(f float64) Q16 {
	f = math.Max(0, math.Min(1, f))
	return Q16(f * Q16Max)
}

// Float64 converts Q16 back to float for HMI/Logging
func (q Q16) Float64() float64 {
	return float64(q) / Q16Max
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