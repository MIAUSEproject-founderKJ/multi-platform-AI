// internal/schema/version.go
package schema

// CurrentVersion defines the active schema version used by the runtime.
const CurrentVersion = 1

type SchemaInfo struct {
	Version int
	Name    string
	Created string
}

var Current = SchemaInfo{
	Version: 1,
	Name:    "environment-schema",
	Created: "2026-03-13",
}

func Migrate(env *EnvConfig) *EnvConfig {

	switch env.SchemaVersion {

	case 1:
		return migrateV1toV2(env)

	}

	return env
}
func migrateV1toV2(env *EnvConfig) *EnvConfig {
	env.SchemaVersion = 2
	return env
}
