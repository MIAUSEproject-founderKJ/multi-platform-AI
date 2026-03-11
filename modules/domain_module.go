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
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/core"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/runtime"
)

type DomainModule interface {
	Name() string
	Category() ModuleCategory
	DependsOn() []string
	Allowed(*runtime.RuntimeContext) bool
	Init(*runtime.RuntimeContext) error
	Start() error
	Stop() error
	Run(context.Context) error

	// Capability introspection
	SupportedPlatforms() []runtime.PlatformClass
	RequiredCapabilities() []string
	Optional() bool
}

type ModuleCategory int

const (
	ModuleCore ModuleCategory = iota
	ModulePlatform
	ModuleDomain
	ModuleCognitive
)

// AI Agent Initialization (Execution Layer)
type ExecutionRouter interface {
	ExecuteIntent(Intent) error
}

type Intent struct {
	Domain     string
	Action     string
	Parameters map[string]interface{}
	Confidence float64
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
	ctx      *core.RuntimeContext
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
