package helpers

import "strings"

// InferVariableName infers the variable name from the given word.
// If the word starts with multiple uppercase letters, it retains the capitalization.
// Otherwise, it decapitalizes only the first letter.
func InferVariableName(word string) string {
	if len(word) == 0 {
		return ""
	}

	// Check if the first two characters are uppercase letters
	if len(word) >= 2 && isUpper(word[0]) && isUpper(word[1]) {
		return word
	}

	// Decapitalize only the first letter
	return strings.ToLower(string(word[0])) + word[1:]
}

// isUpper checks if a byte represents an uppercase letter.
func isUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}
