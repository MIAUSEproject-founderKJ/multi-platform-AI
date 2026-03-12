// boot/identity/bind_hardware.go
package boot

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/platform/classify"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type IdentityManager struct {
	MachineName string
	OS          string
	Arch        string

	EntityType schema.EntityType
	TierType   schema.TierType

	Hardware    schema.HardwareProfile
	Environment *schema.EnvConfig
}

func (id *IdentityManager) BindHardware(env *schema.EnvConfig) {

	// Attach machine identity fields already known by the boot manager
	env.Identity.MachineName = id.MachineName
	env.Identity.OS = id.OS
	env.Identity.Arch = id.Arch

	// Ensure schema version is correct
	env.SchemaVersion = schema.CurrentVersion

	// Ensure generation timestamp exists
	if env.GeneratedAt.IsZero() {
		env.GeneratedAt = time.Now()
	}

	// If hardware profile is present but platform is not resolved,
	// run platform inference.
	if !env.Platform.Locked {
		classify.RunPlatformInference(env)
	}

	// Ensure attestation exists
	if env.Attestation.EnvHash == "" {
		classify.ComputeHardwareFingerprint(env)
	}

	// Attach environment to identity manager
	id.Environment = env
}