//internal/schema/capabilities.go

package schema

import "time"

type CapabilitySet uint64
type Capability = CapabilitySet

const (
	CapDisplay Capability = 1 << iota
	CapKeyboard
	CapTouch
	CapMicrophone
	CapSpeaker
	CapCamera
	CapGPU
	CapSecureEnclave
	CapNetwork
	CapCANBus
	CapBiometric
	CapHighFreqSensor
	CapFileSystem
	CapSafetyCritical
	CapPersistentCloudLink
	CapIndustrialIO
	CapLocalStorage
)

func (c *CapabilitySet) Add(cap Capability) {
	*c |= cap
}

func (c *CapabilitySet) Remove(cap Capability) {
	*c &= ^cap
}

func (c CapabilitySet) Has(cap Capability) bool {
	return c&cap != 0
}

func (c CapabilitySet) HasAll(required CapabilitySet) bool {
	return c&required == required
}

func (c CapabilitySet) HasAny(mask CapabilitySet) bool {
	return c&mask != 0
}

func (c CapabilitySet) IsZero() bool {
	return c == 0
}

type CapabilityStatus int

const (
	CapOK CapabilityStatus = iota
	CapDegraded
	CapUnavailable
)

type CapabilityInfo struct {
	Available bool
	Status    CapabilityStatus
	LastCheck int64
}

type CapabilityProfile struct {
	Set   CapabilitySet
	Stats map[Capability]CapabilityInfo
}

func (cp *CapabilityProfile) RecomputeSet() {
	var set CapabilitySet
	for cap, info := range cp.Stats {
		if info.Status == CapOK {
			set |= cap
		}
	}
	cp.Set = set
}

func NewCapabilityProfile() *CapabilityProfile {
	return &CapabilityProfile{
		Set:   0,
		Stats: make(map[Capability]CapabilityInfo),
	}
}

func (cp *CapabilityProfile) Mark(cap Capability, status CapabilityStatus) {
	info := CapabilityInfo{
		Available: status == CapOK,
		Status:    status,
		LastCheck: time.Now().Unix(),
	}

	cp.Stats[cap] = info

	if status == CapOK {
		cp.Set.Add(cap)
	} else {
		cp.Set.Remove(cap)
	}
}

func (cp *CapabilityProfile) IsHealthy(cap Capability) bool {
	info, ok := cp.Stats[cap]
	return ok && info.Status == CapOK
}
