//core/security/env_key.go

package security

import (
	"fmt"
)

func lastKnownEnvKey(machineID string) string {
	return fmt.Sprintf("env/%s/last-known", machineID)
}
