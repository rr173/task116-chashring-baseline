package api

import (
	"context"
	"testing"

	"task116-chashring/internal/model"
)

func TestBug05_DefaultRingConfigSurvivesPersistenceAndLiveIndex(t *testing.T) {
	a := newTestAPI(t)
	if err := a.CreateRing(context.Background(), model.Ring{ID: "defaults", Name: "defaults"}); err != nil { t.Fatal(err) }
	persisted, err := a.Store().GetRing(context.Background(), "defaults")
	if err != nil { t.Fatal(err) }
	if persisted.Config.Replicas != 100 { t.Fatalf("persisted replicas = %d, want 100", persisted.Config.Replicas) }
	live, ok := a.Ring("defaults")
	if !ok || live.Config().Replicas != 100 { t.Fatalf("live config = %+v", live.Config()) }
}
