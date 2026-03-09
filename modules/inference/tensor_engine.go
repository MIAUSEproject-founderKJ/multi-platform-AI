//modules/inference/tensor_engine.go

package inference

import "context"

type TensorEngine interface {

	Predict(
		ctx context.Context,
		features []float32,
	) ([]float32, error)

	Load(modelPath string) error

	Close() error
}