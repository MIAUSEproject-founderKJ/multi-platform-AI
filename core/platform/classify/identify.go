//core/platform/classify/identify.go

package classify

import (
	"multi-platform-AI/core/platform/probe"
	"multi-platform-AI/internal/logging"
)

// SecurityProfile defines the "Gate Height" for the platform
type SecurityProfile struct {
	PlatformClass    string // Automotive | Workstation | Industrial
	RequiresBiometry bool
	RequiresKeyTelemetry bool
	AuthTimeoutMinutes int
}

// Identify maps hardware identity to a specific operational profile.
func Identify(id *probe.HardwareIdentity) (*SecurityProfile, error) {
	logging.Info("[CLASSIFY] Matching Hardware Identity to Security Profile...")

	profile := &SecurityProfile{}

	switch id.PlatformType {
	case "Automotive":
		// Scenario: FSD / Tractor / Harvester
		profile.PlatformClass = "Automotive"
		profile.RequiresBiometry = true
		profile.RequiresKeyTelemetry = true // Car-Key / Remote handshake
		profile.AuthTimeoutMinutes = 0      // Always re-verify on ignition

	case "Industrial":
		// Scenario: Smart-House / Factory
		profile.PlatformClass = "Industrial"
		profile.RequiresBiometry = false    // Uses NFC / Physical Interlock instead
		profile.RequiresKeyTelemetry = false
		profile.AuthTimeoutMinutes = 1440   // 24-hour persistent trust

	case "Workstation":
		// Scenario: Professional PC / Laptop
		profile.PlatformClass = "Workstation"
		profile.RequiresBiometry = true     // Windows Hello / TouchID
		profile.RequiresKeyTelemetry = false
		profile.AuthTimeoutMinutes = 480    // Standard work shift
	}

	logging.Info("[CLASSIFY] Profile Locked: %s (Biometry: %v)", 
		profile.PlatformClass, profile.RequiresBiometry)
		
	return profile, nil
}