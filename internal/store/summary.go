package store

import (
	"context"
	"task116-chashring/internal/model"
)

// RingSummary is a per-ring summary surfaced by the /summary endpoint for
// monitoring. Nodes is the number of distinct physical backends; Points is the
// total number of virtual routing points the ring would place for those nodes
// (weight * virtual_nodes per node, honoring the same defaults as the ring
// builder). The two are tracked separately so monitoring can distinguish real
// node scale from routing-point scale.
type RingSummary struct {
	Ring     model.Ring `json:"ring"`
	Nodes    int        `json:"nodes"`
	Points   int        `json:"points"`
	Capacity int        `json:"capacity"`
}

func (s *Store) Summary(ctx context.Context, id string) (RingSummary, error) {
	r, err := s.GetRing(ctx, id)
	if err != nil {
		return RingSummary{}, err
	}
	nodes, err := s.ListNodes(ctx, id)
	if err != nil {
		return RingSummary{}, err
	}
	capacity := 0
	for _, node := range nodes {
		capacity += model.Capacity(node)
	}
	return RingSummary{Ring: r, Nodes: len(nodes), Points: capacity, Capacity: capacity}, nil
}
