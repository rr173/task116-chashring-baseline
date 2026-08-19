package ring

import (
	"testing"

	"task116-chashring/internal/model"
)

func TestBug01_DefaultConfigAndExplicitVirtualNodesKeepTheirOwnCapacity(t *testing.T) {
	defaults := New("defaults", model.RingConfig{}, []*model.Node{{ID: "a", Weight: 1}})
	if got := defaults.Config().Replicas; got != 100 {
		t.Fatalf("default replicas = %d, want 100", got)
	}
	explicit := New("explicit", model.RingConfig{Replicas: 100}, []*model.Node{{ID: "b", Weight: 1, VirtualNodes: 3}})
	if got := explicit.Stats().PointCount; got != 3 {
		t.Fatalf("explicit virtual-node capacity = %d, want 3", got)
	}
}
