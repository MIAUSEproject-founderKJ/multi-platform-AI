//MIAUSEproject-founderKJ/multi-platform-AI/core/policy/trust_bayesian.go

package policy

import (
	"fmt"
	"math"
	"multi-platform-AI/configs/configStruct"
)

// TrustDescriptor represents the final judgment of the system's integrity.
type TrustDescriptor struct {
	CurrentScore   float64       // 0.0 to 1.0 (Final Posterior Probability)
	Label          string        // NOMINAL, DEGRADED, CRITICAL
	OperationMode  string        // AUTONOMOUS, ASSISTED, MANUAL_ONLY
	Factors        []TrustFactor // The "Why" behind the score (for HUD)
}

type TrustFactor struct {
	Component   string
	Probability float64 // The confidence contributed by this component
	Weight      float64 // How critical this component is to the system
	Reason      string
}

// TrustEvaluator holds the Bayesian priors and thresholds.
type TrustEvaluator struct {
	MinThreshold float64 // Below this, autonomy is forbidden
}

// Evaluate performs the Bayesian update cycle to determine system trust.
func (te *TrustEvaluator) Evaluate(env *configStruct.EnvConfig) *TrustDescriptor {
	factors := []TrustFactor{}
	
	// 1. PRIOR: Start with a neutral agnostic prior (0.5)
	// We assume the system is trustworthy until proven otherwise, but cautious.
	currentTrust := 0.5 

	// --- EVIDENCE 1: Security Attestation (The Strongest Signal) ---
	// If the binary hash doesn't match, trust plummets.
	attestScore := 0.1 // Default low
	if env.Attestation.Valid {
		if env.Attestation.Level == configStruct.AttestationStrong {
			attestScore = 0.99
		} else if env.Attestation.Level == configStruct.AttestationWeak {
			attestScore = 0.75
		}
	}
	factors = append(factors, TrustFactor{
		Component:   "Security Vault",
		Probability: attestScore,
		Weight:      0.4, // High weight: Security is paramount
		Reason:      fmt.Sprintf("Attestation Level: %s", env.Attestation.Level),
	})

	// --- EVIDENCE 2: Platform Identity Confidence ---
	// How sure are we that this is actually a "Vehicle" or "Drone"?
	platScore := 0.0
	// We extract the confidence from the Q16 value in the config
	if env.Platform.Locked {
		// Find the score for the final selected platform
		for _, candidate := range env.Platform.Candidates {
			if candidate.Class == env.Platform.Final {
				platScore = candidate.Confidence.Float64()
				break
			}
		}
	}
	factors = append(factors, TrustFactor{
		Component:   "Platform Identity",
		Probability: platScore,
		Weight:      0.3, 
		Reason:      fmt.Sprintf("Identified as %s via %s", env.Platform.Final, env.Platform.Source),
	})

	// --- EVIDENCE 3: Hardware Bus Integrity ---
	// Do we have the buses we expect? (e.g., CAN-bus for a car)
	hwScore := 0.5
	if len(env.Hardware.Buses) > 0 {
		// Calculate average confidence of all detected buses
		totalBusConf := 0.0
		for _, b := range env.Hardware.Buses {
			totalBusConf += b.Confidence.Float64()
		}
		hwScore = totalBusConf / float64(len(env.Hardware.Buses))
	}
	factors = append(factors, TrustFactor{
		Component:   "Hardware I/O",
		Probability: hwScore,
		Weight:      0.3,
		Reason:      fmt.Sprintf("%d Active Buses Detected", len(env.Hardware.Buses)),
	})

	// 2. POSTERIOR UPDATE: Recursive Bayesian Update
	// New_Belief = (Likelihood * Prior) / Normalization
	// Here we use a weighted fusion approach which is numerically stable for systems.
	
	numerator := 0.0
	denominator := 0.0

	for _, f := range factors {
		numerator += f.Probability * f.Weight
		denominator += f.Weight
	}

	if denominator > 0 {
		currentTrust = numerator / denominator
	} else {
		currentTrust = 0.0
	}

	// 3. DECISION LOGIC: Determine Operational Mode
	desc := &TrustDescriptor{
		CurrentScore: currentTrust,
		Factors:      factors,
	}

	// Enforce the User's Constraint: Confidence < 0.5 -> Mandatory Manual
	if currentTrust < 0.5 {
		desc.Label = "CRITICAL (UNTRUSTED)"
		desc.OperationMode = "MANUAL_ONLY"
	} else if currentTrust < te.MinThreshold {
		desc.Label = "DEGRADED (CAUTION)"
		desc.OperationMode = "ASSISTED"
	} else {
		desc.Label = "NOMINAL"
		desc.OperationMode = "AUTONOMOUS"
	}

	return desc
}