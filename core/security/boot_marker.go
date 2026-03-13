// core/security/boot_marker.go
package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// --- Marker Logic ---

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

func (v *IsolatedVault) MarkFirstBoot(machineName string) error {

	marker := &schema.FirstBootMarker{
		MachineName: machineName,
		Initialized: true,
		CreatedAt:   time.Now(),
	}

	data, err := json.Marshal(marker)
	if err != nil {
		return fmt.Errorf("failed to serialize first boot marker: %w", err)
	}

	if err := v.Store(firstBootVaultKey, data); err != nil {
		return fmt.Errorf("failed to persist first boot marker: %w", err)
	}

	return nil
}

func (v *IsolatedVault) LoadFirstBootMarker() (*schema.FirstBootMarker, error) {

	raw, err := v.Load(firstBootVaultKey)
	if err != nil {
		return nil, err
	}

	var marker schema.FirstBootMarker
	if err := json.Unmarshal(raw, &marker); err != nil {
		return nil, err
	}

	return &marker, nil
}

func (v *IsolatedVault) Load(firstBootVaultKey string) (any, any) {
	panic("unimplemented")
}
