package model

type Placement struct {
	Key      string   `json:"key"`
	Primary  string   `json:"primary"`
	Replicas []string `json:"replicas"`
}

func UniqueNodes(nodes []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if node != "" && !seen[node] {
			seen[node] = true
			out = append(out, node)
		}
	}
	return out
}
// Capacity reports the number of virtual points a node contributes to a ring.
// When the node does not specify a virtual-node count, it defaults to
// DefaultReplicas — the same default ring.virtualCount applies when generating
// virtual points — so capacity statistics stay consistent with the points
// actually placed on the ring.
func Capacity(n Node) int {
	weight := n.Weight
	if weight <= 0 {
		weight = 1
	}
	virtual := n.VirtualNodes
	if virtual <= 0 {
		virtual = DefaultReplicas
	}
	return weight * virtual
}
