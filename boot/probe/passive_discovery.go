// boot/probe/passive_discovery.go
// PASSIVE PROBE: Minimal hardware identity check
package probe

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot/platform/classify"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func PassiveDiscovery() (*schema.EnvConfig, error) {

	logging.Info("[PROBE] Phase 1: Passive Identity Extraction")

	id := buildMachineIdentity()

	env := &schema.EnvConfig{
		SchemaVersion: schema.CurrentVersion,
		GeneratedAt:   time.Now(),
		Identity:      id,
	}

	// run platform inference
	classify.RunPlatformInference(env)

	logging.Info(
		"[PROBE] Identity confirmed: %s (%s)",
		id.MachineID,
		env.Platform.Final,
	)

	return env, nil
}

func buildMachineIdentity() schema.MachineIdentity {

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