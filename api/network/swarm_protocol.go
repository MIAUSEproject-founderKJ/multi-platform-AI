//api/network/swarm_protocol.go

package network

import (
	"time"
)

type SwarmCommand struct {
	Header struct {
		IssuerID  string    `json:"issuer_id"`
		Timestamp time.Time `json:"timestamp"`
		Signature []byte    `json:"signature"` // HMAC or RSA signature
		Nonce     uint64    `json:"nonce"`     // Prevents replay attacks
	} `json:"header"`

	Payload struct {
		Action string                 `json:"action"` // e.g., "SWARM_RECALL", "GRID_SEARCH"
		Params map[string]interface{} `json:"params"`
		Target string                 `json:"target"` // "ALL" or specific NodeID
	} `json:"payload"`
}