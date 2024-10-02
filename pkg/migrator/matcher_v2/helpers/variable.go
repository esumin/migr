// helpers/helpers.go
package helpers

import "strings"

// blacklistedWords are words that should not be used as variable names
var blacklistedWords = map[string]struct{}{
	"uninstall": {},
}

func isBlacklisted(word string) bool {
	_, isBlacklisted := blacklistedWords[word]
	return isBlacklisted
}

// requiresCombination checks if the word requires concatenation with the pre-last word.
func requiresCombination(word string) bool {
	return word == "name"
}

func getVarNameFromLastWords(lastWords []string) string {
	if len(lastWords) == 0 {
		return ""
	}

	preLastWord := lastWords[0]
	lastWord := ""
	if len(lastWords) == 1 {
		lastWord = preLastWord
		preLastWord = ""
	} else {
		lastWord = lastWords[1]
	}

	if isBlacklisted(lastWord) {
		return ""
	}

	if !requiresCombination(lastWord) {
		return decapitalize(lastWord)
	}

	if preLastWord == "" {
		return "" // Cannot concatenate without pre-last word
	}

	return decapitalize(preLastWord) + capitalize(decapitalize(lastWord))
}

// InferVariableName infers the variable name based on the provided words and varExpr.
// If the last word is blacklisted, it infers the variable name from varExpr.
// If the last word (for instance "name") requires concatenation with pre-last word
// it concatenates the pre last word with the last one.
// If the word starts with multiple uppercase letters, it retains the capitalization.
// Otherwise, it decapitalizes only the first letter.
func InferVariableName(lastWords []string, varExpr string) string {
	varName := getVarNameFromLastWords(lastWords)
	if varName != "" {
		return varName
	}

	// We are unable to guess name from last words, let's try variable expression
	// Check if the word is blacklisted
	return inferFromVarExpr(varExpr)
}

// inferFromVarExpr extracts the last word from varExpr and infers the variable name from it.
func inferFromVarExpr(varExpr string) string {
	// Split varExpr by dots
	parts := strings.Split(varExpr, ".")
	if len(parts) == 0 {
		return ""
	}

	// Take the last part
	lastExprWord := parts[len(parts)-1]
	lastExprWord = strings.Trim(lastExprWord, ":/")

	if len(lastExprWord) == 0 {
		return ""
	}

	return decapitalize(lastExprWord)
}

// isUpper checks if a byte represents an uppercase letter.
func isUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func decapitalize(word string) string {
	if len(word) == 0 {
		return ""
	}

	// Do not decapitalize if more than one first letter are uppercase
	if len(word) >= 2 && isUpper(word[0]) && isUpper(word[1]) {
		return word
	}

	return strings.ToLower(string(word[0])) + word[1:]
}

func capitalize(word string) string {
	if len(word) == 0 {
		return ""
	}

	return strings.ToUpper(string(word[0])) + word[1:]
}

func GetLastWords(s string) []string {
	w := strings.Fields(s)
	if len(w) == 0 {
		return nil
	}

	if len(w) == 1 {
		return w
	}

	return w[len(w)-2:]
}
