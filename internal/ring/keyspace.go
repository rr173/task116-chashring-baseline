package ring

import "task116-chashring/internal/model"

func (c *ConsistentHash) Assignments(keys []string) map[string][]string {
	out := map[string][]string{}
	for _, key := range keys {
		node, err := c.Lookup(key)
		if err == nil {
			out[key] = []string{node}
		}
	}
	return out
}
func (c *ConsistentHash) AssignmentCoverage(keys []string) model.RingStats {
	_ = keys
	return c.Stats()
}
