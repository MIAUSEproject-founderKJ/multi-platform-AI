// core/kernel.go

package core

import (
    "context"
    "fmt"
    
    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform"
    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
    "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

// The Kernel is now strictly a coordinator. 
// It knows NOTHING about Vision, Lidar, or Motors.
type Kernel struct {
    Ctx      context.Context
    cancel   context.CancelFunc
    
    // The "Source of Truth"
    Runtime  RuntimeContext
    
    // The Subsystem Manager
    Loader   *Loader
    
    // Core Infrastructure
    Vault    *security.IsolatedVault
    EventBus *EventBus // Modules communicate via this
}

func Bootstrap(ctx context.Context) (*Kernel, error) {
    logging.Info("[KERNEL] Bootstrapping Nucleus...")

    // 1. Hardware/Platform Discovery (Your existing solid logic)
    v, err := security.OpenVault()
    if err != nil { return nil, err }
    
    bootSeq, err := platform.RunBootSequence(v)
    if err != nil { return nil, err }

    // 2. THE BRIDGE: Create the Context
    runtimeCtx := ResolveRuntimeContext(bootSeq.EnvConfig, bootSeq.Mode)
    
    // 3. Initialize Module Infrastructure
    registry := NewRegistry()
    loader := NewLoader(registry)
    
    // 4. Register Available Modules 
    // (In a real app, these might be dynamically imported or use plugins)
    // registry.Register(&modules.AutonomousKernel{}) 
    // registry.Register(&modules.SimEngine{})
    
    // 5. Load Modules based on Capabilities
    // This is where "SimEngine" gets dropped if we are on an embedded board
    if err := loader.ResolveAndLoad(runtimeCtx); err != nil {
        return nil, fmt.Errorf("module resolution failed: %w", err)
    }

    // 6. Create the Kernel
    kCtx, cancel := context.WithCancel(ctx)
    k := &Kernel{
        Ctx:      kCtx,
        cancel:   cancel,
        Runtime:  runtimeCtx,
        Loader:   loader,
        Vault:    v,
        EventBus: NewEventBus(),
    }

    // 7. Ignite
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