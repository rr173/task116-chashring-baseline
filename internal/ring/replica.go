package ring

import "task116-chashring/internal/model"

type ReplicaReport struct {
	Key      string   `json:"key"`
	Primary  string   `json:"primary"`
	Replicas []string `json:"replicas"`
	Distinct int      `json:"distinct"`
}

func (c *ConsistentHash) ReplicaReport(key string, n int) (ReplicaReport, error) {
	placement, err := c.Placement(key, n)
	if err != nil {
		return ReplicaReport{}, err
	}
	return ReplicaReport{Key: key, Primary: placement.Primary, Replicas: placement.Replicas, Distinct: len(model.UniqueNodes(placement.Replicas))}, nil
}
