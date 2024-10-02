package param_matcher

import (
	"regexp"

	"mig/pkg/migrator/matcher_v2/helpers"
)

// MatchTwoVariablesSimple takes a slice of strings and matches the first element
// against a template containing two variables without curly braces.
// If matched, it returns a new slice with the extracted template and parameters.
// Otherwise, it returns nil.
func MatchTwoVariablesSimple(input []string) []string {
	if len(input) != 3 {
		return nil
	}

	template := input[0]
	param1 := input[1]
	param2 := input[2]

	// Updated regex pattern:
	// - Make the first capturing group non-greedy by using *?
	// Pattern breakdown:
	// ^"([^"]*?)\s*(\w+):?\s+%s,\s+(\w+):?\s+%s"$
	// ^"            : Start with a double quote
	// ([^"]*?)      : Non-greedy capture of any characters except double quotes
	// \s*           : Optional whitespace
	// (\w+)         : Capture first variable name
	// :?            : Optional colon
	// \s+%s         : At least one whitespace followed by %s
	// ,\s+          : Comma followed by at least one whitespace
	// (\w+)         : Capture second variable name
	// :?            : Optional colon
	// \s+%s         : At least one whitespace followed by %s
	// "$            : End with a double quote
	pattern := `^"([^"]*?)\s*(\w+):?\s+%s,\s+(\w+):?\s+%s"$`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(template)
	if match == nil {
		return nil
	}

	// Extract the static part before the variables.
	before := match[1]

	// Extract variable names.
	lastWords1 := []string{match[2]}
	lastWords2 := []string{match[3]}
	varName1 := helpers.InferVariableName(lastWords1, param1)
	varName2 := helpers.InferVariableName(lastWords2, param2)

	// Construct and return the resulting slice.
	return []string{
		`"` + before + `"`,
		`"` + varName1 + `"`,
		param1,
		`"` + varName2 + `"`,
		param2,
	}
}
