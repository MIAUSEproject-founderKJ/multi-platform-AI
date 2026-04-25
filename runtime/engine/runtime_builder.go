// runtime/engine/runtime_builder.go

package engine

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"

	boot_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/phases"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	runtime "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	"go.uber.org/zap"
)

type RuntimeContext struct {
	Router       router.Router
	Bus          *runtime_bus.MessageBus
	DB           *sql.DB
	Logger       *zap.Logger
	BasePath     string
	Session      *schema_identity.UserSession
	Orchestrator *boot_phase.Orchestrator
	Config       *schema_identity.UserConfig
	Context      context.Context
}

func (r *RuntimeContext) SafePath(rel string) string {

	if r.BasePath == "" {
		r.BasePath = "./data" // fallback
	}

	// normalize path
	clean := filepath.Clean(rel)

	// prevent directory traversal
	if strings.Contains(clean, "..") {
		clean = strings.ReplaceAll(clean, "..", "")
	}

	return filepath.Join(r.BasePath, clean)
}

///////////////////////////////////////////////////////////////
// CONSTRUCTOR
///////////////////////////////////////////////////////////////

func NewRuntimeContext(logger *zap.Logger) (*RuntimeContext, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	rtx := &RuntimeContext{
		Bus:    runtime.NewMessageBus(),
		Logger: logger,
	}

	// Optional: initialize router if required
	rtx.Router = router.New() // <-- ONLY if this exists

	return rtx, nil
}
