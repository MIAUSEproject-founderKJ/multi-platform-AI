// core/security/persistence/config_store.go

package verification_persistence

import (
	"encoding/json"
	"fmt"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

func (v *IsolatedVault) SaveConfig(name string, config *internal_environment.EnvConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal env config (%s): %w", name, err)
	}

	return v.writeEncrypted("config_"+name+".json", data)
}

func (v *IsolatedVault) LoadConfig(name string) (*internal_environment.EnvConfig, error) {
	data, err := v.readDecrypted("config_" + name + ".json")
	if err != nil {
		return nil, err
	}

	var cfg internal_environment.EnvConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal env config (%s): %w", name, err)
	}

	return &cfg, nil
}
