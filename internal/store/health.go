package store

import (
	"context"
	"task116-chashring/internal/model"
)

type Health struct {
	Rings int `json:"rings"`
	Nodes int `json:"nodes"`
}

func (s *Store) Health(ctx context.Context) (Health, error) {
	var out Health
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM rings`).Scan(&out.Rings); err != nil {
		return out, err
	}
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM nodes`).Scan(&out.Nodes); err != nil {
		return out, err
	}
	return out, nil
}
func (s *Store) ValidateRing(ctx context.Context, id string) error {
	r, err := s.GetRing(ctx, id)
	if err != nil {
		return err
	}
	if err := r.Validate(); err != nil {
		return err
	}
	nodes, err := s.ListNodes(ctx, id)
	if err != nil {
		return err
	}
	for _, n := range nodes {
		if err := n.Validate(); err != nil {
			return err
		}
	}
	return nil
}

var _ = model.ErrInvalidConfig
