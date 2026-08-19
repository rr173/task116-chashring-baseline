package ring

import "task116-chashring/internal/model"

type Snapshot struct {
	RingID string           `json:"ring_id"`
	Config model.RingConfig `json:"config"`
	Nodes  []model.Node     `json:"nodes"`
	Stats  model.RingStats  `json:"stats"`
}

func (c *ConsistentHash) Snapshot() Snapshot {
	return Snapshot{RingID: c.ringID, Config: c.Config(), Nodes: c.Nodes(), Stats: c.Stats()}
}
