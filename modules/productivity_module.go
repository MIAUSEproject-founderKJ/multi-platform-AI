//go:build workstation

//modules/productivity_module.go
//Productivity Module (Workstation)

package modules

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

type ProductivityModule struct {
	ctx schema.RuntimeContext
}

func (p *ProductivityModule) Name() string {
	return "ProductivityModule"
}

func (p *ProductivityModule) RequiredPermissions() []string {
	return []string{"STANDARD_USE"}
}

func (p *ProductivityModule) Init(ctx schema.BootContext) error {
	p.ctx = ctx
	return nil
}

func (p *ProductivityModule) Start() error {
	return nil
}

func (p *ProductivityModule) Stop() error {
	return nil
}
