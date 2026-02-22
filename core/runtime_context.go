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



type BootProfile struct {
	Type string // FirstBoot | FastBoot | RecoveryBoot
}


type RuntimeContext struct {
    PlatformClass schema.PlatformClass
    Capabilities  CapabilitySet
    Service       ServiceType
    Entity        EntityType
    Tier          TierType
    BootMode      BootMode
    Permissions   PermissionSet
}

