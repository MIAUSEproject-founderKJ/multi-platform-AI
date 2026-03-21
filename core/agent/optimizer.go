//core\agent\optimizer.go

package agent

type Optimizer interface {
	Distill(input []byte) ([]byte, error)
}
