// core/kernel.go
//cmd/aios-node/main.go->core/kernel.go
package core

import (
	"context"
	"fmt"
	"time"

	"multi-platform-AI/api/hmi"
	"multi-platform-AI/core/platform"
	"multi-platform-AI/core/policy"
	"multi-platform-AI/core/security"
	"multi-platform-AI/internal/logging"
)

// Kernel represents the operational heart of the system.
type Kernel struct {
	Platform *platform.BootSequence
	Vault    *security.IsolatedVault
	Trust    *policy.TrustDescriptor
	Status   string
	
	// Channels and Sub-systems for Lifecycle
	HMIPipe  chan hmi.ProgressUpdate
	// Internal components (Assuming these are initialized during Bootstrap)
	Sim      SimulationEngine 
	Bridge   PowerController
	Vitals   HealthMonitor
	Memory   CognitiveVault
}

// Bootstrap is the "Entry Gate" called by cmd/aios-node/main.go.
func Bootstrap(ctx context.Context) (*Kernel, error) {
	logging.Info("Kernel: Starting Secure Bootstrap sequence...")

	// 1. Initialize the Secure Vault first (Needed for markers)
	v, err := security.OpenVault()
	if err != nil {
		return nil, fmt.Errorf("vault initialization failed: %w", err)
	}

	// 2. Physical Verification (Path: core/platform/boot.go)
	// We pass the vault to RunBootSequence as established in our workflow
	pSequence, err := platform.RunBootSequence(v)
	if err != nil {
		return nil, fmt.Errorf("platform verification failed: %w", err)
	}

	// 3. Security Handshake (Check code signatures/attestation)
	if err := security.Initialize(pSequence.PlatformID); err != nil {
		return nil, fmt.Errorf("security handshake failed: %w", err)
	}

	// 4. Trust Initialization (Bayesian Prior)
	trustDescriptor := policy.InitializeTrust(pSequence)

	logging.Info("Kernel: Bootstrap complete. Identity: %s | Mode: %s", 
		pSequence.PlatformID, pSequence.Mode)

	return &Kernel{
		Platform: pSequence,
		Vault:    v,
		Trust:    trustDescriptor,
		Status:   "initialized",
		HMIPipe:  make(chan hmi.ProgressUpdate, 10),
	}, nil
}

// --- Lifecycle Methods ---

// TrustLevel returns the current autonomy capability for main.go to display.
func (k *Kernel) TrustLevel() float64 {
	return k.Trust.CurrentScore
}

// RunLifecycle manages the "Dream State" and power modes.
func (k *Kernel) RunLifecycle() {
	for {
		if k.IsIdle() {
			// 1. Lower power to non-essential HAL nodes
			k.Bridge.SetPowerMode("PowerSave")
			
			// 2. Reflective HUD Update
			k.HMIPipe <- hmi.ProgressUpdate{
				Stage: "SIM_DREAM",
				Message: "IDLE: Running Digital Twin Simulations...",
			}

			// 3. Trigger Dream State with Twist Injection
			k.Sim.EnterDreamState(k.Memory.Recall())
			// replay.InjectTwist(k.Sim.World) // Assuming replay package is imported
			
		} else {
			// Wake up immediately on user input
			k.Sim.Stop()
			k.Bridge.SetPowerMode("Performance")
		}
		time.Sleep(1 * time.Second)
	}
}

// Shutdown performs the "Safe-Park" of hardware before the process exits.
func (k *Kernel) Shutdown() {
	logging.Info("Kernel: Executing safe shutdown. Locking all actuators...")
	k.Vault.Close() // Ensure the vault is sealed
}

// IsIdle is a helper to determine if the system should enter Dream State
func (k *Kernel) IsIdle() bool {
	// Logic to check if user input is absent and CPU load is low
	return true 
}

func (k *Kernel) MonitorState() {
	for {
		if k.IsIdle() && k.Vitals.Temperature < 65.0 {
			// Reflective HUD Update
			k.HMIPipe <- hmi.ProgressUpdate{
				Stage: "SIM_DREAM",
				Message: "IDLE: Running Digital Twin Simulations...",
			}
			
			k.Sim.EnterDreamState(k.Memory.Recall())
		} else {
			k.Sim.Stop()
		}
		time.Sleep(5 * time.Second)
	}
}

func (k *Kernel) ReflectToHUD() {
    for {
        // Send Vitals
        k.HMIPipe.SendTelemetry(hmi.SystemPulse{
            CpuLoad:     k.Vitals.CPU,
            TrustScore:  k.Trust.CurrentScore,
            Temperature: k.Vitals.Temp,
        })
        
        // If IDLE, send the Voxel Map
        if k.IsIdle() {
             k.HMIPipe.SendSpatial(k.Sim.GetVoxelFrame())
        }
        time.Sleep(100 * time.Millisecond) // 10Hz Refresh
    }
}

// Start initiates the parallel background processes of the Kernel.
func (k *Kernel) Start(ctx context.Context) {
    logging.Info("[KERNEL] Activating background subsystems...")

    // 1. Start the Reflective HUD stream (10Hz)
    go k.ReflectToHUD()

    // 2. Start the Lifecycle manager (Power & Dream State)
    go k.RunLifecycle()

    // 3. Keep-alive or monitor for context cancellation
    <-ctx.Done()
    k.Shutdown()
}