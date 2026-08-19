package model

func NormalizeKeys(keys []string) []string {
	out := make([]string, 0, len(keys))
	seen := map[string]bool{}
	for _, key := range keys {
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, key)
		}
	}
	return out
}
func KeyCount(keys []string) int { return len(NormalizeKeys(keys)) }
