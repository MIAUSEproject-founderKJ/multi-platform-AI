//cmd/aios/system_context.go

package main

import (
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

// ============================================================
// SYSTEM CONTEXT
// ============================================================

type SystemContext struct {
	Boot      *bootstrap.BootContext
	Execution runtime_types.ExecutionContext
	Session   *user_setting.UserSession
}

// buildSystemContext attempts to build the system context by running the boot sequence and resolving the execution context. It includes retry logic to handle transient failures during boot. The function returns a fully constructed SystemContext or an error if boot fails after retries.
func buildSystemContext() (*SystemContext, error) {
	var last error

	for i := 0; i < 3; i++ {
		sys, err := attemptBoot()
		if err == nil {
			return sys, nil
		}
		last = err
		time.Sleep(time.Second * time.Duration(i+1))
	}

	return nil, fmt.Errorf("boot failed: %w", last)
}
