package ring

import (
	"testing"

	"task116-chashring/internal/model"
)

func TestBug03_PublicFingerprintMatchesRingPlacementHash(t *testing.T) {
	key := "tenant-42/document-7"
	if got, want := model.KeyFingerprint(key), hashPoint(key); got != want {
		t.Fatalf("fingerprint = %d, ring hash = %d", got, want)
	}
}
