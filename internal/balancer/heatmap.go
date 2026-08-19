package balancer

import "task116-chashring/internal/model"

type Heat struct {
	NodeID   string  `json:"node_id"`
	Expected float64 `json:"expected"`
	Actual   float64 `json:"actual"`
	Delta    float64 `json:"delta"`
}

func Heatmap(stats model.RingStats) []Heat {
	out := make([]Heat, 0, len(stats.Nodes))
	total := 0
	for _, node := range stats.Nodes {
		total += node.Weight
	}
	if total == 0 {
		return out
	}
	for _, node := range stats.Nodes {
		expected := float64(node.Weight) / float64(total)
		out = append(out, Heat{NodeID: node.ID, Expected: expected, Actual: node.PointShare, Delta: node.PointShare - expected})
	}
	return out
}
