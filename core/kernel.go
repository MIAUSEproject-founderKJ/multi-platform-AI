// core/kernel.go
// cmd/aios-node/main.go->core/kernel.go
package core

import (
	"context"
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/hmi"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/policy"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// --- Updated Interfaces to match usage ---

type SimulationEngine interface {
	InjectFault(env *schema.EnvConfig)
	EnterDreamState(data interface{})
	Stop()
	GetVoxelFrame() interface{}
}

type PowerController interface {
	TransitionTo(state string)
	SyncWithTrust(trust *policy.TrustDescriptor)
	SetPowerMode(mode string) // Added to fix "undefined" error
	WriteActuator(name string, value float64) error
}

type CognitiveVault interface {
	Store(id string, entry interface{})
	Recall(id string) (interface{}, bool)
}

// Added missing types for Step() and ReflectToHUD
type VisionSystem interface {
	ProcessFrame(frame interface{}) interface{}
}

type HardwareBridge interface {
	GetCameraFrame() interface{}
	GetPowerProfile() PowerProfile
}

type PowerProfile struct {
	BatteryLevel float64
}

type Vitals struct {
	CPU         float64
	Temperature float64
	Temp        float64 // Matching your inconsistent naming in ReflectToHUD
}

// --- The Kernel Struct ---

type Kernel struct {
	Platform  *platform.BootSequence
	Vault     *security.IsolatedVault
	Trust     *policy.TrustDescriptor
	EnvConfig *schema.EnvConfig
	Evaluator *policy.TrustEvaluator

	// Subsystems
	Sim      SimulationEngine
	Bridge   PowerController
	Memory   CognitiveVault
	Vision   VisionSystem
	Hardware HardwareBridge
	Vitals   Vitals // Fixed "k.Vitals undefined"

	// Communication
	HMIPipe hmi.HMIPipe // Changed from chan to Interface
	Status  string
	ctx     context.Context
}

// Bootstrap is the "Entry Gate" called by cmd/aios-node/main.go.
func Bootstrap(ctx context.Context) (*Kernel, error) {
	logging.Info("Kernel: Starting Secure Bootstrap sequence...")

	// 1. Initialize the Secure Vault first (Needed for markers)
	/* OpenVault return
	type IsolatedVault struct {
	BaseDir string
	Key     []byte // Reserved for future AES-GCM encryption implementation}
	*/
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
		// FIX: Use a constructor for the interface instead of make(chan)
		HMIPipe: hmi.NewBufferedPipe(10),
	}, nil
}

// --- Lifecycle Methods ---

/*calculating the "Truth" every time the function is called. calling this in a high-frequency loop could stutter the CPU.
func (k *Kernel) TrustLevel() float64 {
	// Quick-access for the main.go logger
	res := k.Evaluator.Evaluate(k.EnvConfig)
	return res.CurrentScore
}
*/

// If the HMILoop hasn't refreshed the k.Trust pointer recently, this value might be "lying" about the current state of the hardware.
func (k *Kernel) TrustLevel() float64 {
	res := k.Evaluator.Evaluate(k.EnvConfig)
	return res.CurrentScore
}

// RunLifecycle manages the "Dream State" and power modes.
func (k *Kernel) RunLifecycle() {
	for {
		if k.IsIdle() {
			k.Bridge.SetPowerMode("PowerSave")

			// Fix: Recall needs a key string
			dreamData, _ := k.Memory.Recall("last_world_state")
			k.Sim.EnterDreamState(dreamData)

			// MonitorProgress stages should match your hmi provider
			k.HMIPipe.SendProgress(hmi.MonitorProgress{
				Task:     "SIM_DREAM",
				Progress: 0.5,
			})
		} else {
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
	// 1. Check if CPU load is below a "quiet" threshold (e.g., 15%)
	isQuiet := k.Vitals.CPU < 0.15

	// 2. You could also check a "LastCommandTime" field (optional)
	// isRecentlyActive := time.Since(k.LastCommandTime) < 5 * time.Second

	return isQuiet
}

func (k *Kernel) MonitorState() {
	for {
		if k.IsIdle() && k.Vitals.Temperature < 65.0 {
			// Reflective HUD Update
			k.HMIPipe <- hmi.MonitorProgress{
				Stage:   "SIM_DREAM",
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
	ticker := time.NewTicker(100 * time.Millisecond)
	for range ticker.C {
		// Fix: SendTelemetry now exists on the interface
		k.HMIPipe.SendTelemetry(hmi.SystemPulse{
			CPUUsage:  k.Vitals.CPU,
			Status:    k.Status,
			Timestamp: time.Now(),
		})

		if k.IsIdle() {
			k.HMIPipe.SendSpatial(k.Sim.GetVoxelFrame())
		}
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

func (k *Kernel) OnPeerUpdate(pulse NodePulse) {
	// 1. Verify Peer Attestation
	if !k.Vault.VerifyRemote(pulse.SourceID, pulse.Trust.EnvHash) {
		logging.Warn("[SECURITY] Rejected untrusted peer pulse from %s", pulse.SourceID)
		return
	}

	// 2. Update Global Map
	// This data can be recalled by the CognitiveVault for swarm-level learning
}

func (k *Kernel) Step() {
	// 1. See: Vision processing
	objects := k.Vision.ProcessFrame(k.Hardware.GetCameraFrame())

	// 2. Feel: Energy check
	p := k.Hardware.GetPowerProfile()

	// 3. Think: Bayesian Update
	k.Trust = k.Evaluator.Evaluate(k.EnvConfig, objects, p)

	// 4. Socialize: Negotiate with Swarm
	if k.Trust.CurrentScore < 0.5 || p.BatteryLevel < 0.1 {
		k.Swarm.RequestHelp("REPLACEMENT_NEEDED")
	}
}
