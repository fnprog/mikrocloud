package containers

import (
	"strings"
)

// SanitizeDockerName sanitizes a string to be used as a Docker container or image name.
// Docker names must match the pattern: [a-zA-Z0-9][a-zA-Z0-9_.-]*
// This function:
// - Converts to lowercase
// - Replaces spaces with hyphens
// - Replaces any invalid characters with hyphens
// - Ensures the result only contains: a-z, 0-9, hyphens, underscores, and periods
func SanitizeDockerName(name string) string {
	result := strings.ToLower(name)
	result = strings.ReplaceAll(result, " ", "-")
	result = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			return r
		}
		return '-'
	}, result)
	return result
}
