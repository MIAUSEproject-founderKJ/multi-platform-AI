// boot/probe/heartbeat.go
// PASSIVE PROBE: Minimal hardware identity check
package probe

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/platform/classify"
)

// PassiveScan performs a minimal, non-invasive identity extraction and resolves platform.
func PassiveScan(env *schema.EnvConfig) (*schema.MachineIdentity, error) {
    logging.Info("[PROBE] Phase 1: Passive Identity Extraction...")

    // 1. Detect rough platform first
    roughPlatform := classify.DetectPlatformClass(&env.Hardware)
    env.Platform.Final = roughPlatform
    logging.Info("[PROBE] Rough Platform Detected: %s", roughPlatform)

    // 2. Build full Machine Identity
    id := getMachineUUID()
    env.Identity = id

	if env.Platform.Final == "" {
    env.Platform.Final = classify.DetectPlatformClass(&env.Hardware)
    logging.Warn("[PROBE] Using fallback platform detection: %s", env.Platform.Final)
}

    // 3. Run full heuristic scoring to finalize platform
    classify.RunPlatformInference(env)

    logging.Info("[PROBE] Identity Confirmed: %s (%s)", id.MachineID, env.Platform.Final)
    return &id, nil
}

// getMachineUUID generates a unique machine ID from hostname, OS, and architecture
func getMachineUUID() schema.MachineIdentity {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "unknown-node"
	}

	raw := fmt.Sprintf("%s-%s-%s", hostname, runtime.GOOS, runtime.GOARCH)
	hash := sha256.Sum256([]byte(raw))

	return schema.MachineIdentity{
		MachineID:   hex.EncodeToString(hash[:]),
		MachineName: hostname,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
	}
}