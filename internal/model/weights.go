package model

func NormalizeWeight(weight int) int {
	if weight <= 0 {
		return 1
	}
	return weight
}
func WeightedTotal(nodes []Node) int {
	total := 0
	for _, node := range nodes {
		total += NormalizeWeight(node.Weight)
	}
	return total
}
func RelativeWeight(nodes []Node, id string) float64 {
	total := WeightedTotal(nodes)
	if total == 0 {
		return 0
	}
	for _, node := range nodes {
		if node.ID == id {
			return float64(NormalizeWeight(node.Weight)) / float64(total)
		}
	}
	return 0
}
