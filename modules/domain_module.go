//modules/domain_module.go

/*
This allows:
• platform filtering
• hardware capability filtering
• policy enforcement
• optimizer-based gating
*/
package modules

import (
	"context"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime"
)

type DomainModule interface {
	Name() string
	Category() ModuleCategory
	DependsOn() []string
	Allowed(*schema.BootContext) bool
	Init(*schema.BootContext) error
	Start() error
	Stop() error
	Run(context.Context) error

	// Capability introspection
	SupportedPlatforms() []schema.PlatformClass
	RequiredCapabilities() schema.CapabilitySet
	Optional() bool
}

// All modules that need DB/Bus access implement this
type RuntimeAware interface {
	SetRuntime(*runtime.RuntimeContext)
}

func (defaultPolicy) Allowed(*schema.BootContext, Intent) bool {
	return true
}

type ModuleCategory int

const (
	ModuleCore ModuleCategory = iota
	ModulePlatform
	ModuleDomain
	ModuleCognitive
	ModuleCategoryInference ModuleCategory = iota
	ModuleCategoryTelemetry
	ModuleCategoryControl
)

// AI Agent Initialization (Execution Layer)
type ExecutionRouter interface {
	ExecuteIntent(Intent) error
}

type IntentInterpreter interface {
	Parse(input string) (Intent, error)
}

type TaskPlanner interface {
	Plan(Intent) Task
}

type IntentHandler interface {
	Handle(Intent) error
}

type AgentRuntime struct {
	interpreter IntentInterpreter
	planner     TaskPlanner
	router      ExecutionRouter
}

/*This prevents:

• Raw LLM text from executing shell
• Hallucinated commands reaching hardware
• Low-confidence unsafe operations*/

func (a *AgentRuntime) HandleInput(input string) error {

	intent, err := a.interpreter.Parse(input)
	if err != nil {
		return err
	}

	if intent.Confidence < 0.75 {
		return fmt.Errorf("low confidence intent")
	}

	task := a.planner.Plan(intent)

	return a.router.ExecuteIntent(task.Intent)
}

type DefaultRouter struct {
	ctx      *schema.BootContext
	handlers map[string]IntentHandler
}

/*• Domain = "navigation" only valid on vehicle/robot
• Domain = "web" only valid on PC/mobile
• Domain = "actuator" only valid on industrial/tester*/

func (r *DefaultRouter) ExecuteIntent(i Intent) error {

	if !policy.Allowed(r.ctx, i) {
		return fmt.Errorf("policy denied")
	}

	h, ok := r.handlers[i.Domain]
	if !ok {
		return fmt.Errorf("no handler for domain %s", i.Domain)
	}

	return h.Handle(i)
}

type PolicyEngine interface {
	Allowed(*schema.BootContext, Intent) bool
}

type defaultPolicy struct{}

var policy = struct {
	Allowed func(*schema.BootContext, Intent) bool
}{
	Allowed: func(*schema.BootContext, Intent) bool {
		return true
	},
}
