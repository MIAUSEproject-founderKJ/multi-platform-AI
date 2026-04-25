//internal/network/network_discovery_service.go

package network

import (
	"context"
	"os"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/hashicorp/mdns"
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

// ListenForPeers searches for other StrataCore nodes on the network.
func ListenForPeers(ctx context.Context, peerFound func(entry *mdns.ServiceEntry)) {
	entriesCh := make(chan *mdns.ServiceEntry, 10)

	go func() {
		for {
			select {
			case entry := <-entriesCh:
				logging.Info("[MESH] Found potential peer: %s at %s:%d", entry.Name, entry.AddrV4, entry.Port)
				peerFound(entry)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Query for our specific AIOS service type
	params := mdns.DefaultParams("_strata-aios._tcp")
	params.Entries = entriesCh
	params.WantUnicastResponse = true

	mdns.Query(params)
}
