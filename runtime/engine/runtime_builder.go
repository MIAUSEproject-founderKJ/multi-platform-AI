// runtime/engine/runtime_builder.go
package runtime_engine

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	runtime_supervisor "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/supervisor"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"

	"go.uber.org/zap"
)

type RuntimePolicy struct {
	Permissions map[user_setting.PermissionKey]bool

	MaxInferenceRate     int
	MaxActuatorAuthority float64

	AllowExternalNetwork bool
	RequireHumanApproval bool

	AllowedModules map[string]bool
}

type AppContainer struct {
	Logger   *zap.Logger
	DB       *sql.DB
	Bus      *runtime_bus.MessageBus
	Router   router.Router
	BasePath string
}

type RuntimeContainer struct {
	Execution runtime_types.ExecutionContext
	User      *user_setting.UserSession
	Policy    *RuntimePolicy
	Infra     *AppContainer
	Ctx       context.Context
}

func (r *RuntimeContainer) SafePath(rel string) (string, error) {
	if r == nil || r.Infra == nil {
		return "", errors.New("runtime infrastructure is nil")
	}

	base := filepath.Clean(r.Infra.BasePath)

	if base == "." || base == "" {
		return "", errors.New("runtime base path is not configured")
	}

	if filepath.IsAbs(rel) {
		return "", errors.New("absolute paths are not allowed")
	}

	target := filepath.Clean(filepath.Join(base, rel))

	relative, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}

	if relative == ".." ||
		strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("path traversal detected")
	}

	return target, nil
}

func Build(
	exec runtime_types.ExecutionContext,
	user *user_setting.UserSession,
	logger *zap.Logger,
	db *sql.DB,
	rt router.Router,
) (*RuntimeContainer, error) {
	if exec == nil {
		return nil, errors.New("execution context is required")
	}

	if logger == nil {
		return nil, errors.New("logger is required")
	}

	bus := runtime_bus.NewMessageBus()

	infra := &AppContainer{
		Logger: logger,
		DB:     db,
		Bus:    bus,
		Router: rt,
	}

	policy := buildRuntimePolicy(exec)

	return &RuntimeContainer{
		Execution: exec,
		User:      user,
		Policy:    policy,
		Infra:     infra,
		Ctx:       context.Background(),
	}, nil
}

func buildRuntimePolicy(
	exec runtime_types.ExecutionContext,
) *RuntimePolicy {
	policy := &RuntimePolicy{
		Permissions:    make(map[user_setting.PermissionKey]bool),
		AllowedModules: make(map[string]bool),
	}

	// Keep policy construction based on the canonical execution
	// context. Do not reconstruct permissions from the user session here.
	for _, permission := range []user_setting.PermissionKey{
		user_setting.PermUser,
		user_setting.PermDiagnostics,
		user_setting.PermConfigEdit,
		user_setting.PermHardwareIO,
		user_setting.PermAdmin,
		user_setting.PermSafetyOverride,
	} {
		if exec.HasPermission(permission) {
			policy.Permissions[permission] = true
		}
	}

	return policy
}

func (r *RuntimeContainer) Context() *RuntimeContext {
	if r == nil || r.Infra == nil {
		return nil
	}

	return &RuntimeContext{
		Execution: r.Execution,
		User:      r.User,
		Policy:    r.Policy,

		Router: r.Infra.Router,
		Bus:    r.Infra.Bus,

		Ctx: r.Ctx,

		Modules: make(map[string]runtime_supervisor.Module),
	}
}
