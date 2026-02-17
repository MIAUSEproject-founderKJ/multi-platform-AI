//core/runtime_context.go
//RuntimeContext must be constructed once and treated as immutable.


package core

type Capability string

const (
	CapCANBus              Capability = "CAN_BUS"
	CapBiometric           Capability = "BIOMETRIC"
	CapHighFreqSensor      Capability = "HIGH_FREQ_SENSOR"
	CapFileSystem          Capability = "FILE_SYSTEM"
	CapMicrophone          Capability = "MICROPHONE"
	CapSafetyCritical      Capability = "SAFETY_CRITICAL"
	CapPersistentCloudLink Capability = "PERSISTENT_CLOUD"
)

type RuntimeContext struct {
	Platform   PlatformProfile
	Service    ServiceProfile
	Identity   IdentityProfile
	Tier       TierProfile
	Boot       BootProfile
	Policy     PolicyProfile
}

type PlatformProfile struct {
	Name         string
	Capabilities map[Capability]bool
}

type ServiceProfile struct {
	Name string
}

type IdentityProfile struct {
	Entity string // Personal | Organization | Stranger | Tester
}

type TierProfile struct {
	Name string // Funder | Non-Funder
}

type BootProfile struct {
	Type string // FirstBoot | FastBoot | RecoveryBoot
}

type PolicyProfile struct {
	Permissions map[string]bool
}
