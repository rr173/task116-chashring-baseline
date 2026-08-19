package ring

import "task116-chashring/internal/model"

func NodeCapacities(nodes []model.Node) map[string]int {
	out := map[string]int{}
	for _, node := range nodes {
		out[node.ID] = model.Capacity(node)
	}
	return out
}
func TotalCapacity(nodes []model.Node) int {
	total := 0
	for _, node := range nodes {
		total += model.Capacity(node)
	}
	return total
}
