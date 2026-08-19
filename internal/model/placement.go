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
func Capacity(n Node) int {
	weight := n.Weight
	if weight <= 0 {
		weight = 1
	}
	virtual := n.VirtualNodes
	if virtual <= 0 {
		virtual = 1
	}
	return weight * virtual
}
