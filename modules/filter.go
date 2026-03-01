// modules/filter.go
package modules

import (
	"log/slog"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/runtime"
)

func FilterModules(registry []DomainModule, ctx *runtime.ExecutionContext) []DomainModule {

	var filtered []DomainModule

	for _, m := range registry {

		// 1. Platform compatibility
		if !platformSupported(m, ctx.PlatformClass) {
			continue
		}

		// 2. Hardware / capability compatibility
		if !capabilitiesSatisfied(m, ctx.Capabilities) {
			if m.Optional() {
				continue
			}
			// hard requirement missing
			panic("required capability missing for module: " + m.Name())
		}

		// 3. Optimizer-based gating
		if !optimizerAllows(m, ctx) {
			continue
		}

		filtered = append(filtered, m)
	}

	return filtered
}

//platform filtering
func platformSupported(m DomainModule, p runtime.PlatformClass) bool {
	for _, sp := range m.SupportedPlatforms() {
		if sp == p {
			return true
		}
	}
	return false
}

//capability filtering
/*This prevents:
• DatabaseSink without persistent storage
• GPU inference module without GPU
• CAN adapter without can_bus*/
func capabilitiesSatisfied(m DomainModule, caps map[string]bool) bool {

	for _, req := range m.RequiredCapabilities() {
		if !caps[req] {
			return false
		}
	}
	return true
}

func optimizerAllows(m DomainModule, ctx *runtime.ExecutionContext) bool {

	switch ctx.Optimizer.PrecisionMode() {

	case optimization.PrecisionAggressive:
		// disable non-critical analytics
		if m.Name() == "telemetry" {
			return false
		}

	case optimization.PrecisionReduced:
		// allow most modules
		return true

	case optimization.PrecisionFull:
		return true
	}

	return true
}

