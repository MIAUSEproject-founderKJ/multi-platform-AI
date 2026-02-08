// core/kernel_swarm.go
package core

import (
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/commands"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/api/network"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/bridge"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/logging"
)

func (k *Kernel) OnSwarmCommandReceived(cmd network.SwarmCommand) {
	// 1. REPLAY PROTECTION: Check if we've seen this Nonce/Timestamp before
	if time.Since(cmd.Header.Timestamp) > 5*time.Second {
		return // Ignore stale commands
	}

	// 2. SIGNATURE VERIFICATION: Is the Issuer a trusted Commander?
	if !k.Vault.VerifySignature(cmd.Header.IssuerID, cmd.Payload, cmd.Header.Signature) {
		logging.Error("[SECURITY] Rejected UNSIGNED or FORGED swarm command!")
		return
	}

	// 3. EXECUTION: Map swarm intent to local Bridge actions
	logging.Warn("[SWARM] Executing Signed Command: %s", cmd.Payload.Action)

	switch cmd.Payload.Action {
	case "SWARM_RECALL":
		k.ProcessCommand(commands.Task{Type: commands.CmdNavigate, Params: map[string]interface{}{"point": "HOME"}})
	case "EMERGENCY_STOP":
		k.Bridge.TransitionTo(bridge.StateEmergencyOff)
	}
}
