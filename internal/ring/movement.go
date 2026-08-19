package ring

import "task116-chashring/internal/model"

type Movement struct {
	Key  string `json:"key"`
	From string `json:"from"`
	To   string `json:"to"`
}

func (c *ConsistentHash) Compare(keys []string, other *ConsistentHash) []Movement {
	out := make([]Movement, 0)
	for _, key := range keys {
		from, _ := c.Lookup(key)
		to, _ := other.Lookup(key)
		if from != to {
			out = append(out, Movement{Key: key, From: from, To: to})
		}
	}
	return out
}
func (c *ConsistentHash) Placement(key string, replicas int) (model.Placement, error) {
	primary, err := c.Lookup(key)
	if err != nil {
		return model.Placement{}, err
	}
	reps, err := c.Replicas(key, replicas)
	if err != nil {
		return model.Placement{}, err
	}
	return model.Placement{Key: key, Primary: primary, Replicas: model.UniqueNodes(reps)}, nil
}
