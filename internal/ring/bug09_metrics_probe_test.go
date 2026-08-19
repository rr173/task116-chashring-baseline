package ring

import (
	"testing"

	"task116-chashring/internal/model"
)

func TestBug09_StatsAndMetricsUsePhysicalNodesAndVirtualPoints(t *testing.T) {
	rh := New("r", model.RingConfig{Replicas: 5}, []*model.Node{{ID: "a", Address: "a:1"}, {ID: "b", Address: "b:1"}})
	stats, metrics := rh.Stats(), rh.Metrics()
	if stats.NodeCount != 2 || stats.PointCount != 10 { t.Fatalf("stats = %+v, want 2 nodes and 10 points", stats) }
	if metrics.Nodes != 2 || metrics.Points != 10 { t.Fatalf("metrics = %+v, want 2 nodes and 10 points", metrics) }
}
