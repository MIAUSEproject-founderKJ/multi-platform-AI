//boot/phase_attestation.go

package boot

import (
	"fmt"

)

func PhaseAttestation(
    v *security.IsolatedVault,
    identity *schema.MachineIdentity,
    bs *schema.BootSequence,
) (*schema.UserSession, error) {

    am := &auth.AuthManager{
        Vault:    v,
        Identity: identity,
        Platform: identity.Platform,
        Entity:   bs.Env.Identity.EntityType,
        Tier:     bs.Env.Identity.TierType,
    }

    return am.LoginOrSignUp()
}