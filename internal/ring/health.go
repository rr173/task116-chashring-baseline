package ring

import "task116-chashring/internal/model"

type Health struct {
	RingID   string `json:"ring_id"`
	Nodes    int    `json:"nodes"`
	Points   int    `json:"points"`
	Empty    bool   `json:"empty"`
	Balanced bool   `json:"balanced"`
}

func (c *ConsistentHash) Health() Health {
	stats := c.Stats()
	return Health{RingID: c.ringID, Nodes: stats.NodeCount, Points: stats.PointCount, Empty: stats.NodeCount == 0, Balanced: stats.Balanced}
}
func (c *ConsistentHash) HasNode(id string) bool {
	for _, node := range c.Nodes() {
		if node.ID == id {
			return true
		}
	}
	return false
}

var _ = model.ErrEmptyRing
