package model

import "testing"

func TestRingConfigNormalize(t *testing.T) {
	c := RingConfig{}.Normalize()
	if c.Replicas != 100 || c.HashFunc != "fnv1a" || c.Replication != 2 {
		t.Fatalf("unexpected defaults: %+v", c)
	}
	c2 := RingConfig{Replicas: 50, HashFunc: "fnv1a", Replication: 3}.Normalize()
	if c2.Replicas != 50 || c2.Replication != 3 {
		t.Fatalf("normalize mutated explicit values: %+v", c2)
	}
}
