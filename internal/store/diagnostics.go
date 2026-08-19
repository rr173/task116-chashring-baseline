package store

import (
	"context"
	"task116-chashring/internal/model"
)

type Diagnostics struct {
	Health Health       `json:"health"`
	Valid  bool         `json:"valid"`
	Rings  []model.Ring `json:"rings"`
}

func (s *Store) Diagnostics(ctx context.Context) (Diagnostics, error) {
	health, err := s.Health(ctx)
	if err != nil {
		return Diagnostics{}, err
	}
	rings, err := s.ListRings(ctx)
	if err != nil {
		return Diagnostics{}, err
	}
	valid := s.ValidateAll(ctx) == nil
	return Diagnostics{Health: health, Valid: valid, Rings: rings}, nil
}
