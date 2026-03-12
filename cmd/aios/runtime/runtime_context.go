//cmd/aios/runtime/runtime_context.go
//RuntimeContext must be constructed once and treated as immutable.

package runtime

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/optimization"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/router"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

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
	Capabilities  core.CapabilitySet
	Service       core.ServiceType
	Entity        core.EntityType
	Tier          core.TierType
	BootMode      core.BootMode
	Permissions   core.PermissionSet
	Router        *router.Router
	Optimizer     optimization.Optimizer
}
