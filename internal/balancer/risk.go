package balancer

import "task116-chashring/internal/model"

type Risk struct {
	Critical  bool    `json:"critical"`
	Imbalance float64 `json:"imbalance"`
	Nodes     int     `json:"nodes"`
}

func RiskOf(stats model.RingStats) Risk {
	imbalance := Imbalance(stats)
	return Risk{Critical: imbalance > 0.25, Imbalance: imbalance, Nodes: stats.NodeCount}
}
func RequiresRebalance(stats model.RingStats) bool {
	return RiskOf(stats).Critical || stats.NodeCount == 0
}
