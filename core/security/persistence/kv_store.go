//core/security/persistence/kv_store.go

package verification_persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func (v *IsolatedVault) Read(collection, key string, out interface{}) (bool, error) {
	filename := fmt.Sprintf("%s_%s.json", collection, key)
	data, err := v.readDecrypted(filename)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	if err := json.Unmarshal(data, out); err != nil {
		return false, fmt.Errorf("failed to unmarshal kv record (%s/%s): %w", collection, key, err)
	}

	return true, nil
}

func (v *IsolatedVault) Write(collection, key string, value interface{}) error {
	filename := fmt.Sprintf("%s_%s.json", collection, key)
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal kv record (%s/%s): %w", collection, key, err)
	}

	return v.writeEncrypted(filename, data)
}

func (v *IsolatedVault) Exists(collection, key string) (bool, error) {
	filename := fmt.Sprintf("%s_%s.json", collection, key)
	targetPath, err := v.resolvePath(filename)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
