package balancer

import "task116-chashring/internal/model"

type WeightReport struct {
	Total int            `json:"total"`
	Nodes map[string]int `json:"nodes"`
}

// Weights summarizes per-node capacity and the ring's total capacity. The
// total is the sum of node capacities (not raw weights) so that WeightShare
// expresses each node's fraction of the ring's virtual-point capacity and the
// reported total matches the virtual points generated on the ring.
func Weights(nodes []model.Node) WeightReport {
	out := WeightReport{Nodes: map[string]int{}}
	for _, node := range nodes {
		cap := model.Capacity(node)
		out.Nodes[node.ID] = cap
		out.Total += cap
	}
	return out
}
func WeightShare(nodes []model.Node, id string) float64 {
	report := Weights(nodes)
	if report.Total == 0 {
		return 0
	}
	return float64(report.Nodes[id]) / float64(report.Total)
}
