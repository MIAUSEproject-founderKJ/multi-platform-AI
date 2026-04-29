//core/security/persistence/marker_store.go

package verification_persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap"
)

const firstBootVaultKey = "machine_first_boot_marker"

func (v *IsolatedVault) IsMissingMarker(name string) bool {
	_, err := os.Stat(filepath.Join(v.BaseDir, name))
	return os.IsNotExist(err)
}

func (v *IsolatedVault) WriteMarker(name string) error {
	path := filepath.Join(v.BaseDir, name)
	logging.Info("[VAULT] Sealing state marker: %s", name)
	return os.WriteFile(path, []byte("PROVISIONED"), 0600)
}

func (v *IsolatedVault) MarkFirstBoot(machineID string) error {
	marker := &internal_boot.FirstBootMarker{
		MachineID:   machineID,
		Initialized: true,
		CreatedAt:   time.Now(),
	}

	return v.SaveFirstBootMarker(marker)
}

func (v *IsolatedVault) LoadFirstBootMarker() (*internal_boot.FirstBootMarker, error) {
	path := filepath.Join(v.BaseDir, firstBootVaultKey+".json")

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load first bootstrap marker: %w", err)
	}

	var marker internal_boot.FirstBootMarker
	if err := json.Unmarshal(raw, &marker); err != nil {
		return nil, err
	}

	return &marker, nil
}

func (v *IsolatedVault) SaveFirstBootMarker(marker *internal_boot.FirstBootMarker) error {
	path := filepath.Join(v.BaseDir, firstBootVaultKey+".json")

	data, err := json.MarshalIndent(marker, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
