//core/policy/scheduler.go

package policy

type PowerProfile struct {
	BatteryLevel float64 // 0.0 to 1.0
	IsCharging   bool
	EstimatedRun uint32 // Minutes remaining
}

func (e *TrustEvaluator) CalculateUtility(taskRequirement float64, p PowerProfile) float64 {
	// Utility = (CurrentTrust * BatteryLevel) / TaskDifficulty
	// High trust + High battery = High Utility for the swarm
	if p.BatteryLevel < 0.15 && !p.IsCharging {
		return 0.0 // Node is too weak to help
	}
	return (e.LastScore * p.BatteryLevel) / taskRequirement
}
