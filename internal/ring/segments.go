package ring

type Segment struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
	Node  string `json:"node"`
}

func (c *ConsistentHash) SegmentCount() int { c.mu.RLock(); defer c.mu.RUnlock(); return len(c.points) }
func (c *ConsistentHash) Empty() bool       { return c.SegmentCount() == 0 }
