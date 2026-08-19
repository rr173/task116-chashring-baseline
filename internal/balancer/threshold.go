package balancer

import "task116-chashring/internal/model"

func WithinImbalance(stats model.RingStats, threshold float64) bool {
	if threshold < 0 {
		return false
	}
	return Imbalance(stats) <= threshold
}
func Critical(stats model.RingStats) bool { return !WithinImbalance(stats, 0.25) }
