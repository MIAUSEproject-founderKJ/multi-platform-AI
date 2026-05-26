//core/security/persistence/marker_store.go

package verification_persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
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

func (v *IsolatedVault) MarkFirstBoot(marker *internal_boot.FirstBootMarker) error {
	path := filepath.Join(v.BaseDir, firstBootVaultKey+".json")

	data, err := json.MarshalIndent(marker, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
