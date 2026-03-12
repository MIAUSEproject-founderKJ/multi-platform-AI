//internal/schema/session.go

package schema

import "time"

type UserSession struct {
	SessionID   string
	Platform    PlatformClass
	Entity      EntityType
	Tier        TierType
	Service     ServiceType
	Permissions map[string]bool
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Permissions core.PermissionSet
}
