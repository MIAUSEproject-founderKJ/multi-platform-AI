//runtime/runtime_context.go

package runtime

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"go.uber.org/zap"
)

type RuntimeContext struct {
	Router   router.Router
	Bus      *MessageBus
	DB       *sql.DB
	Logger   *zap.Logger
	BasePath string
	Session      *schema.UserSession
	Orchestrator *interaction.Orchestrator
	Config   *schema.UserConfig
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
		Bus:    NewMessageBus(),
		Logger: logger,
	}

	// Optional: initialize router if required
	rtx.Router = router.New() // <-- ONLY if this exists

	return rtx, nil
}
