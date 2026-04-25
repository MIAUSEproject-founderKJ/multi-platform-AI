//internal/mathutil/q16.go

package mathutil

import "math"

// Q16 represents a fixed-point value in [0, 1] using 16-bit precision.
// 0     = 0.0
// 65535 = 1.0
type Q16 uint16

const (
	Max  Q16 = 65535
	Min  Q16 = 0
	Half Q16 = 32768
)

// ============================================================
// CONSTRUCTORS (BOUNDARY INPUTS)
// ============================================================

// FromFloat64 converts float64 → Q16 with clamping and rounding.
// Use only at system boundaries (ML outputs, external inputs).
func FromFloat64(f float64) Q16 {
	if f <= 0 {
		return Min
	}
	if f >= 1 {
		return Max
	}
	return Q16(math.Round(f * float64(Max)))
}

// MustFromFloat64 panics if out of range.
// Use only in trusted/internal code paths.
func MustFromFloat64(f float64) Q16 {
	if f < 0 || f > 1 {
		panic("mathutil.Q16: value out of range")
	}
	return Q16(math.Round(f * float64(Max)))
}

// ============================================================
// INTEGER-BASED CONVERSION (CORE PATH)
// ============================================================

// FromRatio computes (num / denom) → Q16 using integer math.
// This is the preferred method for scoring systems.
func FromRatio(num, denom uint32) Q16 {
	if denom == 0 {
		return Min
	}
	if num >= denom {
		return Max
	}

	// (num * Max) fits in uint32? → yes (max 65535 * 65535 safe in uint32)
	return Q16((num*uint32(Max) + denom/2) / denom)
}

// FromPercentage converts integer percentage (0–100) → Q16.
func FromPercentage(p uint8) Q16 {
	if p >= 100 {
		return Max
	}
	return Q16((uint32(p) * uint32(Max)) / 100)
}

// ============================================================
// NORMALIZATION / SCALING
// ============================================================

// Scale maps a value from [min, max] → Q16.
// Useful for sensor normalization and AI signal preprocessing.
func Scale(value, min, max float64) Q16 {
	if max <= min {
		return Min
	}
	if value <= min {
		return Min
	}
	if value >= max {
		return Max
	}

	ratio := (value - min) / (max - min)
	return FromFloat64(ratio)
}

// ============================================================
// CONVERSION (OUTPUT)
// ============================================================

// Float64 converts Q16 → float64 (for logging/UI/export only).
func (q Q16) Float64() float64 {
	return float64(q) / float64(Max)
}

// Percentage returns integer percentage (0–100).
func (q Q16) Percentage() int {
	return int((uint32(q) * 100) / uint32(Max))
}

// ============================================================
// COMPARISON
// ============================================================

func (q Q16) GT(o Q16) bool  { return q > o }
func (q Q16) GTE(o Q16) bool { return q >= o }
func (q Q16) LT(o Q16) bool  { return q < o }
func (q Q16) LTE(o Q16) bool { return q <= o }
func (q Q16) EQ(o Q16) bool  { return q == o }

// ============================================================
// ARITHMETIC (DETERMINISTIC, SATURATING)
// ============================================================

// Add performs saturating addition.
func (q Q16) Add(o Q16) Q16 {
	sum := uint32(q) + uint32(o)
	if sum > uint32(Max) {
		return Max
	}
	return Q16(sum)
}

// Sub performs saturating subtraction.
func (q Q16) Sub(o Q16) Q16 {
	if q < o {
		return Min
	}
	return q - o
}

// Mul performs fixed-point multiplication with rounding.
// (q * o) / Max
func (q Q16) Mul(o Q16) Q16 {
	product := uint32(q) * uint32(o)
	return Q16((product + uint32(Max)/2) / uint32(Max))
}

// Div performs fixed-point division with rounding.
// (q * Max) / o
func (q Q16) Div(o Q16) Q16 {
	if o == 0 {
		return Min // safe default for confidence systems
	}
	num := uint32(q) * uint32(Max)
	return Q16((num + uint32(o)/2) / uint32(o))
}

// ============================================================
// AI / SIGNAL FUSION UTILITIES
// ============================================================

// WeightedBlend computes:
// result = a*weight + b*(1-weight)
func WeightedBlend(a, b, weight Q16) Q16 {
	inv := Max - weight
	return a.Mul(weight).Add(b.Mul(inv))
}

// Invert returns (1 - q)
func (q Q16) Invert() Q16 {
	return Max - q
}

// ============================================================
// FAST CHECKS
// ============================================================

func (q Q16) IsZero() bool { return q == Min }
func (q Q16) IsMax() bool  { return q == Max }

// Clamp ensures value is within valid bounds.
func (q Q16) Clamp() Q16 {
	if q > Max {
		return Max
	}
	return q
}

// WeightedSum calculates the weighted sum of values with weights.
// Assumes values[i] and weights[i] are Q16 (0..65535)
func WeightedSum(values []Q16, weights []Q16) Q16 {
	if len(values) != len(weights) || len(values) == 0 {
		return 0
	}

	var total, sumWeights uint32 // use 32-bit for safety

	for i := 0; i < len(values); i++ {
		// Multiply value and weight, then divide by Max to scale back to Q16
		total += uint32(values[i]) * uint32(weights[i]) / uint32(Max)
		sumWeights += uint32(weights[i])
	}

	if sumWeights == 0 {
		return 0
	}

	// Normalize total by sum of weights
	return Q16(uint32(total) * uint32(Max) / sumWeights)
}
