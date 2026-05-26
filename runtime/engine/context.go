//runtime/engine/context.go

package runtime_engine

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	runtime_supervisor "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/supervisor"
)

type RuntimeContext struct {
	Router     router.Router
	Bus        *runtime_bus.MessageBus
	Supervisor *runtime_supervisor.Supervisor

	Modules map[string]runtime_supervisor.Module
}
