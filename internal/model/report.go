package model

type Distribution struct {
	NodeID string  `json:"node_id"`
	Points int     `json:"points"`
	Share  float64 `json:"share"`
}

func Shares(stats RingStats) []Distribution {
	out := make([]Distribution, 0, len(stats.Nodes))
	for _, node := range stats.Nodes {
		out = append(out, Distribution{NodeID: node.ID, Points: node.VirtualNodes, Share: node.PointShare})
	}
	return out
}
func TotalPoints(stats RingStats) int {
	total := 0
	for _, node := range stats.Nodes {
		total += node.VirtualNodes
	}
	return total
}
