// core/security/persistence/config_store.go
package security_persistence

import (
	"encoding/json"
	"os"
	"path/filepath"

	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

func (v *IsolatedVault) SaveConfig(name string, config *schema_system.EnvConfig) error {
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

func (v *IsolatedVault) LoadConfig(name string) (*schema_system.EnvConfig, error) {
	path := filepath.Join(v.BaseDir, name+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg schema_system.EnvConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
