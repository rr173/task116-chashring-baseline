package balancer

import "task116-chashring/internal/model"

type Coverage struct {
	Assigned int `json:"assigned"`
	Missing  int `json:"missing"`
	Distinct int `json:"distinct"`
}

func CoverageOf(assignments []string) Coverage {
	seen := map[string]bool{}
	out := Coverage{}
	for _, id := range assignments {
		if id == "" {
			out.Missing++
			continue
		}
		out.Assigned++
		seen[id] = true
	}
	out.Distinct = len(seen)
	return out
}
func NodeOrder(nodes []model.Node) []string {
	out := make([]string, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, node.ID)
	}
	return out
}
