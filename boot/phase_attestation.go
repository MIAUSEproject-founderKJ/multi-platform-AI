//boot/phase_attestation.go

package boot

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func PhaseAttestation(
	vault security.VaultStore,
	identity *schema.MachineIdentity,
	bootSeq *schema.BootSequence,
	preSession *schema.UserSession,
) (*schema.UserSession, error) {
	// Load credentials from Vault
	var cred struct {
		UserID   string
		Password string
	}

	found, err := v.Read("credentials", identity.MachineID, &cred)
	if err != nil || !found {
		return nil, err
	}

	am := &auth.AuthManager{
		Vault:    v,
		Identity: identity,
		Platform: identity.PlatformType,
		Entity:   bs.Env.EntityType,
		Tier:     bs.Env.TierType,
	}

	return am.LoginOrSignUpInteractive()
}
