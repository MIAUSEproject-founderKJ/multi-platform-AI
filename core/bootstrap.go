// core/bootstrap.go

package core

import (
	"context"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

func Bootstrap(ctx context.Context) (*Kernel, error) {

	vault, err := security.NewIsolatedVault()
	if err != nil {
		return nil, fmt.Errorf("vault initialization failed: %w", err)
	}

	bootSeq, err := platform.RunBootSequence(vault)
	if err != nil {
		return nil, err
	}

	runtimeCtx := RuntimeContext{
		PlatformClass: bootSeq.EnvConfig.PlatformClass,
		Mode:          bootSeq.Mode,
		Identity:      bootSeq.Identity,
	}

	registry := NewRegistry()
	loader := NewLoader(registry)

	if err := loader.ResolveAndLoad(runtimeCtx); err != nil {
		return nil, fmt.Errorf("module resolution failed: %w", err)
	}

	kCtx, cancel := context.WithCancel(ctx)

	k := &Kernel{
		Ctx:      kCtx,
		cancel:   cancel,
		Runtime:  runtimeCtx,
		Loader:   loader,
		Vault:    vault,
		EventBus: NewEventBus(),
	}

	if err := k.Loader.StartAll(); err != nil {
		k.Shutdown()
		return nil, fmt.Errorf("module ignition failed: %w", err)
	}

	return k, nil
}

func (k *Kernel) Shutdown() {
    logging.Info("[KERNEL] Initiating shutdown sequence...")
    k.Loader.StopAll()
    k.Vault.Close()
    k.cancel()
}