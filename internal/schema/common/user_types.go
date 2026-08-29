//internal/schema/common/user_types.go

package internal_common

// ------------------------------------------------------------
// Tier System
// ------------------------------------------------------------
// Use TierType (string) externally for readability and compatibility. Use EntityType (uint8) internally for speed and clarity.
type TierType string

const (
	TierUnknown    TierType = "unknown"
	TierPersonal   TierType = "personal"
	TierEnterprise TierType = "enterprise"
	TierTester     TierType = "tester"
)
