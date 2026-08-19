package balancer

import "task116-chashring/internal/model"

type WeightReport struct {
	Total int            `json:"total"`
	Nodes map[string]int `json:"nodes"`
}

func Weights(nodes []model.Node) WeightReport {
	out := WeightReport{Nodes: map[string]int{}}
	for _, node := range nodes {
		out.Nodes[node.ID] = model.Capacity(node)
		out.Total += node.Weight
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
