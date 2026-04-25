// core/security/persistence/golden_hash_store.go
package security_persistence

// GoldenHashStore is a simple in-memory store for golden hashes used in security verification.

import (
	"os"
	"path/filepath"
)

func (v *IsolatedVault) SealGoldenHash(machine string, hash []byte) error {
	path := filepath.Join(v.BaseDir, "golden-"+machine)
	return os.WriteFile(path, hash, 0600)
}

func (v *IsolatedVault) LoadGoldenHash(machine string) (string, error) {
	path := filepath.Join(v.BaseDir, "golden-"+machine)

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
