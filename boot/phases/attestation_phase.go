//boot/phases/attestation_phase.go

package boot_phase

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	security_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	schema_identity "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/identity"
	schema_system "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
)

func PhaseAttestation(
	vault security_persistence.VaultStore,
	identity *schema_system.MachineIdentity,
	bootSeq *schema_system.BootSequence,
	preSession *schema_identity.UserSession,
) (*schema_identity.UserSession, error) {

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
