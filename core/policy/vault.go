//core/security/vault.go

package security

import (
	"os"
	"path/filepath"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/apppath"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type IsolatedVault struct {
	BaseDir string
}

func NewVault() *IsolatedVault {
	// Use your apppath package to determine where the vault lives
	dataDir := apppath.GetDataDir()
	vaultPath := filepath.Join(dataDir, "vault")

	// Ensure the directory exists
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		os.MkdirAll(vaultPath, 0700) // Restricted permissions
	}

	return &IsolatedVault{BaseDir: vaultPath}
}

// IsMissingMarker checks if this is the "First Boot"
func (v *IsolatedVault) IsMissingMarker(name string) bool {
	markerPath := filepath.Join(v.BaseDir, name)
	_, err := os.Stat(markerPath)
	return os.IsNotExist(err)
}

// WriteMarker seals the "First Boot" phase
func (v *IsolatedVault) WriteMarker(name string) error {
	markerPath := filepath.Join(v.BaseDir, name)
	logging.Info("Sealing vault marker: %s", name)
	return os.WriteFile(markerPath, []byte("PROVISIONED"), 0600)
}
