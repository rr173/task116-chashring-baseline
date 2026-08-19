package store

import (
	"context"
	"fmt"
	"task116-chashring/internal/model"
)

func (s *Store) ValidateAll(ctx context.Context) error {
	rings, nodes, err := s.LoadAll(ctx)
	if err != nil {
		return err
	}
	for id, ring := range rings {
		if err := ring.Validate(); err != nil {
			return fmt.Errorf("ring %s: %w", id, err)
		}
		for _, node := range nodes[id] {
			if err := node.Validate(); err != nil {
				return fmt.Errorf("node %s: %w", node.ID, err)
			}
		}
	}
	return nil
}

var _ = model.ErrRingNotFound
