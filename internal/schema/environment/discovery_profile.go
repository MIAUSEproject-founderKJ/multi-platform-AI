//internal/schema/system/discovery_profile.go

package internal_environment

import "time"

type DiscoveryProfile struct {
	Physical          PhysicalProfile      `json:"physical"`
	Signal            SignalProfile        `json:"signal"`
	Nodes             []NodeDescriptor     `json:"nodes"`
	Protocol          ProtocolProfile      `json:"protocol"`
	Capabilities      CapabilityDescriptor `json:"capabilities"`
	DiscoveryDuration time.Duration        `json:"discovery_duration"`
}
