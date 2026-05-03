// runtime/interface_adapter/orchestrator.go
package interface_adapter

import (
	bootstrap_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/phases"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type Orchestrator struct {
	adapters []bootstrap_phase.InterfaceAdapter
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{}
}

func (orch *Orchestrator) Add(adapter bootstrap_phase.InterfaceAdapter) {
	orch.adapters = append(orch.adapters, adapter)
}

func (orch *Orchestrator) StartAll(s *user_setting.UserSession) {
	for _, adapt := range orch.adapters {
		go adapt.Start(s)
	}
}

func (orch *Orchestrator) Broadcast(msg string) {
	for _, adapt := range orch.adapters {
		adapt.Notify(msg)
	}
}
