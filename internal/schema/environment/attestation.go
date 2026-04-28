//internal/schema/system/attestation.go

package internal_environment

// EnvAttestation defines the cryptographic seal of the environment
type EnvAttestation struct {
	Locked        bool          `json:"locked"`
	PlatformClass PlatformClass `json:"platform_class,omitempty"`
	Valid         bool          `json:"valid"`
	Level         BootTrust     `json:"level"` // "strong" | "weak" | "invalid"
	EnvHash       string        `json:"env_hash"`
	SessionToken  string        `json:"session_token,omitempty"`
}

type SchemaInfo struct {
	Version int
	Name    string
	Created string
}

type BootTrust uint8

const (
	TrustInvalid BootTrust = iota
	TrustWeak
	TrustStrong
)

// CurrentVersion defines the active schema version used by the runtime.
const CurrentVersion = 2

var Current = SchemaInfo{
	Version: 1,
	Name:    "environment-schema",
	Created: "2026-03-13",
}

func Migrate(env *EnvConfig) *EnvConfig {
	if env == nil {
		return nil
	}

	switch env.SchemaVersion {
	case 1:
		return migrateV1toV2(env)
	case 2:
		return env
	default:
		return env
	}
}

func migrateV1toV2(env *EnvConfig) *EnvConfig {
	env.SchemaVersion = 2
	return env
}
