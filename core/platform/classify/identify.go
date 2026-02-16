// core/platform/classify/identify.go
package classify

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// SecurityProfile defines the "Gate Height" for the platform
type SecurityProfile struct {
	PlatformClass        string // Automotive | Workstation | Industrial | Embedded
	RequiresBiometry     bool
	RequiresKeyTelemetry bool
	AuthTimeoutMinutes   int
}

// Identify maps hardware identity to a specific operational profile.
func Identify(id *probe.HardwareIdentity) (*SecurityProfile, error) {
	logging.Info("[CLASSIFY] Matching Hardware Identity to Security Profile...")

	profile := &SecurityProfile{}

	switch id.PlatformType {

	// ------------------------------------------
	// AUTOMOTIVE / VEHICLE
	// ------------------------------------------
	case schema.PlatformVehicle, "Automotive":
		profile.PlatformClass = "Automotive"
		profile.RequiresBiometry = true
		profile.RequiresKeyTelemetry = true
		profile.AuthTimeoutMinutes = 0 // always re-verify on ignition

	// ------------------------------------------
	// INDUSTRIAL / EMBEDDED / ROBOT / DRONE
	// ------------------------------------------
	case schema.PlatformIndustrial, schema.PlatformEmbedded, schema.PlatformRobot, schema.PlatformDrone, schema.PlatformGamePad:
		profile.PlatformClass = "Industrial" // safe default for embedded
		profile.RequiresBiometry = false
		profile.RequiresKeyTelemetry = false
		profile.AuthTimeoutMinutes = 1440 // 24h persistent trust

	// ------------------------------------------
	// HIGH-LEVEL COMPUTE DEVICES (PC / LAPTOP / MOBILE / TABLET)
	// ------------------------------------------
	case schema.PlatformComputer, schema.PlatformLaptop, schema.PlatformTablet, schema.PlatformMobile:
		profile.PlatformClass = "Workstation"
		profile.RequiresBiometry = true
		profile.RequiresKeyTelemetry = false
		profile.AuthTimeoutMinutes = 480 // standard work shift

	// ------------------------------------------
	// UNKNOWN / FALLBACK
	// ------------------------------------------
	default:
		profile.PlatformClass = "Embedded" // minimal safe profile
		profile.RequiresBiometry = false
		profile.RequiresKeyTelemetry = false
		profile.AuthTimeoutMinutes = 1440
		logging.Warn("[CLASSIFY] Unknown platform type %s, applying default embedded profile", id.PlatformType)
	}

	logging.Info("[CLASSIFY] Profile Locked: %s (Biometry: %v, KeyTelemetry: %v)",
		profile.PlatformClass, profile.RequiresBiometry, profile.RequiresKeyTelemetry)

	return profile, nil
}
