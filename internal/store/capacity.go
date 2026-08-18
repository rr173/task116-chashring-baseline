package store

import (
	"context"
	"task116-chashring/internal/model"
)

func (s *Store) NodeCapacities(ctx context.Context, ringID string) (map[string]int, error) {
	nodes, err := s.ListNodes(ctx, ringID)
	if err != nil {
		return nil, err
	}
	out := map[string]int{}
	for _, n := range nodes {
		out[n.ID] = model.Capacity(n)
	}
	return out, nil
}
