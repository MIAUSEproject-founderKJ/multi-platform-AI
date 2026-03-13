// core/optimization/distiller.go
package optimization

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot"

type ModelOptimizer interface {
	Distill(model Model) (Model, error)
	Quantize(model Model) (Model, error)
	Prune(model Model, threshold float64) (Model, error)
}

func (m *NLUModule) Init(ctx *boot.RuntimeContext) error {

	rawModel := LoadModelForPlatform(ctx.PlatformClass)

	optimized, err := ctx.Optimizer.Distill(rawModel)
	if err != nil {
		return err
	}

	m.model = optimized
	return nil
}
