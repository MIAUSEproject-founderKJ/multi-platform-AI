// bootstrap/probe/discovery.go

// ClassifyPlatform is the only function that sets env.Platform.Final.
// Cheap, bounded (~3s via CollectHardwareFingerprint's internal timeout),
// safe to run on every cold boot without a separate gate.

package probe

import (
	"context"
	"fmt"

	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/pkg/logging"
	"go.uber.org/zap"
)

func ClassifyPlatform(ctx context.Context, env *internal_environment.EnvConfig) error {
	fp, probeErrors := CollectHardwareFingerprint(ctx)
	if len(probeErrors) > 0 {
		logging.Warn("hardware probe errors", zap.Strings("errors", probeErrors))
	}
	env.Hardware = buildHardwareProfile(fp)
	runPlatformInference(env, fp)
	return nil
}

// EnrichHardwareProfile adds platform-specific deep detail (GPU/VRAM,
// bus nodes, signal properties) via subprocess calls. Requires
// env.Platform.Final to already be set by ClassifyPlatform.
func EnrichHardwareProfile(env *internal_environment.EnvConfig) error {
	if env.Platform.Final == "" {
		return fmt.Errorf("enrichment requires a classified platform; call ClassifyPlatform first")
	}
	switch env.Platform.Final {
	case internal_common.PlatformComputer, internal_common.PlatformMobile:
		populateCompute(env)
	case internal_common.PlatformVehicle, internal_common.PlatformRobot,
		internal_common.PlatformIndustrial, internal_common.PlatformEmbedded:
		populateEmbedded(env)
	default:
		phy, _ := discoverPhysical()
		env.Discovery.Physical = phy
		env.Discovery.Capabilities.SensorOnly = true
	}
	return nil
}
