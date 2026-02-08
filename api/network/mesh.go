//api/network/mesh.go

package network

import (
	"time"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/policy"
)

type NodePulse struct {
	SourceID    string                 `json:"source_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Identity    string                 `json:"platform_class"`
	Trust       policy.TrustDescriptor `json:"trust_state"`
	Position    [3]float64             `json:"position"` // x, y, z or lat, lng, alt
}

type MeshNetwork struct {
	Peers       map[string]NodePulse
	DiscoveryIP string
}