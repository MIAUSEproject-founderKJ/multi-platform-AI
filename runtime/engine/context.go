//runtime/engine/context.go

package runtime_engine

import (
	"context"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	runtime_supervisor "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/supervisor"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

type RuntimeContext struct {
	Execution runtime_types.ExecutionContext
	User      *user_setting.UserSession
	Policy    *RuntimePolicy

	Router     router.Router
	Bus        *runtime_bus.MessageBus
	Supervisor *runtime_supervisor.Supervisor

	Modules map[string]runtime_supervisor.Module
	Ctx     context.Context
}
