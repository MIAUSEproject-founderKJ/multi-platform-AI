//internal/keys/env_keys.go

package keys

import (
	"fmt"
)

func LastKnownEnvKey(machineID string) string {
	return fmt.Sprintf("env/%s/last-known", machineID)
}
