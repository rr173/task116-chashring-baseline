package ring

func NormalizeReplicaCount(requested, available int) int {
	if requested < 0 {
		return 0
	}
	if available >= 0 && requested > available {
		return available
	}
	return requested
}

func ReplicaSufficient(replicas []string, requested int) bool {
	return len(replicas) >= NormalizeReplicaCount(requested, requested)
}

func ReplicaCount(replicas []string) int { return len(replicas) }
func ReplicaMissing(replicas []string, requested int) int {
	missing := requested - len(replicas)
	if missing < 0 {
		return 0
	}
	return missing
}

func ReplicaPolicy(requested, available int) map[string]int {
	return map[string]int{
		"requested": requested,
		"available": available,
		"effective": NormalizeReplicaCount(requested, available),
		"missing":   ReplicaMissing(make([]string, available), requested),
	}
}

func ReplicaAllowed(requested, available int) bool {
	return requested >= 0 && available >= 0 && NormalizeReplicaCount(requested, available) > 0
}

func ReplicaLimit(requested, available int) int {
	if requested <= 0 || available <= 0 {
		return 0
	}
	if requested < available {
		return requested
	}
	return available
}

func ReplicaSummary(requested, available int) map[string]bool {
	effective := ReplicaLimit(requested, available)
	return map[string]bool{
		"requested_valid":  requested >= 0,
		"available_valid":  available >= 0,
		"has_primary":      effective > 0,
		"fully_replicated": effective == requested,
	}
}

func ReplicaRequestValid(requested int) bool {
	return requested >= 0
}

func ReplicaOverflow(requested, available int) bool {
	if requested < 0 || available < 0 {
		return false
	}
	return requested > available
}

func ReplicaBalanced(requested, available int) bool {
	return !ReplicaOverflow(requested, available)
}
func ReplicaExact(requested, available int) bool { return requested == available }
func ReplicaDifference(requested, available int) int {
	return requested - available
}
