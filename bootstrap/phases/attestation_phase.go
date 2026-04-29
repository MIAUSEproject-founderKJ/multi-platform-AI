//bootstrap/phases/attestation_phase.go

package bootstrap_phase

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

func PhaseAttestation(
	vault verification_persistence.VaultStore,
	identity *internal_environment.MachineIdentity,
	bootSeq *internal_environment.BootSequence,
	preSession *user_setting.UserSession,
) (*user_setting.UserSession, error) {

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
