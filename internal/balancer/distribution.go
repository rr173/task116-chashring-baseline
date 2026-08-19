package balancer

import "task116-chashring/internal/model"

type DistributionReport struct {
	Total    int                  `json:"total"`
	Nodes    []model.Distribution `json:"nodes"`
	Balanced bool                 `json:"balanced"`
}

func Report(stats model.RingStats) DistributionReport {
	return DistributionReport{Total: model.TotalPoints(stats), Nodes: model.Shares(stats), Balanced: stats.Balanced}
}
func Imbalance(stats model.RingStats) float64 {
	if len(stats.Nodes) == 0 {
		return 0
	}
	min, max := 1.0, 0.0
	for _, node := range stats.Nodes {
		if node.PointShare < min {
			min = node.PointShare
		}
		if node.PointShare > max {
			max = node.PointShare
		}
	}
	return max - min
}
