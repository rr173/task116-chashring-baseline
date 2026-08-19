// Package balancer computes the effect of topology changes on key placement.
package balancer

import (
	"task116-chashring/internal/model"
	"task116-chashring/internal/ring"
)

// Move describes a key whose responsible node changes after a topology change.
type Move struct {
	Key  string `json:"key"`
	From string `json:"from"`
	To   string `json:"to"`
}

// Plan is the result of a rebalance preview.
type Plan struct {
	Added   []model.Node `json:"added"`
	Removed []string     `json:"removed"`
	Moves   []Move       `json:"moves"`
	Count   int          `json:"move_count"`
}

// Preview estimates how key assignments shift if `added` nodes are inserted and
// `removed` node ids are deleted, evaluated over a sample of keys.
func Preview(cfg model.RingConfig, current []model.Node, sample []string, added []model.Node, removed []string) Plan {
	rm := make(map[string]bool, len(removed))
	for _, id := range removed {
		rm[id] = true
	}
	candNodes := make([]model.Node, 0, len(current)+len(added))
	for _, n := range current {
		candNodes = append(candNodes, n)
	}
	candNodes = append(candNodes, added...)

	base := ring.New("preview", cfg, toPtrs(current))
	cand := ring.New("preview", cfg, toPtrs(candNodes))

	plan := Plan{Added: added, Removed: removed}
	for _, key := range sample {
		before, _ := base.Lookup(key)
		after, _ := cand.Lookup(key)
		if before != after {
			plan.Moves = append(plan.Moves, Move{Key: key, From: before, To: after})
		}
	}
	plan.Count = len(plan.Moves)
	return plan
}

func toPtrs(in []model.Node) []*model.Node {
	out := make([]*model.Node, len(in))
	for i := range in {
		out[i] = &in[i]
	}
	return out
}
