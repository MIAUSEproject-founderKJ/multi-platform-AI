//api/network/broadcaster.go

package network

import (
	"encoding/json"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

type SwarmBroadcaster struct {
	engine *NetworkEngine
}

// Dispatch sends a command to the entire trusted mesh.
func (sb *SwarmBroadcaster) Dispatch(action string, params map[string]interface{}) error {
	cmd := SwarmCommand{}
	cmd.Payload.Action = action
	cmd.Payload.Params = params
	cmd.Payload.Target = "ALL"
	cmd.Header.Timestamp = time.Now()

	// 1. Sign the command using the Local Vault's Private Key
	// signature, _ := sb.engine.Kernel.Vault.Sign(cmd.Payload)
	// cmd.Header.Signature = signature

	data, _ := json.Marshal(cmd)

	logging.Info("[SWARM] Dispatching Global Command: %s", action)

	// 2. Broadcast over UDP/TCP to all discovered peers in the SwarmMap
	for id, peer := range sb.engine.Kernel.Swarm.Peers {
		go sb.engine.sendToPeer(peer.Addr, data)
		logging.Debug("[SWARM] -> Propagating to peer %s", id)
	}

	return nil
}
