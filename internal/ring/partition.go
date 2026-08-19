package ring

import "task116-chashring/internal/model"

type Partition struct {
	Node string   `json:"node"`
	Keys []string `json:"keys"`
}

func (c *ConsistentHash) Partition(keys []string) []Partition {
	assignments := c.Assignments(model.NormalizeKeys(keys))
	groups := map[string][]string{}
	for key, nodes := range assignments {
		if len(nodes) > 0 {
			groups[nodes[0]] = append(groups[nodes[0]], key)
		}
	}
	out := make([]Partition, 0, len(groups))
	for node, items := range groups {
		out = append(out, Partition{Node: node, Keys: items})
	}
	return out
}
