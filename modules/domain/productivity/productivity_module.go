//go:build workstation

//modules/domain/productivity/productivity_module.go
//Productivity Module (Workstation)

package module_productivity

import (
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
)

type ProductivityModule struct {
	ctx engine.RuntimeContext
}

func (p *ProductivityModule) Name() string {
	return "ProductivityModule"
}

func (p *ProductivityModule) RequiredPermissions() []string {
	return []string{"STANDARD_USE"}
}

func (p *ProductivityModule) Init(ctx internal_boot.BootContext) error {
	p.ctx = ctx
	return nil
}

func (p *ProductivityModule) Start() error {
	return nil
}

func (p *ProductivityModule) Stop() error {
	return nil
}
