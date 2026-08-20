// Package store persists rings and nodes in SQLite using the pure-Go driver.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"

	"task116-chashring/internal/model"
)

// Store is a thin persistence layer over a SQLite database.
type Store struct {
	db *sql.DB
}

// Open opens (and migrates) a SQLite database at the given path.
func Open(path string) (*Store, error) {
	dsn := fmt.Sprintf(
		"file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)",
		path,
	)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("chashring: open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("chashring: ping sqlite: %w", err)
	}
	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

// Close releases the underlying database handle.
func (s *Store) Close() error { return s.db.Close() }

func migrate(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS rings (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    replicas INTEGER NOT NULL,
    hash_func TEXT NOT NULL,
    replication INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS nodes (
    ring_id TEXT NOT NULL,
    id TEXT NOT NULL,
    address TEXT NOT NULL,
    weight INTEGER NOT NULL,
    virtual_nodes INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (ring_id, id)
);`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("chashring: migrate: %w", err)
	}
	return nil
}

// SaveRing upserts a ring row.
func (s *Store) SaveRing(ctx context.Context, r model.Ring) error {
	if r.ID == "" {
		return model.ErrInvalidConfig
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO rings (id,name,replicas,hash_func,replication,created_at,updated_at)
         VALUES (?,?,?,?,?,?,?)
         ON CONFLICT(id) DO UPDATE SET name=excluded.name, replicas=excluded.replicas,
           hash_func=excluded.hash_func, replication=excluded.replication, updated_at=excluded.updated_at`,
		r.ID, r.Name, r.Config.Replicas, r.Config.HashFunc, r.Config.Replication, r.CreatedAt, r.UpdatedAt)
	if err != nil {
		return fmt.Errorf("chashring: save ring: %w", err)
	}
	return nil
}

// GetRing loads a single ring by id.
func (s *Store) GetRing(ctx context.Context, id string) (model.Ring, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id,name,replicas,hash_func,replication,created_at,updated_at FROM rings WHERE id=?`, id)
	return scanRing(row)
}

func scanRing(row *sql.Row) (model.Ring, error) {
	var r model.Ring
	if err := row.Scan(&r.ID, &r.Name, &r.Config.Replicas, &r.Config.HashFunc, &r.Config.Replication, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Ring{}, model.ErrRingNotFound
		}
		return model.Ring{}, fmt.Errorf("chashring: scan ring: %w", err)
	}
	return r, nil
}

// ListRings returns all rings ordered by creation time.
func (s *Store) ListRings(ctx context.Context) ([]model.Ring, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id,name,replicas,hash_func,replication,created_at,updated_at FROM rings ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("chashring: list rings: %w", err)
	}
	defer rows.Close()
	out := []model.Ring{}
	for rows.Next() {
		r, err := scanRingRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func scanRingRow(rows *sql.Rows) (model.Ring, error) {
	var r model.Ring
	if err := rows.Scan(&r.ID, &r.Name, &r.Config.Replicas, &r.Config.HashFunc, &r.Config.Replication, &r.CreatedAt, &r.UpdatedAt); err != nil {
		return model.Ring{}, fmt.Errorf("chashring: scan ring row: %w", err)
	}
	return r, nil
}

// DeleteRing removes a ring and its nodes.
func (s *Store) DeleteRing(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM nodes WHERE ring_id=?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM rings WHERE id=?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

// SaveNode inserts a node, returning ErrNodeDuplicate on a primary-key clash.
func (s *Store) SaveNode(ctx context.Context, ringID string, n model.Node) error {
	if ringID == "" || n.ID == "" {
		return model.ErrInvalidConfig
	}
	var exists int
	if err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM nodes WHERE ring_id=? AND id=?`, ringID, n.ID).Scan(&exists); err != nil {
		return fmt.Errorf("chashring: check node: %w", err)
	}
	if exists > 0 {
		return model.ErrNodeDuplicate
	}
	if _, err := s.db.ExecContext(ctx,
		`INSERT INTO nodes (ring_id,id,address,weight,virtual_nodes,created_at) VALUES (?,?,?,?,?,?)`,
		ringID, n.ID, n.Address, n.Weight, n.VirtualNodes, n.CreatedAt); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return model.ErrNodeDuplicate
		}
		return fmt.Errorf("chashring: save node: %w", err)
	}
	return nil
}

// GetNode loads a single node.
func (s *Store) GetNode(ctx context.Context, ringID, id string) (model.Node, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id,address,weight,virtual_nodes,created_at FROM nodes WHERE ring_id=? AND id=?`, ringID, id)
	var n model.Node
	if err := row.Scan(&n.ID, &n.Address, &n.Weight, &n.VirtualNodes, &n.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Node{}, model.ErrNodeNotFound
		}
		return model.Node{}, fmt.Errorf("chashring: scan node: %w", err)
	}
	return n, nil
}

// ListNodes returns all nodes of a ring ordered by creation time.
func (s *Store) ListNodes(ctx context.Context, ringID string) ([]model.Node, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id,address,weight,virtual_nodes,created_at FROM nodes WHERE ring_id=? ORDER BY created_at`, ringID)
	if err != nil {
		return nil, fmt.Errorf("chashring: list nodes: %w", err)
	}
	defer rows.Close()
	out := []model.Node{}
	for rows.Next() {
		var n model.Node
		if err := rows.Scan(&n.ID, &n.Address, &n.Weight, &n.VirtualNodes, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("chashring: scan node row: %w", err)
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// DeleteNode removes a single node.
func (s *Store) DeleteNode(ctx context.Context, ringID, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM nodes WHERE ring_id=? AND id=?`, ringID, id)
	if err != nil {
		return fmt.Errorf("chashring: delete node: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return model.ErrNodeNotFound
	}
	return nil
}

// LoadAll returns all rings and their nodes keyed by ring id.
func (s *Store) LoadAll(ctx context.Context) (map[string]model.Ring, map[string][]model.Node, error) {
	rings, err := s.ListRings(ctx)
	if err != nil {
		return nil, nil, err
	}
	ringMap := make(map[string]model.Ring, len(rings))
	nodeMap := make(map[string][]model.Node, len(rings))
	for _, r := range rings {
		ringMap[r.ID] = r
		nodes, err := s.ListNodes(ctx, r.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("chashring: load nodes for %s: %w", r.ID, err)
		}
		nodeMap[r.ID] = nodes
	}
	return ringMap, nodeMap, nil
}
