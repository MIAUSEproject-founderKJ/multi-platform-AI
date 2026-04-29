//core/security/persistence/kv_store.go

package verification_persistence

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func (v *IsolatedVault) Read(collection, key string, out interface{}) (bool, error) {
	path := filepath.Join(v.BaseDir, collection+"_"+key+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if err := json.Unmarshal(data, out); err != nil {
		return false, err
	}

	return true, nil
}

func (v *IsolatedVault) Write(collection, key string, value interface{}) error {
	path := filepath.Join(v.BaseDir, collection+"_"+key+".json")

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (v *IsolatedVault) Exists(collection, key string) (bool, error) {
	path := filepath.Join(v.BaseDir, collection+"_"+key+".json")

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
