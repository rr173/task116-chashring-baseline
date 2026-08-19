package balancer

import "task116-chashring/internal/model"

type Thresholds struct {
	Warning  float64 `json:"warning"`
	Critical float64 `json:"critical"`
}

func DefaultThresholds() Thresholds { return Thresholds{Warning: 0.15, Critical: 0.25} }
func Evaluate(stats model.RingStats, t Thresholds) string {
	imbalance := Imbalance(stats)
	if imbalance >= t.Critical {
		return "critical"
	}
	if imbalance >= t.Warning {
		return "warning"
	}
	return "healthy"
}
