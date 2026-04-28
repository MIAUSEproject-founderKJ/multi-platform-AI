//core/verification/persistence/vault_store.go

package verification_persistence

import (
	"fmt"
	"os"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/apppath"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/bootstrap"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

type IsolatedVault struct {
	BaseDir string
	Key     []byte
}

func OpenVault() (*IsolatedVault, error) {
	path := apppath.GetVaultPath()

	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, fmt.Errorf("vault init failed: %w", err)
	}

	logging.Info("[VAULT] initialized at %s", path)

	return &IsolatedVault{
		BaseDir: path,
		Key:     LoadSecureKey(),
	}, nil
}

type VaultStore interface {
	LoadConfig(key string) (*internal_environment.EnvConfig, error)
	SaveConfig(key string, cfg *internal_environment.EnvConfig) error
	MarkFirstBoot(machineID string) error
	LoadGoldenHash(machine string) (string, error)
	SealGoldenHash(machine string, hash []byte) error
	LoadFirstBootMarker() (*internal_boot.FirstBootMarker, error)
	SaveFirstBootMarker(*internal_boot.FirstBootMarker) error
	Read(key, id string, out interface{}) (bool, error)
	Write(key, id string, value interface{}) error
	Exists(collection string, key string) (bool, error)
}
