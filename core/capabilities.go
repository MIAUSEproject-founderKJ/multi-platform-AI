//core/capabilities.go

package core

const (
	CapCANBus CapabilitySet = 1 << iota
	CapBiometric
	CapSecureEnclave
	CapIndustrialIO
	CapNetwork
	CapLocalStorage
)

const (
	PermUser PermissionSet = 1 << iota
	PermAdmin
	PermDeviceControl
	PermFleetControl
)