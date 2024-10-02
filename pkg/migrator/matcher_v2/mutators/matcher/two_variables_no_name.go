// param_matcher/two_variables_no_name.go
package param_matcher

import (
	"regexp"
	"strings"
)

// MatchTwoVariablesNoName takes a slice of strings and matches the first element
// against a template containing two unnamed variables without curly braces.
// If matched, it returns a new slice with the extracted static template and parameters.
// Otherwise, it returns nil.
//
// The function infers variable names from the variable expressions by
// extracting the last part after any dot and converting the first letter to lowercase.
//
// Example:
// Input: []string{`"Failed to validate if PVC %s:%s exists"`, `namespace`, `claimName`}
// Output: []string{`"Failed to validate if PVC exists"`, `"namespace"`, `namespace`, `"claimName"`, `claimName`}
func MatchTwoVariablesNoName(input []string) []string {
	if len(input) != 3 {
		return nil
	}

	template := input[0]
	param1 := input[1]
	param2 := input[2]

	// Updated regex pattern:
	// ^"(.+?)\s*%s[/:]\s*%s\s*(.+)"$
	// Explanation:
	// ^"            : Start with a double quote
	// (.+?)         : Non-greedy capture of any characters (static part before first %s)
	// \s*%s         : Optional whitespace followed by first %s
	// [/:]          : Mandatory separator (either : or /)
	// \s*%s         : Optional whitespace followed by second %s
	// \s*(.+)"$     : Optional whitespace followed by the remaining static part and end with a double quote
	pattern := `^"([^/:]+).*\s*%s\s*[/:]\s*%s\s*(.*)"$`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(template)
	if match == nil {
		return nil
	}

	// Extract the static parts before and after the placeholders
	staticBefore := match[1]
	staticAfter := match[2]

	// Combine the static parts, ensuring proper spacing
	// Trim any trailing spaces from staticBefore and leading spaces from staticAfter
	staticBefore = strings.TrimSpace(staticBefore)
	staticAfter = strings.TrimSpace(staticAfter)

	// If there's text after the placeholders, include it with a space
	if staticAfter != "" {
		staticBefore = staticBefore + " " + staticAfter
	}

	// Function to extract variable name from expression
	// Function to extract variable name from expression
	getVarName := func(varExpr string) string {
		parts := strings.Split(varExpr, ".")
		lastPart := parts[len(parts)-1]
		if len(lastPart) == 0 {
			return ""
		}

		// isUpper checks if a byte represents an uppercase letter.
		isUpper := func(b byte) bool {
			return b >= 'A' && b <= 'Z'
		}
		// Check if the first two letters are uppercase (e.g., "PVCName")
		if len(lastPart) >= 2 && isUpper(lastPart[0]) && isUpper(lastPart[1]) {
			return lastPart
		}

		// Decapitalize only the first letter
		return strings.ToLower(string(lastPart[0])) + lastPart[1:]
	}

	varName1 := getVarName(param1)
	varName2 := getVarName(param2)

	// Construct and return the resulting slice
	return []string{
		`"` + staticBefore + `"`,
		`"` + varName1 + `"`,
		param1,
		`"` + varName2 + `"`,
		param2,
	}
}
