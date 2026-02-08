//internal/network/listener.go

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