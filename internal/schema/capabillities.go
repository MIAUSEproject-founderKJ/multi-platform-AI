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

type CapabilityDescriptor struct {
	SupportsGoalControl        bool
	SupportsRegisterControl    bool
	SensorOnly                 bool
	HasSafetyEnvelope          bool
	SupportsAcceleratedCompute bool
}

func (c CapabilitySet) Has(flag CapabilitySet) bool {
	return c&flag != 0
}

func (c *CapabilitySet) Add(flag CapabilitySet) {
	*c |= flag
}
func (c *CapabilitySet) Remove(flag CapabilitySet) {
	*c &^= flag
}
