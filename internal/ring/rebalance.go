package ring

import "task116-chashring/internal/model"

type RebalanceReport struct {
	Before model.RingStats `json:"before"`
	After  model.RingStats `json:"after"`
	Moves  []Movement      `json:"moves"`
}

func PreviewRebalance(current *ConsistentHash, next *ConsistentHash, keys []string) RebalanceReport {
	return RebalanceReport{Before: current.Stats(), After: next.Stats(), Moves: current.Compare(keys, next)}
}
