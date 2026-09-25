//internal/schema/environment/attestation.go

package internal_environment

import internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"

type EnvAttestation struct {
	Locked       bool                      `json:"locked"`
	Valid        bool                      `json:"valid"`
	Level        internal_common.BootTrust `json:"level"`
	EnvHash      string                    `json:"env_hash"`
	SessionToken string                    `json:"session_token,omitempty"`
}

type SchemaInfo struct {
	Version int
	Name    string
	Created string
}

const (
	schemaName    = "environment-schema"
	schemaCreated = "2026-03-13"
)

var Current = SchemaInfo{Version: 2, Name: schemaName, Created: schemaCreated}

const CurrentVersion = 2 // keep exactly one number

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
