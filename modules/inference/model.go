//modules/inference/model.go

package inference

import "context"

type Model interface {

	Predict(
		ctx context.Context,
		req PredictionRequest,
	) (*PredictionResult, error)

	PredictBatch(
		ctx context.Context,
		req []PredictionRequest,
	) ([]PredictionResult, error)

	Health(ctx context.Context) error

	Warmup(ctx context.Context) error
}