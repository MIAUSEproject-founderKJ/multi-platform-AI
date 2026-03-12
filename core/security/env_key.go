//core/security/env_key.go

package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

)


func lastKnownEnvKey(machineID string) string {
	return fmt.Sprintf("env/%s/last-known", machineID)
}
