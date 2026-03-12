// core/types.go
package core

type CapabilitySet uint64
type PermissionSet uint64
type ServiceType uint8
type EntityType uint8
type TierType uint8

const (
	ServiceUnknown ServiceType = iota
	ServicePersonal
	ServiceEnterprise
	ServiceSystem
)

const (
	EntityUnknown EntityType = iota
	EntityUser
	EntityAdmin
	EntityDevice
)

const (
	TierFree TierType = iota
	TierPro
	TierEnterprise
)
