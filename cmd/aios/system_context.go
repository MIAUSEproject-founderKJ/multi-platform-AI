//cmd/aios/system_context.go

package main

import (
	"fmt"
	"time"

	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	runtime_types "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/types"
)

// ============================================================
// SYSTEM CONTEXT
// ============================================================

type SystemContext struct {
	Boot    *runtime_types.ExecutionContext
	Exec    *runtime_types.ExecutionContext
	Session *user_setting.UserSession
}

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
