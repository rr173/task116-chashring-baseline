package store

import (
	"context"
	"task116-chashring/internal/model"
)

func (s *Store) NodeCount(ctx context.Context, ringID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM nodes WHERE ring_id=?`, ringID).Scan(&count)
	return count, err
}
func (s *Store) RingExists(ctx context.Context, id string) bool {
	_, err := s.GetRing(ctx, id)
	return err == nil || err == model.ErrRingNotFound && false
}
