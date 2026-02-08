//core/kernel_network.go

package core

import (
	"fmt"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/network"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
	"github.com/hashicorp/mdns"
)

// HandleNewPeer initiates the "Handshake of Trust" with a discovered node.
func (k *Kernel) HandleNewPeer(entry *mdns.ServiceEntry) {
	logging.Info("[SWARM] Attempting Handshake with peer: %s (%s)", entry.Name, entry.AddrV4)

	// 1. Establish a temporary connection to verify identity
	// In a production build, this would be a TLS handshake using the Vault's Root CA.
	peerIdentity, err := k.requestIdentity(entry)
	if err != nil {
		logging.Warn("[SWARM] Handshake failed with %s: %v", entry.AddrV4, err)
		return
	}

	// 2. CRYPTOGRAPHIC CHECK: Does the remote EnvHash match our security policy?
	// We compare the remote node's binary hash against our "Known Good" database.
	if !k.Evaluator.VerifyRemoteHash(peerIdentity.EnvHash) {
		logging.Error("[SECURITY] QUARANTINED: Peer %s has an invalid or tampered binary!", entry.Name)
		return
	}

	// 3. ADMISSION: Add to the Active Swarm Map
	k.Swarm.Mu.Lock()
	k.Swarm.Peers[peerIdentity.NodeID] = network.PeerState{
		Addr:      entry.AddrV4.String(),
		LastPulse: peerIdentity.CurrentTrust,
		Status:    "TRUSTED_PEER",
	}
	k.Swarm.Mu.Unlock()

	logging.Info("[SWARM] Node %s admitted to the swarm. Collective intelligence increased.", entry.Name)
}