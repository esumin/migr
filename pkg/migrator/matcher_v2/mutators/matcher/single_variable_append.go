package param_matcher

import (
	"regexp"
	"strings"

	"mig/pkg/migrator/matcher_v2/helpers"
)

// MatchSingleVariableAppend matches templates where a single variable is appended to the message
// using the + operator and transforms them into the correct format where the variable is passed separately.
// E.g. "failed to read env from dir:" + dirName => "failed to read env from dir", "dir", dirName
func MatchSingleVariableAppend(input []string) []string {
	// Ensure exactly one input string
	if len(input) != 1 {
		return nil
	}

	str := input[0]

	// Ensure exactly one '+' operator to prevent matching multiple variable appends
	if strings.Count(str, "+") != 1 {
		return nil
	}

	// Define a regex to match a single variable appended to the string with the + operator
	// Pattern breakdown:
	// ^"(.+?)\s*:\s*"\s*\+\s*(\w+)$
	// ^"        : Start with a double quote
	// (.+?)     : Non-greedy match to capture the static part of the string before ":" (message)
	// \s*:\s*   : Match the colon ":" with optional spaces around it
	// "\s*\+    : Match the + operator with optional spaces
	// (\w+)     : Match the variable name (alphanumeric word)
	// $         : End of the string (ensures that no additional variables or content follows)
	pattern := `^"(.+?)\s*:\s*"\s*\+\s*(\w+)$`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(str)

	// If no match is found, return nil
	if match == nil {
		return nil
	}

	// Extract the static message and the variable name
	message := match[1]
	varName := match[2]

	// Trim any trailing spaces from the message
	message = strings.TrimSpace(message)

	// Infer the variable label from the last word of the message
	words := strings.Fields(message)
	if len(words) == 0 {
		return nil
	}
	lastWord := words[len(words)-1]
	varLabel := helpers.InferVariableName(lastWord)
	if varLabel == "" {
		return nil
	}

	// Construct and return the resulting slice
	// Example: ["failed to read env from dir", "dir", dirName]
	return []string{
		`"` + message + `"`,
		`"` + varLabel + `"`,
		varName,
	}
}
