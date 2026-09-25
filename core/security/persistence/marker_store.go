//core/security/persistence/marker_store.go

package verification_persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	internal_boot "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/boot"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
	"go.uber.org/zap"
)

const firstBootVaultKey = "machine_first_boot_marker"

func (v *IsolatedVault) IsMissingMarker(name string) bool {
	targetPath, err := v.resolvePath(name)
	if err != nil {
		return true
	}
	_, err = os.Stat(targetPath)
	return os.IsNotExist(err)
}

func (v *IsolatedVault) WriteMarker(name string) error {
	logging.Info("Sealing state marker", zap.String("marker", name))
	return v.writeEncrypted(name, []byte("PROVISIONED"))
}

func (v *IsolatedVault) LoadFirstBootMarker() (*internal_boot.FirstBootMarker, error) {
	raw, err := v.readDecrypted(firstBootVaultKey + ".json")
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load first bootstrap marker: %w", err)
	}

	var marker internal_boot.FirstBootMarker
	if err := json.Unmarshal(raw, &marker); err != nil {
		return nil, fmt.Errorf("failed to unmarshal first bootstrap marker: %w", err)
	}

	return &marker, nil
}

func (v *IsolatedVault) MarkFirstBoot(marker *internal_boot.FirstBootMarker) error {
	data, err := json.Marshal(marker)
	if err != nil {
		return fmt.Errorf("failed to marshal first boot marker: %w", err)
	}

	return v.writeEncrypted(firstBootVaultKey+".json", data)
}
