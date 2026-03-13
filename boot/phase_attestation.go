//boot/phase_attestation.go

package boot

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

func PhaseAttestation(
	v *security.IsolatedVault,
	identity *schema.MachineIdentity,
	bs *schema.BootSequence,
) (*schema.UserSession, error) {

	am := &auth.AuthManager{
		Vault:    v,
		Identity: identity,
		Platform: identity.PlatformType,
		Entity:   bs.Env.EntityType,
		Tier:     bs.Env.TierType,
	}

	return am.LoginOrSignUp()
}
