//bootstrap\phases\runtime_phase.go

package bootstrap_phase

import (
	bootstrap_resolver "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/resolver"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/mutual_interaction"
)

func PhaseRuntime(session *user_setting.UserSession) error {
	cp, err := bootstrap_resolver.DeviceCapabilitiesResolver()
	if err != nil {
		return err
	}

	orch := BuildOrchestrator(cp)

	orch.StartAll(session)

	return nil
}

func ResolveAuthMode(
	caps *internal_environment.CapabilityProfile,
) user_setting.InteractionMode {
	if caps == nil {
		return user_setting.ModeCLI
	}

	return mutual_interaction.ResolveInteractionMode(
		nil,
		caps.Set,
	)

}

func (am *AuthManager) initializeRuntime(session *user_setting.UserSession) error {

	cp, err := bootstrap_resolver.DeviceCapabilitiesResolver()
	if err != nil {
		return err
	}

	orch := bootstrap_phase.BuildOrchestrator(cp)
	orch.StartAll(session)

	mode := mutual_interaction.ResolveInteractionMode(session.Config, cp.Set)

	bootstrap.Capabilities = cp.Set
	bootstrap.CapProfile = cp
	bootstrap.Mode = string(mode)
	bootstrap.Orchestrator = orch

	orch.Broadcast("Session initialized successfully")

	return nil
}
