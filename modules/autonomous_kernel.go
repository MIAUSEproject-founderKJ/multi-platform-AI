//go:build automotive

//modules/autonomous_kernel.go
//Domain Module. If this binary runs on a laptop: Capabilities won't match → module never loads.

package modules

type AutonomousKernel struct {
	ctx schema_boot.BootContext
}

func (a *AutonomousKernel) Name() string {
	return "AutonomousKernel"
}

func (a *AutonomousKernel) RequiredPermissions() []string {
	return []string{"AUTONOMOUS_EXECUTION"}
}

func (a *AutonomousKernel) Init(ctx schema_boot.BootContext) error {
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
