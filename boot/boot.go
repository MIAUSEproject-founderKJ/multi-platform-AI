//boot/boot.go

package boot

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/platform/probe"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// RunBootSequence performs full boot → verification → session creation
func RunBootSequence(v *security.IsolatedVault) (*schema.BootSequence, *schema.UserSession, error) {

	// 1. Passive identity scan (non-invasive)
	raw, err := probe.PassiveScan()
	if err != nil {
		return nil, nil, fmt.Errorf("passive scan failed: %w", err)
	}

	// 2. Normalize into canonical identity model
	identity := &schema.MachineIdentity{
		MachineName: raw.InstanceID,
		Platform:    raw.PlatformType,
		OS:          raw.OS,
		Arch:        raw.Architecture,
	}

	// 3. Instantiate BootManager AFTER identity binding
	bm := &BootManager{
		Vault:    v,
		Identity: identity,
	}

	// 4. Decide boot path (cold / fast boot)
	bs, err := bm.DecideBootPath()
	if err != nil {
		return nil, nil, fmt.Errorf("boot sequence failed: %w", err)
	}

	// 5. Formalize FirstBootMarker if cold boot
	if bs.Mode == schema.BootCold {
		if err := v.MarkFirstBoot(identity.MachineName); err != nil {
			return nil, nil, fmt.Errorf("failed to set FirstBootMarker: %w", err)
		}
	}

	// 6. AuthManager: login / sign-up based on platform
	am := &auth.AuthManager{
		Vault:    v,
		Identity: identity,
		Platform: raw.PlatformType,
		Entity:   bs.Env.Identity.EntityType, // provisional, can update based on registration
		Tier:     bs.Env.Identity.TierType,
	}

	session, err := am.LoginOrSignUp()
	if err != nil {
		return nil, nil, fmt.Errorf("authentication failed: %w", err)
	}

	// 7. Attach verified session to boot sequence for runtime context
	bs.Env.Attestation.SessionToken = session.SessionID

	return bs, session, nil
}
