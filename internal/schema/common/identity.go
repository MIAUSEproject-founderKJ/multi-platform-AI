// internal/schema/common/identity.go
package internal_common

type EntityKind uint8

const (
	EntityPersonal EntityKind = iota
	EntityOrganization
	EntityStranger
	EntityTester
)

type PlatformClass string

const (
	PlatformComputer   PlatformClass = "computer"
	PlatformMobile     PlatformClass = "mobile"
	PlatformEmbedded   PlatformClass = "embedded"
	PlatformIndustrial PlatformClass = "industrial"
	PlatformVehicle    PlatformClass = "vehicle"
	PlatformRobot      PlatformClass = "robot"
	PlatformUnknown    PlatformClass = "unknown"
)
