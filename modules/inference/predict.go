//modules/inference/predict.go

package inference

import (
	"context"
)

func (m *ModelAdapter) Predict(
	ctx context.Context,
	req PredictionRequest,
) (*PredictionResult, error) {

	features := buildFeatures(req)

	output, err := m.engine.Predict(ctx, features)

	if err != nil {
		return nil, err
	}

	result := &PredictionResult{
		DeviceID:   req.DeviceID,
		Timestamp:  req.Timestamp,
		Prediction: "value",
		Confidence: float64(output[0]),
	}

	return result, nil
}

func (m *ModelAdapter) PredictBatch(
	ctx context.Context,
	req []PredictionRequest,
) ([]PredictionResult, error) {

	results := make([]PredictionResult, 0, len(req))

	for _, r := range req {

		pred, err := m.Predict(ctx, r)

		if err != nil {
			return nil, err
		}

		results = append(results, *pred)
	}

	return results, nil
}

func (m *ModelAdapter) Health(ctx context.Context) error {

	_, err := m.engine.Predict(ctx, []float32{0, 0})

	return err
}

func (m *ModelAdapter) Warmup(ctx context.Context) error {

	for i := 0; i < 10; i++ {

		_, err := m.engine.Predict(ctx, []float32{0, 0})

		if err != nil {
			return err
		}
	}

	return nil
}
