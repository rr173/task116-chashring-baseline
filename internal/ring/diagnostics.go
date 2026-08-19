package ring

import "task116-chashring/internal/model"

type Diagnostic struct {
	Health       Health               `json:"health"`
	Distribution []model.Distribution `json:"distribution"`
	Capacities   map[string]int       `json:"capacities"`
}

func (c *ConsistentHash) Diagnostic() Diagnostic {
	stats := c.Stats()
	nodes := c.Nodes()
	return Diagnostic{Health: c.Health(), Distribution: model.Shares(stats), Capacities: NodeCapacities(nodes)}
}
