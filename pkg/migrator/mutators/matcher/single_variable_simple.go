// param_matcher/one_variable_simple.go
package param_matcher

import (
	"regexp"
	"strings"

	"mig/pkg/helpers"
)

// MatchOneVariableSimple takes a slice of strings and matches the first element
// against a template containing exactly one unnamed variable without curly braces.
// If matched, it returns a new slice with the extracted static template,
// the inferred variable name, and the variable expression.
// Otherwise, it returns nil.
//
// The variable name is inferred from the last word of the static part before the placeholder.
// If the last word starts with multiple uppercase letters (e.g., "PVC"), it retains the capitalization.
// Otherwise, it decapitalizes the first letter.
func MatchOneVariableSimple(input []string) []string {
	if len(input) != 2 {
		return nil
	}

	template := input[0]
	varExpr := input[1]

	// Define a regex pattern to match templates with exactly one placeholder (%s or %v)
	// Pattern breakdown:
	// ^"(.+?)\s*%[sv]\s*(.*)"$
	// ^"            : Start with a double quote
	// (.+?)         : Non-greedy capture of any characters (static part before placeholder)
	// \s*%[sv]      : Optional whitespace followed by %s or %v
	// \s*(.*)"$     : Optional whitespace followed by the remaining static part and end with a double quote
	pattern := `^"(.+?)\s*%[sv]\s*(.*)"$`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(template)
	if match == nil {
		return nil
	}

	// Extract static parts before and after the placeholder
	staticBefore := match[1]
	staticAfter := match[2]

	// Trim any trailing spaces from staticBefore and leading/trailing spaces from staticAfter
	staticBefore = strings.TrimSpace(staticBefore)
	staticAfter = strings.TrimSpace(staticAfter)

	// Combine the static parts, ensuring proper spacing
	if staticAfter != "" {
		staticBefore = staticBefore + " " + staticAfter
	}

	// Infer variable name from the last word before the placeholder
	// Split the staticBefore into words and take the last one
	words := strings.Fields(staticBefore)
	if len(words) == 0 {
		return nil
	}
	lastWord := words[len(words)-1]

	// Remove any trailing punctuation from the last word
	lastWord = strings.Trim(lastWord, ":/")

	// Infer the variable name
	varName := helpers.InverVariableName(lastWord)

	// If variable name is empty after inference, return nil
	if varName == "" {
		return nil
	}

	// Construct and return the resulting slice
	return []string{
		`"` + staticBefore + `"`,
		`"` + varName + `"`,
		varExpr,
	}
}
