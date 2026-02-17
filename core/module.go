//core/module.go
//No module can access global state. Everything flows through RuntimeContext.

package core

type Module interface {
	Name() string
	RequiredCapabilities() []Capability
	RequiredPermissions() []string
	Init(ctx RuntimeContext) error
	Start() error
	Stop() error
}
