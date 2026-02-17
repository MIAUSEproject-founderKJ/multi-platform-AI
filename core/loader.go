//core/loader.go
//Module Loader (Kernel Orchestrator)
/*
No platform checks.
No entity branching.
Only capability + permission resolution.
*/

package core

type Loader struct {
	registry *ModuleRegistry
	active   []Module
}

func NewLoader(registry *ModuleRegistry) *Loader {
	return &Loader{
		registry: registry,
		active:   []Module{},
	}
}

func (l *Loader) ResolveAndLoad(ctx RuntimeContext) error {
	for _, m := range l.registry.All() {

		if !hasCapabilities(ctx, m.RequiredCapabilities()) {
			continue
		}

		if !hasPermissions(ctx, m.RequiredPermissions()) {
			continue
		}

		if err := m.Init(ctx); err != nil {
			return err
		}

		l.active = append(l.active, m)
	}

	return nil
}

func (l *Loader) StartAll() error {
	for _, m := range l.active {
		if err := m.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Loader) StopAll() {
	for _, m := range l.active {
		_ = m.Stop()
	}
}
