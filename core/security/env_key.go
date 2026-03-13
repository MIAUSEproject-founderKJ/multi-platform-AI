//core/security/env_key.go

package security

import (
	"fmt"
)

func LastKnownEnvKey(machineID string) string {
	return fmt.Sprintf("env/%s/last-known", machineID)
}
