package sanitizer

import (
	"regexp"
	"strings"
)

// Precompile regex patterns for efficiency
var (
	// Matches one or more space characters
	spaceRegex = regexp.MustCompile(` {2,}`)

	// Matches one or more space characters followed by a dot or comma
	spaceBeforePunctRegex = regexp.MustCompile(` +(?:([.,]))`)
)

// SanitizeString removes unnecessary spaces from the input string.
// Unnecessary spaces include:
// 1. Leading and trailing spaces (tabs and newlines are preserved).
// 2. Multiple consecutive spaces replaced by a single space.
// 3. Spaces before dots and commas.
func SanitizeString(input string) string {
	// Step 1: Trim leading and trailing spaces only
	trimmed := strings.TrimLeft(input, " ")
	trimmed = strings.TrimRight(trimmed, " ")

	// Step 2: Replace multiple spaces with a single space
	singleSpaced := spaceRegex.ReplaceAllString(trimmed, " ")

	// Step 3: Remove spaces before dots and commas
	cleaned := spaceBeforePunctRegex.ReplaceAllString(singleSpaced, "$1")

	return cleaned
}
