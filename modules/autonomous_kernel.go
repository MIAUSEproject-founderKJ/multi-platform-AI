//go:build automotive


//modules/autonomous_kernel.go 
//Domain Module. If this binary runs on a laptop: Capabilities won't match → module never loads.

package modules

import "aios/core"

type AutonomousKernel struct {
	ctx core.RuntimeContext
}

func (a *AutonomousKernel) Name() string {
	return "AutonomousKernel"
}

func (a *AutonomousKernel) RequiredCapabilities() []core.Capability {
	return []core.Capability{
		core.CapSafetyCritical,
		core.CapCANBus,
		core.CapHighFreqSensor,
	}
}

func (a *AutonomousKernel) RequiredPermissions() []string {
	return []string{"AUTONOMOUS_EXECUTION"}
}

func (a *AutonomousKernel) Init(ctx core.RuntimeContext) error {
	a.ctx = ctx
	return nil
}

func (a *AutonomousKernel) Start() error {
	// enter control loop
	return nil
}

func (a *AutonomousKernel) Stop() error {
	return nil
}
