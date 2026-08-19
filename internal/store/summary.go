package store

import (
	"context"
	"task116-chashring/internal/model"
)

type RingSummary struct {
	Ring     model.Ring `json:"ring"`
	Nodes    int        `json:"nodes"`
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
	return RingSummary{Ring: r, Nodes: len(nodes), Capacity: capacity}, nil
}
