// cmd/aios/main.go
package main

import (
	"log"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/cmd/aios/runtime"
	boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
)

func main() {

	vault, err := security.OpenIsolatedVault()
	if err != nil {
		log.Fatalf("vault initialization failed: %v", err)
	}

	bootSeq, err := boot.RunBootSequence(vault)
	if err != nil {
		log.Fatalf("boot failed: %v", err)
	}

	ctx, err := runtime.ResolveExecutionContext(bootSeq, vault)
	if err != nil {
		log.Fatalf("context resolution failed: %v", err)
	}

	registry := module.DefaultRegistry()
	active := []module.DomainModule{}

	for _, m := range registry {

		if !m.Allowed(ctx) {
			continue
		}

		if err := m.Init(ctx); err != nil {
			log.Fatalf("module %s init failed: %v", m.Name(), err)
		}

		if err := m.Start(); err != nil {
			log.Fatalf("module %s start failed: %v", m.Name(), err)
		}

		active = append(active, m)
		log.Printf("[MODULE] Activated: %s", m.Name())
	}

	session := runtime.NewSession(ctx)
	if err := session.Start(); err != nil {
		log.Fatalf("session start failed: %v", err)
	}

	go operationalLoop(session)

	waitForShutdown()

	for _, m := range active {
		_ = m.Stop()
	}

	_ = session.Stop()
}
