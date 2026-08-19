package ring

import "task116-chashring/internal/model"

func (c *ConsistentHash) ForEachNode(fn func(model.Node) bool) {
	if fn == nil {
		return
	}
	for _, node := range c.Nodes() {
		if !fn(node) {
			return
		}
	}
}
func (c *ConsistentHash) NodeIDs() []string {
	out := []string{}
	c.ForEachNode(func(node model.Node) bool { out = append(out, node.ID); return true })
	return out
}
