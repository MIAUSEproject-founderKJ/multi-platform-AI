//boot\phases\attestation_phase.go

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

	var cred struct {
		UserID   string
		Password string
	}

	found, err := vault.Read("credentials", identity.MachineID, &cred)
	if err != nil || !found {
		return nil, err
	}

	am := &auth.AuthManager{
		Vault:    vault,
		Identity: identity,
		Platform: identity.PlatformType,
		Entity:   bootSeq.Env.EntityType,
		Tier:     bootSeq.Env.TierType,
	}

	// Prefer pre-authenticated session
	if preSession != nil {
		return preSession, nil
	}

	return am.Login(cred.UserID, cred.Password)
}
