package model

import "hash/fnv"

func KeyFingerprint(key string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return h.Sum64()
}
func SameBucket(a, b uint64, buckets uint64) bool {
	if buckets == 0 {
		return false
	}
	return a%buckets == b%buckets
}
