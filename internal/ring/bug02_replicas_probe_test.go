package ring

import (
	"testing"

	"task116-chashring/internal/model"
)

func TestBug02_ReplicaResultsRemainDistinctAcrossRingAndPlacement(t *testing.T) {
	rh := New("r", model.RingConfig{Replicas: 4}, sampleNodes())
	p, err := rh.Placement("tenant-42", 3)
	if err != nil { t.Fatal(err) }
	if len(p.Replicas) != 3 || len(model.UniqueNodes(p.Replicas)) != 3 { t.Fatalf("placement replicas = %v, want three distinct nodes", p.Replicas) }
	if got := model.UniqueNodes([]string{"a", "a", "b"}); len(got) != 2 { t.Fatalf("unique helper returned %v", got) }
}
