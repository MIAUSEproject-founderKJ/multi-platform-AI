// modules\domain\cognition\cognition_module.go
package module_cognition

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/math_convert"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
	domain_shared "github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules/domain/shared"
	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	runtime_engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
)

// CognitionModule handles incoming intents and converts them to tasks.
type CognitionModule struct {
	ctx *runtime_engine.RuntimeContext
}

// NewCognitionModule returns a DomainModule instance
func NewCognitionModule() domain_shared.DomainModule {
	return &CognitionModule{}
}

// SetRuntime sets the runtime context directly
func (m *CognitionModule) SetRuntime(rtx *runtime_engine.RuntimeContext) {
	m.ctx = rtx
}

// Init subscribes to the "audio.intent" topic
func (m *CognitionModule) Init(ctx *internal_boot.BootContext) error {
	if m.ctx == nil {
		ctx.Logger.Error("runtime context not set")
		return errors.New("runtime context not set")
	}
	m.ctx.Bus.Subscribe("audio.intent")
	return nil
}

// Handle processes incoming payloads and publishes tasks
func (m *CognitionModule) Handle(ctx context.Context, payload []byte) error {
	intent := parseIntent(payload)
	if intent.Confidence < math_convert.FromFloat64(0.75) {
		return nil
	}

	task := plan(intent)
	data, _ := json.Marshal(task)

	// Return the error from Publish
	m.ctx.Bus.Publish(runtime_bus.Message{
		Topic: "cognition.task",
		Data:  data,
	})
	return nil
}

// Run keeps the module alive until context cancellation
func (m *CognitionModule) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// DomainModule interface methods
func (m *CognitionModule) Name() string                                             { return "CognitionModule" }
func (m *CognitionModule) Category() ModuleCategory                                 { return ModuleDomain }
func (m *CognitionModule) DependsOn() []string                                      { return []string{"AudioModule"} }
func (m *CognitionModule) Allowed(ctx *internal_boot.BootContext) bool              { return true }
func (m *CognitionModule) Start() error                                             { return nil }
func (m *CognitionModule) Stop() error                                              { return nil }
func (m *CognitionModule) SupportedPlatforms() []internal_environment.PlatformClass { return nil }

// DomainModule implementation
func (m *CognitionModule) RequiredCapabilities() internal_environment.CapabilitySet {
	// This module doesn’t require any capabilities, so return 0
	return 0
}
func (m *CognitionModule) Optional() bool { return false }

// parseIntent converts payload to Intent (placeholder logic)
func parseIntent(payload []byte) *Intent {
	return &Intent{
		Name:       string(payload),
		Confidence: 1.0, // placeholder
	}
}

// plan converts an Intent to a Task (placeholder logic)
func plan(intent *Intent) *Task {
	return &Task{
		Action: intent.Name,
	}
}
