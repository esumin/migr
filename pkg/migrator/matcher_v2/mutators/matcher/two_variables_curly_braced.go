package param_matcher

import (
	"regexp"
)

// MatchCurlyBracedTwoVariables takes a slice of strings and matches the first element
// against a template containing two variables within curly braces. If matched,
// it returns a new slice with the extracted template and parameters.
// Otherwise, it returns nil.
func MatchCurlyBracedTwoVariables(input []string) []string {
	if len(input) != 3 {
		return nil
	}

	template := input[0]
	param1 := input[1]
	param2 := input[2]

	// Define a regex pattern to match the template.
	// The pattern ensures:
	// 1. The string starts and ends with a double quote.
	// 2. Contains exactly two placeholders within curly braces.
	pattern := `^"([^{}]*)\{(\w+):?\s+%s,\s+(\w+):?\s+%s\}"$`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(template)
	if match == nil {
		return nil
	}

	// Extract the part before the curly braces.
	before := match[1]

	// Construct and return the resulting slice.
	return []string{
		`"` + before + `"`,
		`"` + param1 + `"`,
		param1,
		`"` + param2 + `"`,
		param2,
	}
}
