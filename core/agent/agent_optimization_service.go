//core/agent/agent_optimization_service.go

package agent

type Optimizer interface {
	Distill(input []byte) ([]byte, error)
}
