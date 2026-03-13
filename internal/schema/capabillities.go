//internal/schema/service.go

package schema

type CapabilitySet uint64

const (
	CapCANBus CapabilitySet = 1 << iota
	CapSecureEnclave
	CapIndustrialIO
	CapNetwork
	CapLocalStorage
	CapBiometric
	CapPersistentCloudLink
	CapSafetyCritical
)
