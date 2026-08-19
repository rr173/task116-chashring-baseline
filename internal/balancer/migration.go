package balancer

import "task116-chashring/internal/model"

type MigrationSummary struct {
	Moves        int      `json:"moves"`
	Added        int      `json:"added"`
	Removed      int      `json:"removed"`
	AffectedKeys []string `json:"affected_keys"`
}

func Summarize(plan Plan) MigrationSummary {
	keys := make([]string, 0, len(plan.Moves))
	for _, move := range plan.Moves {
		keys = append(keys, move.Key)
	}
	return MigrationSummary{Moves: len(plan.Moves), Added: len(plan.Added), Removed: len(plan.Removed), AffectedKeys: keys}
}
func ValidNodeSet(nodes []model.Node) bool {
	seen := map[string]bool{}
	for _, node := range nodes {
		if node.ID == "" || seen[node.ID] {
			return false
		}
		seen[node.ID] = true
	}
	return true
}
