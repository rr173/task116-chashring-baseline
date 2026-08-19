package ring

type HashRange struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
	Wraps bool   `json:"wraps"`
}

func RangeContains(r HashRange, value uint64) bool {
	if !r.Wraps {
		return value >= r.Start && value < r.End
	}
	return value >= r.Start || value < r.End
}
func (c *ConsistentHash) HashRangeFor(key string) (HashRange, error) {
	idx, err := c.lookupIndex(key)
	if err != nil {
		return HashRange{}, err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	point := c.points[idx]
	next := c.points[(idx+1)%len(c.points)]
	return HashRange{Start: point.hash, End: next.hash, Wraps: next.hash <= point.hash}, nil
}
