// Package ring implements a weighted consistent-hash ring with virtual nodes.
package ring

import (
	"errors"
	"fmt"
	"hash/fnv"
	"sort"
	"sync"

	"task116-chashring/internal/model"
)

// ErrInvalidReplicas is returned when a non-negative replica count is required.
var ErrInvalidReplicas = errors.New("chashring: invalid replicas count")

type point struct {
	hash   uint64
	nodeID string
}

// ConsistentHash is an in-memory consistent-hash ring for one logical ring.
type ConsistentHash struct {
	mu     sync.RWMutex
	ringID string
	cfg    model.RingConfig
	nodes  map[string]*model.Node
	points []point
}

// New builds a ring from a config and an initial set of nodes.
func New(ringID string, cfg model.RingConfig, nodes []*model.Node) *ConsistentHash {
	c := &ConsistentHash{
		ringID: ringID,
		cfg:    cfg.Normalize(),
		nodes:  make(map[string]*model.Node, len(nodes)),
	}
	for _, n := range nodes {
		if n != nil {
			c.nodes[n.ID] = n
		}
	}
	c.rebuild()
	return c
}

func (c *ConsistentHash) virtualCount(n *model.Node) int {
	base := c.cfg.Replicas
	if base <= 0 {
		base = 100
	}
	vn := n.VirtualNodes
	if vn <= 0 {
		vn = base
	}
	w := n.Weight
	if w <= 0 {
		w = 1
	}
	return vn * w
}

// rebuild recomputes the sorted point list from the current node set.
// Callers must hold c.mu (write lock).
func (c *ConsistentHash) rebuild() {
	pts := make([]point, 0, len(c.nodes)*c.cfg.Replicas*4+1)
	for _, n := range c.nodes {
		vcount := c.virtualCount(n)
		for i := 0; i < vcount; i++ {
			pts = append(pts, point{hash: hashPoint(fmt.Sprintf("%s:%d", n.ID, i)), nodeID: n.ID})
		}
	}
	sort.Slice(pts, func(i, j int) bool { return pts[i].hash < pts[j].hash })
	c.points = pts
}

func hashPoint(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}

// Config returns a copy of the ring configuration.
func (c *ConsistentHash) Config() model.RingConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cfg
}

// Nodes returns the registered nodes sorted by ID.
func (c *ConsistentHash) Nodes() []model.Node {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]model.Node, 0, len(c.nodes))
	for _, n := range c.nodes {
		out = append(out, *n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// AddNode registers a node and rebuilds the ring.
func (c *ConsistentHash) AddNode(n *model.Node) error {
	if n == nil || n.ID == "" {
		return model.ErrNodeNotFound
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.nodes[n.ID]; ok {
		return model.ErrNodeDuplicate
	}
	c.nodes[n.ID] = n
	c.rebuild()
	return nil
}

// RemoveNode drops a node and rebuilds the ring.
func (c *ConsistentHash) RemoveNode(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.nodes[id]; !ok {
		return model.ErrNodeNotFound
	}
	delete(c.nodes, id)
	c.rebuild()
	return nil
}

// lookupIndex returns the clockwise index of the first point >= hash(key).
func (c *ConsistentHash) lookupIndex(key string) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.points) == 0 {
		return 0, model.ErrEmptyRing
	}
	h := hashPoint(key)
	idx := sort.Search(len(c.points), func(i int) bool { return c.points[i].hash >= h })
	if idx == len(c.points) {
		idx = 0
	}
	return idx, nil
}

// Lookup returns the primary node responsible for key.
func (c *ConsistentHash) Lookup(key string) (string, error) {
	idx, err := c.lookupIndex(key)
	if err != nil {
		return "", err
	}
	return c.points[idx].nodeID, nil
}

// Replicas returns up to n distinct node ids responsible for key, ordered by
// clockwise distance from the primary.
func (c *ConsistentHash) Replicas(key string, n int) ([]string, error) {
	if n < 0 {
		return nil, ErrInvalidReplicas
	}
	if n == 0 {
		return []string{}, nil
	}
	c.mu.RLock()
	m := len(c.points)
	c.mu.RUnlock()
	if m == 0 {
		return []string{}, nil
	}
	idx, err := c.lookupIndex(key)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool, n)
	out := make([]string, 0, n)
	for i := 0; i < m && len(out) < n; i++ {
		cand := c.points[(idx+i)%m].nodeID
		if !seen[cand] {
			seen[cand] = true
			out = append(out, cand)
		}
	}
	return out, nil
}

// Stats computes a per-node virtual-point distribution summary.
func (c *ConsistentHash) Stats() model.RingStats {
	c.mu.RLock()
	pts := c.points
	c.mu.RUnlock()

	per := make(map[string]int, len(pts))
	total := 0
	for _, p := range pts {
		per[p.nodeID]++
		total++
	}
	nodes := make([]model.NodeStat, 0, len(per))
	for id, cnt := range per {
		share := 0.0
		if total > 0 {
			share = float64(cnt) / float64(total)
		}
		nodes = append(nodes, model.NodeStat{ID: id, VirtualNodes: cnt, PointShare: share})
	}
	balanced := true
	if len(nodes) > 1 {
		minS, maxS := 1.0, 0.0
		for _, nd := range nodes {
			if nd.PointShare < minS {
				minS = nd.PointShare
			}
			if nd.PointShare > maxS {
				maxS = nd.PointShare
			}
		}
		balanced = (maxS - minS) <= 0.25
	}
	return model.RingStats{
		NodeCount:  len(per),
		PointCount: total,
		Balanced:   balanced,
		Nodes:      nodes,
	}
}
