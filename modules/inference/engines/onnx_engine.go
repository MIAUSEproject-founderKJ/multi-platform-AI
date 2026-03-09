//modules/inference/engines/onnx_engine.go

package engines

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	onnx "github.com/yalue/onnxruntime_go"
)

type ONNXEngine struct {
	session *onnx.Session

	inputName  string
	outputName string

	mu sync.RWMutex

	warm bool
}

func NewONNXEngine() *ONNXEngine {
	return &ONNXEngine{}
}

func (e *ONNXEngine) Load(modelPath string) error {

	session, err := onnx.NewSession(modelPath)
	if err != nil {
		return fmt.Errorf("failed to load ONNX model: %w", err)
	}

	e.session = session

	inputs := session.InputNames()
	outputs := session.OutputNames()

	if len(inputs) == 0 || len(outputs) == 0 {
		return errors.New("invalid ONNX model IO")
	}

	e.inputName = inputs[0]
	e.outputName = outputs[0]

	return nil
}

func (e *ONNXEngine) Predict(
	ctx context.Context,
	features []float32,
) ([]float32, error) {

	e.mu.RLock()
	session := e.session
	e.mu.RUnlock()

	if session == nil {
		return nil, errors.New("model not loaded")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	inputTensor, err := onnx.NewTensor(
		[]int64{1, int64(len(features))},
		features,
	)

	if err != nil {
		return nil, err
	}

	outputs, err := session.Run(
		map[string]*onnx.Tensor{
			e.inputName: inputTensor,
		},
	)

	if err != nil {
		return nil, err
	}

	resultTensor := outputs[e.outputName]

	return resultTensor.Float32Data(), nil
}

func (e *ONNXEngine) PredictBatch(
	ctx context.Context,
	batch [][]float32,
) ([][]float32, error) {

	if len(batch) == 0 {
		return nil, errors.New("empty batch")
	}

	featureLen := len(batch[0])

	flat := make([]float32, 0, len(batch)*featureLen)

	for _, f := range batch {
		flat = append(flat, f...)
	}

	inputTensor, err := onnx.NewTensor(
		[]int64{int64(len(batch)), int64(featureLen)},
		flat,
	)

	if err != nil {
		return nil, err
	}

	e.mu.RLock()
	session := e.session
	e.mu.RUnlock()

	outputs, err := session.Run(
		map[string]*onnx.Tensor{
			e.inputName: inputTensor,
		},
	)

	if err != nil {
		return nil, err
	}

	resultTensor := outputs[e.outputName]

	raw := resultTensor.Float32Data()

	results := make([][]float32, len(batch))

	stride := len(raw) / len(batch)

	for i := range results {

		start := i * stride
		end := start + stride

		results[i] = raw[start:end]
	}

	return results, nil
}

func (e *ONNXEngine) Warmup(ctx context.Context) error {

	if e.warm {
		return nil
	}

	dummy := []float32{0, 0}

	for i := 0; i < 10; i++ {

		_, err := e.Predict(ctx, dummy)

		if err != nil {
			return err
		}
	}

	e.warm = true

	return nil
}

func (e *ONNXEngine) Health(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := e.Predict(ctx, []float32{0, 0})

	return err
}

func (e *ONNXEngine) Close() error {

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.session != nil {
		e.session.Close()
		e.session = nil
	}

	return nil
}