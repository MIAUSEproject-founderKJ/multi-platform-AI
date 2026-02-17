//core/registry.go
//Domain modules register themselves during build or init().


package core

type ModuleRegistry struct {
	modules []Module
}

func NewRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: []Module{},
	}
}

func (r *ModuleRegistry) Register(m Module) {
	r.modules = append(r.modules, m)
}

func (r *ModuleRegistry) All() []Module {
	return r.modules
}
