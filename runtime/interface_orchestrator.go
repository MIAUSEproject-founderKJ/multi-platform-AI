//runtime/interface_orchestrator.go


package runtimectx

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type InterfaceAdapter interface {
	Start() error
	Notify(msg string)
}

type Orchestrator struct {
	adapters []InterfaceAdapter
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{}
}

func (o *Orchestrator) Add(adapter InterfaceAdapter) {
	o.adapters = append(o.adapters, adapter)
}

func (o *Orchestrator) StartAll() {
	for _, a := range o.adapters {
		go a.Start()
	}
}

func (o *Orchestrator) Broadcast(msg string) {
	for _, a := range o.adapters {
		a.Notify(msg)
	}
}

func BuildOrchestrator(cp *schema.CapabilityProfile) *Orchestrator {
	orch := NewOrchestrator()

	// Screen interface
	if cp.IsHealthy(schema.CapDisplay) {
		orch.Add(NewScreenAdapter())
	}

	// Voice interface (requires both mic + speaker)
	if cp.IsHealthy(schema.CapMicrophone) &&
		cp.IsHealthy(schema.CapSpeaker) {
		orch.Add(NewVoiceAdapter())
	}

	return orch
}

//===================Screen Adapter
type ScreenAdapter struct{}

func NewScreenAdapter() *ScreenAdapter {
	return &ScreenAdapter{}
}

func (s *ScreenAdapter) Start() error {
	fmt.Println("[UI] Screen interface started")
	return nil
}

func (s *ScreenAdapter) Notify(msg string) {
	fmt.Println("[SCREEN]", msg)
}

//============VoiceAdapter
type VoiceAdapter struct {
	engine *VoiceEngine
}

func NewVoiceAdapter() *VoiceAdapter {
	return &VoiceAdapter{
		engine: NewVoiceEngine(),
	}
}

func (v *VoiceAdapter) Start() error {
	fmt.Println("[VOICE] Full-duplex voice engine started")
	v.engine.Start()
	return nil
}

func (v *VoiceAdapter) Notify(msg string) {
	v.engine.outputChan <- msg
}
