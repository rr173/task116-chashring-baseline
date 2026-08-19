package model

import "fmt"

func (c RingConfig) Validate() error {
	if c.Replicas <= 0 {
		return fmt.Errorf("replicas must be positive")
	}
	if c.HashFunc != "fnv1a" {
		return fmt.Errorf("unsupported hash function %q", c.HashFunc)
	}
	if c.Replication <= 0 {
		return fmt.Errorf("replication must be positive")
	}
	return nil
}
func (n Node) Validate() error {
	if n.ID == "" || n.Address == "" {
		return fmt.Errorf("node id and address required")
	}
	if n.Weight <= 0 {
		return fmt.Errorf("node weight must be positive")
	}
	return nil
}
func (r Ring) Validate() error {
	if r.ID == "" || r.Name == "" {
		return ErrInvalidConfig
	}
	return r.Config.Normalize().Validate()
}
