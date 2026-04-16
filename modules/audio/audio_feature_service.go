// modules\audio\audio_feature_service.go
package audio

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"gonum.org/v1/gonum/dsp/fourier"
)

type FeatureExtractor struct {
	sampleRate int
}

func fft(x []float64) []complex128 {
	fft := fourier.NewFFT(len(x))
	return fft.Coefficients(nil, x)
}

func NewFeatureExtractor(rate int) *FeatureExtractor {
	return &FeatureExtractor{sampleRate: rate}
}

func (f *FeatureExtractor) ProcessPCM(pcm []byte) ([]float64, error) {

	samples := schema.BytesToFloat64(pcm)

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
		features[i] = math.Log(cmplx.Abs(spectrum[i]) + 1e-6)
	}

	return features[:40], nil // first 40 bins
}
