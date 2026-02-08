//MIAUSEproject-founderKJ/multi-platform-AI/core/policy/trust_bayesian.go

package policy

import (
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// TrustDescriptor now includes raw bits for low-level system diagnostic
type TrustDescriptor struct {
	CurrentScore float64 `json:"current_score"` 
	RawScoreQ16  uint16  `json:"raw_score_q16"` // The 0-65535 representation
	Label        string  `json:"label"`
	OperationMode string `json:"operation_mode"`
	Factors      []TrustFactor `json:"factors"`
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
func (te *TrustEvaluator) Evaluate(env *schema.EnvConfig) *TrustDescriptor {
	factors := []TrustFactor{}

	// --- EVIDENCE 1: Security Attestation ---
	// We call the helper and convert back to float for the Bayesian fusion
	rawAttest := te.calculateSecurityWeight(env)
	attestProb := float64(rawAttest) / 65535.0

	factors = append(factors, TrustFactor{
		Component:   "Security Vault",
		Probability: attestProb,
		Weight:      0.4,
		Reason:      fmt.Sprintf("Attestation Level: %s (Raw: %d)", env.Attestation.Level, rawAttest),
	})

	// ... [Identity and Hardware Logic] ...

	// Final Calculation (Weighted Fusion)
	numerator := 0.0
	denominator := 0.0
	for _, f := range factors {
		numerator += f.Probability * f.Weight
		denominator += f.Weight
	}

	finalTrust := numerator / denominator

	return &TrustDescriptor{
		CurrentScore: finalTrust,
		RawScoreQ16:  uint16(finalTrust * 65535), // Seal the final judgment in bits
		Factors:      factors,
		// ... [Labeling Logic] ...
	}
}


// calculateSecurityWeight serves as a translation layer between 
// Security Enums and Bayesian Probabilities.
func (te *TrustEvaluator) calculateSecurityWeight(env *schema.EnvConfig) uint16 {
	var floatScore float64 = 0.1 // Default: Untrusted/Initial state

	if env.Attestation.Valid {
		switch env.Attestation.Level {
		case "strong": // Match your platforms.AttestationStrong
			floatScore = 0.99
		case "weak":   // Match your platforms.AttestationWeak
			floatScore = 0.75
		}
	}

	// Numerical Safety: Ensure we don't exceed 65535
	if floatScore > 1.0 { floatScore = 1.0 }
	return uint16(floatScore * 65535)
}