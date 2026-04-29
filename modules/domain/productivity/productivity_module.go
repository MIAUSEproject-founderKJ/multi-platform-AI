//go:build workstation

//modules/domain/productivity/productivity_module.go
//Productivity Module (Workstation)

package module_productivity

import (
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	runtime_engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
)

type ProductivityModule struct {
	ctx runtime_engine.RuntimeContext
}

func (p *ProductivityModule) Name() string {
	return "ProductivityModule"
}

func (p *ProductivityModule) RequiredPermissions() []string {
	return []string{"STANDARD_USE"}
}

func (p *ProductivityModule) Init(ctx bootstrap.BootContext) error {
	p.ctx = ctx
	return nil
}

func (p *ProductivityModule) Start() error {
	return nil
}

func (p *ProductivityModule) Stop() error {
	return nil
}
