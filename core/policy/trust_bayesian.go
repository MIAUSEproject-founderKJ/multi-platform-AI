//MIAUSEproject-founderKJ/multi-platform-AI/core/policy/trust_bayesian.go

package policy

import (
	"multi-platform-AI/configs/defaults"
	"multi-platform-AI/internal/mathutil"
	"multi-platform-AI/internal/logging"
)

type TrustEvaluator struct {
	MinThreshold mathutil.Q16EnvConfd
}

// EvaluateCalculate takes the current EnvConfig and determines the system's "Trust State"
func (te *TrustEvaluator) Evaluate(env *defaults.EnvConfig) mathutil.Q16EnvConfd {
	logging.Info("[POLICY] Evaluating Bayesian Trust Matrix...")

	// 1. Prior Probability: Start with Attestation
	// If the code is tampered with, trust is immediately zero.
	if !env.Attestation.Valid {
		return 0
	}

	// 2. Evidence Integration: Combine Bus Confidence and Identity Confidence
	var combinedEvidence float64 = 1.0

	// Deduct trust based on Hardware Bus confidence (Q16)
	for _, bus := range env.Hardware.Buses {
		// We use a weighted average: lower confidence in critical buses drops the score faster
		combinedEvidence *= bus.Confidence.Float64()
	}

	// 3. Platform Specific Penalties
	// If the platform is "Portable" (USB/Stranger), we apply a logical penalty
	if env.Platform.Source == "probabilistic_match" {
		combinedEvidence *= 0.8
	}

	finalTrust := mathutil.Q16FromFloat(combinedEvidence)
	logging.Info("[POLICY] Bayesian Trust Score: %.2f%%", finalTrust.Float64()*100)

	return finalTrust
}