//cmd/aios/main.go
package main

import (
	"fmt"
	"os"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/classify"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func main() {
    ctx := context.Background()

    k, err := core.Bootstrap(ctx)
    if err != nil {
        logging.Error("Bootstrap failed: %v", err)
        os.Exit(1)
    }

    k.Run()
}

