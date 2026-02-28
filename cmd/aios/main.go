// cmd/aios/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
	boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
)

func main() {

	// ------------------------------------------------------------
	// 1. Root Context (graceful shutdown control)
	// ------------------------------------------------------------

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go listenForShutdown(cancel)

	// ------------------------------------------------------------
	// 2. Vault + Boot
	// ------------------------------------------------------------

	vault, err := security.OpenIsolatedVault()
	if err != nil {
		log.Fatalf("vault initialization failed: %v", err)
	}

	bootSeq, err := boot.RunBootSequence(vault)
	if err != nil {
		log.Fatalf("boot failed: %v", err)
	}

	ctx, err := runtime.ResolveExecutionContext(bootSeq)
	if err != nil {
		log.Fatalf("context resolution failed: %v", err)
	}

	// ------------------------------------------------------------
	// 3. Module Registry + Dependency Resolution
	// ------------------------------------------------------------
ctx.Optimizer = optimization.NewDefaultOptimizer(ctx.PlatformClass)
/*Example behavior:
• Vehicle → aggressive pruning
• Industrial → deterministic inference mode
• PC → full precision*/

registry := module.DefaultRegistry()

filtered := module.FilterModules(registry, ctx)

ordered, err := module.ResolveDependencies(filtered)
if err != nil {
    log.Fatalf("dependency resolution failed: %v", err)
}

active := []module.DomainModule{}

for _, m := range ordered {

    if err := m.Init(ctx); err != nil {
        rollbackModules(active)
        log.Fatalf("init failed: %v", err)
    }

    if err := m.Start(); err != nil {
        rollbackModules(active)
        log.Fatalf("start failed: %v", err)
    }

    active = append(active, m)

		log.Printf("[MODULE] Activated: %s", m.Name())
	}

	// ------------------------------------------------------------
	// 5. Session
	// ------------------------------------------------------------

	router := NewDefaultRouter(ctx)
agent := NewAgentRuntime(router)

session := runtime.NewSession(ctx, agent)

	if err := session.Start(); err != nil {
		rollbackModules(active)
		log.Fatalf("session start failed: %v", err)
	}

	go operationalLoop(rootCtx, session)

	// ------------------------------------------------------------
	// 6. Block Until Shutdown
	// ------------------------------------------------------------

	<-rootCtx.Done()

	// ------------------------------------------------------------
	// 7. Graceful Shutdown (Reverse Order)
	// ------------------------------------------------------------

	if err := session.Stop(); err != nil {
		log.Printf("session stop error: %v", err)
	}

	for i := len(active) - 1; i >= 0; i-- {
		if err := active[i].Stop(); err != nil {
			log.Printf("module %s stop error: %v", active[i].Name(), err)
		}
	}
}

func rollbackModules(active []modules.DomainModule) {
	for i := len(active) - 1; i >= 0; i-- {
		_ = active[i].Stop()
	}
}

func listenForShutdown(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	cancel()
}