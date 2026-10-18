package validator

import "strings"

type idTracker struct {
	seen map[string]struct{}
}

func (tracker *idTracker) isDuplicate(id string) bool {
	id = strings.TrimSpace(id)

	if id == "" {
		return false
	}

	if _, exists := tracker.seen[id]; exists {
		return true
	}

	if tracker.seen == nil {
		tracker.seen = make(map[string]struct{})
	}

	tracker.seen[id] = struct{}{}
	return false
}
