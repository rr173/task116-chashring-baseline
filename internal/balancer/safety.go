package balancer

import (
	"fmt"
	"task116-chashring/internal/model"
)

func ValidatePlan(plan Plan) error {
	seen := map[string]bool{}
	for _, node := range plan.Added {
		if err := node.Validate(); err != nil {
			return err
		}
		if seen[node.ID] {
			return fmt.Errorf("duplicate added node %s", node.ID)
		}
		seen[node.ID] = true
	}
	for _, id := range plan.Removed {
		if id == "" {
			return model.ErrNodeNotFound
		}
	}
	return nil
}
