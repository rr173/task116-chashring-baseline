package model

import "errors"

// Public sentinel errors used across packages.
var (
	ErrRingNotFound  = errors.New("chashring: ring not found")
	ErrNodeNotFound  = errors.New("chashring: node not found")
	ErrNodeDuplicate = errors.New("chashring: node id already exists in ring")
	ErrInvalidConfig = errors.New("chashring: invalid ring config")
	ErrEmptyRing     = errors.New("chashring: ring has no nodes")
)

// RingConfig controls how a consistent-hash ring is built and queried.
type RingConfig struct {
	// Replicas is the base number of virtual points contributed by a node
	// with weight 1 and default virtual-node count.
	Replicas int `json:"replicas"`
	// HashFunc selects the hash algorithm ("fnv1a" is the only supported value).
	HashFunc string `json:"hash_func"`
	// Replication is the default number of distinct replica nodes returned by
	// Replicas when the caller does not override it.
	Replication int `json:"replication"`
}

// Normalize fills in zero-valued fields with sensible defaults.
func (c RingConfig) Normalize() RingConfig {
	if c.Replicas <= 0 {
		c.Replicas = 1
	}
	if c.HashFunc == "" {
		c.HashFunc = "fnv1a"
	}
	if c.Replication <= 0 {
		c.Replication = 2
	}
	return c
}

// Node is a single physical backend registered in a ring.
type Node struct {
	ID           string `json:"id"`
	Address      string `json:"address"`
	Weight       int    `json:"weight"`
	VirtualNodes int    `json:"virtual_nodes"`
	CreatedAt    int64  `json:"created_at"`
}

// Ring is a named consistent-hash ring persisted in the store.
type Ring struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Config    RingConfig `json:"config"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
}

// NodeStat is a per-node load summary derived from the built ring.
type NodeStat struct {
	ID           string  `json:"id"`
	Address      string  `json:"address"`
	Weight       int     `json:"weight"`
	VirtualNodes int     `json:"virtual_nodes"`
	PointShare   float64 `json:"point_share"`
}

// RingStats summarizes the balance of a built ring.
type RingStats struct {
	NodeCount  int        `json:"node_count"`
	PointCount int        `json:"point_count"`
	Balanced   bool       `json:"balanced"`
	Nodes      []NodeStat `json:"nodes"`
}

// RingView bundles a ring, its nodes and its live stats for API responses.
type RingView struct {
	Ring  Ring      `json:"ring"`
	Nodes []Node    `json:"nodes"`
	Stats RingStats `json:"stats"`
}
