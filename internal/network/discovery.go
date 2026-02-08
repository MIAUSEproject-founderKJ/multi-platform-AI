//internal/network/discovery.go

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type DiscoveryService struct {
	server *mdns.Server
}

// StartBroadcasting makes this node visible to other nodes in the swarm.
func StartBroadcasting(nodeID string, port int) (*DiscoveryService, error) {
	// 1. Setup service metadata (including the Platform Class for quick filtering)
	host, _ := os.Hostname()
	info := []string{"version=1.0", "node_id=" + nodeID}
	
	service, err := mdns.NewMDNSService(host, "_strata-aios._tcp", "", "", port, nil, info)
	if err != nil {
		return nil, err
	}

	// 2. Start the mDNS server
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}

	logging.Info("[MESH] Discovery Broadcast started for node: %s", nodeID)
	return &DiscoveryService{server: server}, nil
}