package model

func NormalizeKeys(keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, key)
	}
	return out
}
func KeyCount(keys []string) int { return len(NormalizeKeys(keys)) }
