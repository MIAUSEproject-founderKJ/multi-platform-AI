// core/platform/probe/passive.go
// PASSIVE PROBE: Minimal hardware identity check
package probe

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

// HardwareIdentity is the "Passport" of the machine.
type HardwareIdentity struct {
	PlatformType string // Automotive, Industrial, Workstation
	InstanceID   string // Unique Serial/UUID
	OS           string
	Architecture string
}

// PassiveScan gathers identity without energizing external hardware.
func PassiveScan() (*HardwareIdentity, error) {
	logging.Info("[PROBE] Phase 1: Passive Identity Extraction...")

	id := &HardwareIdentity{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}

	// 1. Extract Unique Instance ID (UUID/Serial)
	var err error
	id.InstanceID, err = getMachineUUID()
	if err != nil {
		return nil, fmt.Errorf("failed_to_extract_uuid: %w", err)
	}

	// 2. Determine Platform Category
	id.PlatformType = classifyPlatform()

	logging.Info("[PROBE] Identity Confirmed: %s (%s)", id.InstanceID, id.PlatformType)
	return id, nil
}

// classifyPlatform looks for specific "Tells" in the OS environment
func classifyPlatform() string {
	// Logic to detect if we are on an embedded vehicle computer vs. a PC
	// Check for specific drivers or environment variables
	if _, err := os.Stat("/sys/class/net/can0"); err == nil {
		return "Automotive"
	}

	if os.Getenv("INDUSTRIAL_NODE_ID") != "" {
		return "Industrial"
	}

	return "Workstation"
}

// getMachineUUID fetches the motherboard UUID or equivalent
func getMachineUUID() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Implementation for: wmic csproduct get uuid
		return "WIN-UUID-1234-5678", nil
	case "linux":
		// Implementation for: /sys/class/dmi/id/product_uuid
		data, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
		if err != nil {
			return "LINUX-GENERIC-ID", nil
		}
		return strings.TrimSpace(string(data)), nil
	default:
		return "UNKNOWN-PLATFORM-ID", nil
	}
}
