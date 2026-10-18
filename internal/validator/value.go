package validator

import "strings"

func isBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}
