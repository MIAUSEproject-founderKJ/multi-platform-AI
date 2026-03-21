// modules/audio/feature.go
package audio

import (
	"fmt"
	"math"
)

type FeatureExtractor struct {
	sampleRate int
}

func NewFeatureExtractor(rate int) *FeatureExtractor {
	return &FeatureExtractor{sampleRate: rate}
}

func (f *FeatureExtractor) ProcessPCM(pcm []byte) ([]float64, error) {

	samples := bytesToFloat64(pcm)

	// 1. Pre-emphasis
	for i := 1; i < len(samples); i++ {
		samples[i] = samples[i] - 0.97*samples[i-1]
	}

	// 2. Windowing (Hamming)
	for i := range samples {
		samples[i] *= 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(len(samples)-1))
	}

	// 3. FFT
	spectrum := fft(samples)
	if len(spectrum) < 40 {
		return nil, fmt.Errorf("insufficient spectrum size")
	}
	// 4. Log energy
	features := make([]float64, len(spectrum))
	for i := range spectrum {
		features[i] = math.Log(math.Abs(spectrum[i]) + 1e-6)
	}

	return features[:40], nil // first 40 bins
}
