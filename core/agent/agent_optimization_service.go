//core/agent/agent_optimization_service.go

package core_agent

type Optimizer interface {
	Distill(input []byte) ([]byte, error)
}
