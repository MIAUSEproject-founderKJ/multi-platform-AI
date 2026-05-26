//runtime/engine/runtime_builder.go

package runtime_engine

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
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

	SafetyMode SafetyMode
}

type AppContainer struct {
	Logger *zap.Logger
	DB     *sql.DB
	Bus    *runtime_bus.MessageBus
	Router router.Router
}

type RuntimeContainer struct {
	Policy *RuntimePolicy
	Infra  *AppContainer
	Ctx    context.Context
}

func (r *RuntimeContainer) SafePath(rel string) (string, error) {
	base := filepath.Clean(r.Infra.BasePath)

	target := filepath.Join(base, rel)
	target = filepath.Clean(target)

	if !strings.HasPrefix(target, base) {
		return "", errors.New("path traversal detected")
	}

	return target, nil
}

///////////////////////////////////////////////////////////////
// CONSTRUCTOR
///////////////////////////////////////////////////////////////

func Build(
	exec runtime_types.ExecutionContext,
	user *user_setting.UserSession,
	logger *zap.Logger,
	db *sql.DB,
	rt router.Router,
) (*RuntimeContainer, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	rtx := &AppContainer{
		Bus:    runtime.NewMessageBus(),
		Logger: logger,
		DB:     db,
		Router: rt,
	}

	return &RuntimeContainer{
		Infra: rtx,
	}, nil
}

func (app *RuntimeContainer) Start(ctx context.Context) error
