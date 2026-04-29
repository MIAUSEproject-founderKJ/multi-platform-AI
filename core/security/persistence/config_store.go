// core/security/persistence/config_store.go
package verification_persistence

import (
	"encoding/json"
	"os"
	"path/filepath"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

func (v *IsolatedVault) SaveConfig(name string, config *internal_environment.EnvConfig) error {
	path := filepath.Join(v.BaseDir, name+".json")

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (v *IsolatedVault) LoadConfig(name string) (*internal_environment.EnvConfig, error) {
	path := filepath.Join(v.BaseDir, name+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg internal_environment.EnvConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
