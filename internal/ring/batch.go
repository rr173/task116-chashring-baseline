package ring

import "task116-chashring/internal/model"

type LookupResult struct {
	Key   string `json:"key"`
	Node  string `json:"node"`
	Error string `json:"error,omitempty"`
}

func (c *ConsistentHash) LookupMany(keys []string) []LookupResult {
	out := make([]LookupResult, 0, len(keys))
	for _, key := range keys {
		node, err := c.Lookup(key)
		result := LookupResult{Key: key, Node: node}
		if err != nil {
			result.Error = err.Error()
		}
		out = append(out, result)
	}
	return out
}
func (c *ConsistentHash) PlacementMany(keys []string, replicas int) []model.Placement {
	out := make([]model.Placement, 0, len(keys))
	for _, key := range keys {
		placement, err := c.Placement(key, replicas)
		if err == nil {
			out = append(out, placement)
		}
	}
	return out
}
