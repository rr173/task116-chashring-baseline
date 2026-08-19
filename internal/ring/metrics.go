package ring

import "task116-chashring/internal/model"

type Metrics struct {
	Nodes    int     `json:"nodes"`
	Points   int     `json:"points"`
	Capacity int     `json:"capacity"`
	Balance  float64 `json:"balance"`
}

func (c *ConsistentHash) Metrics() Metrics {
	stats := c.Stats()
	nodes := c.Nodes()
	balance := 0.0
	if len(stats.Nodes) > 0 {
		min, max := 1.0, 0.0
		for _, node := range stats.Nodes {
			if node.PointShare < min {
				min = node.PointShare
			}
			if node.PointShare > max {
				max = node.PointShare
			}
		}
		balance = max - min
	}
	return Metrics{Nodes: len(nodes), Points: stats.PointCount, Capacity: TotalCapacity(nodes), Balance: balance}
}

var _ = model.RingStats{}
