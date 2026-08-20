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

func TestNormalizeKeys(t *testing.T) {
	got := NormalizeKeys([]string{"alpha", " alpha ", "", "alpha", "\t", "beta"})
	want := []string{"alpha", "beta"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
	if KeyCount([]string{"alpha", "alpha", " "}) != 1 {
		t.Fatalf("KeyCount should report unique normalized keys")
	}
}
