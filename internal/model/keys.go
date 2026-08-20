package model

import "strings"

// NormalizeKeys puts incoming keys into a canonical form: surrounding
// whitespace is trimmed, keys that are empty afterwards are dropped, and
// duplicates collapse to their first occurrence. " alpha ", "alpha", and a
// repeated "alpha" therefore all refer to the same key.
func NormalizeKeys(keys []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		k := strings.TrimSpace(key)
		if k == "" || seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, k)
	}
	return out
}

func KeyCount(keys []string) int { return len(NormalizeKeys(keys)) }
