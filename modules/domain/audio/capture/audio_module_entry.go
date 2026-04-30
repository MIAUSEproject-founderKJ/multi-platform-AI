// modules/domain/audio/capture/audio_module_entry.go
package audio_capture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/modules"
	runtime_bus "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/bus"
	runtime_engine "github.com/MIAUSEproject-founderKJ/multi-platform-AI/runtime/engine"
	"go.uber.org/zap"
)

var (
	ErrRuntimeNotSet    = errors.New("runtime context not set")
	ErrModuleNotRunning = errors.New("module not running")
)

// --------------------------------------------------
// AudioModule
// --------------------------------------------------

type AudioModule struct {
	kernel_lifecycle.BaseModule

	runtime   *runtime_engine.RuntimeContext
	writer    *WAVWriter
	extractor *FeatureExtractor
	repo      *AudioRepository

	running atomic.Bool
}

// --------------------------------------------------
// Constructor
// --------------------------------------------------

func NewAudioModule() kernel_lifecycle.DomainModule {
	m := &AudioModule{}
	m.SetName("AudioModule")
	m.SetDeps([]string{"StorageModule"})
	return m
}

// --------------------------------------------------
// Runtime Injection
// --------------------------------------------------

func (m *AudioModule) SetRuntime(rtx *runtime_engine.RuntimeContext) {
	m.runtime = rtx
}

// --------------------------------------------------
// Init (no goroutines here)
// --------------------------------------------------

func (m *AudioModule) Init(ctx *bootstrap.BootContext) error {
	if m.runtime == nil {
		return ErrRuntimeNotSet
	}

	// Ensure storage dependency exists
	if m.runtime.SafePath == nil {
		return fmt.Errorf("storage capability not available")
	}

	outputPath := m.runtime.SafePath("audio/audio_output.wav")

	m.writer = NewWAVWriter(outputPath)
	m.extractor = NewFeatureExtractor(16000)
	m.repo = NewAudioRepository()

	return nil
}

// --------------------------------------------------
// Run (bus-driven processing loop)
// --------------------------------------------------

func (m *AudioModule) Run(ctx context.Context) error {
	if m.runtime == nil {
		return ErrRuntimeNotSet
	}

	ch := m.runtime.Bus.Subscribe("audio.raw")

	m.running.Store(true)
	defer m.running.Store(false)

	for {
		select {

		case payload, ok := <-ch:
			if !ok {
				return nil
			}

			if err := m.process(payload.Data); err != nil {
				if m.runtime.Logger != nil {
					m.runtime.Logger.Error("audio processing failed", zap.Error(err))
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}

// --------------------------------------------------
// Core Processing Pipeline
// --------------------------------------------------

func (m *AudioModule) process(payload []byte) error {

	// 1. Persist raw audio
	if err := m.writer.AppendPCM(payload); err != nil {
		return err
	}

	// 2. Extract features
	features, err := m.extractor.ProcessPCM(payload)
	if err != nil {
		return err
	}

	// 3. Serialize features
	data, err := json.Marshal(features)
	if err != nil {
		return err
	}

	// 4. Publish downstream
	msg := runtime_bus.Message{
		Topic: "audio.features",
		Data:  data,
	}

	m.runtime.Bus.Publish(msg)

	return nil
}

// --------------------------------------------------
// Capability Enforcement
// --------------------------------------------------

func (m *AudioModule) RequiredCapabilities() internal_environment.CapabilitySet {
	return internal_environment.CapLocalStorage | internal_environment.CapNetwork
}

// --------------------------------------------------
// Policy Enforcement (CRITICAL)
// --------------------------------------------------

func (m *AudioModule) Allowed(ctx *bootstrap.BootContext) bool {

	// Must have runtime execution rights
	if !ctx.Permissions[user_setting.PermBasicRuntime] {
		return false
	}

	// Audio capture requires hardware permission
	if !ctx.Permissions[user_setting.PermHardwareIO] {
		return false
	}

	return true
}

// --------------------------------------------------
// Lifecycle
// --------------------------------------------------

func (m *AudioModule) Start() error {
	m.running.Store(true)
	return nil
}

func (m *AudioModule) Stop() error {
	m.running.Store(false)
	return nil
}

// --------------------------------------------------
// Metadata
// --------------------------------------------------

func (m *AudioModule) Category() modules.ModuleCategory {
	return modules.ModuleDomain
}

func (m *AudioModule) SupportedPlatforms() []internal_environment.PlatformClass {
	return nil // capability-driven only
}

func (m *AudioModule) Optional() bool {
	return true // do not block system bootstrap
}

func (m *AudioModule) DependsOn() []string {
	return []string{"StorageModule"}
}

// --------------------------------------------------
// Repository Placeholder
// --------------------------------------------------

type AudioRepository struct{}

func NewAudioRepository() *AudioRepository {
	return &AudioRepository{}
}
