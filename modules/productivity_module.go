//go:build workstation


//modules/productivity_module.go
//Productivity Module (Workstation)

package modules

import "aios/core"

type ProductivityModule struct {
	ctx boot.RuntimeContext
}

func (p *ProductivityModule) Name() string {
	return "ProductivityModule"
}


func (p *ProductivityModule) RequiredPermissions() []string {
	return []string{"STANDARD_USE"}
}

func (p *ProductivityModule) Init(ctx boot.RuntimeContext) error {
	p.ctx = ctx
	return nil
}

func (p *ProductivityModule) Start() error {
	return nil
}

func (p *ProductivityModule) Stop() error {
	return nil
}
