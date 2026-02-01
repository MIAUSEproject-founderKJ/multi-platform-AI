//MIAUSEproject-founderKJ/multi-platform-AI/core/platform/types.go
package platform

import "multi-platform-AI/configs/platforms"

// BootSequence represents the finalized state of the system after 
// the Nucleus has completed its initialization.
type BootSequence struct {
	PlatformID platforms.PlatformClass `json:"platform_id"` // e.g., Vehicle, Laptop
	TrustScore float64                `json:"trust_score"` // 0.0 to 1.0 (Bayesian)
	IsVerified bool                   `json:"is_verified"` // Attestation result
	Mode       string                 `json:"mode"`        // Autonomous | Discovery | Safe
	UserRole   string                 `json:"user_role"`   // Operator | Admin
}