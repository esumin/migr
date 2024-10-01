package mutators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// HandleErrorParams processes the error template and its arguments to extract
// static messages and parameter names. It returns a transformed slice of strings
// suitable for constructing the new error handling call.
func HandleErrorParams(args []string) []string {
	// Special case: no arguments or one argument, return as is.
	if len(args) <= 1 {
		return args
	}

	firstArg := args[0]

	if !isQuoted(firstArg) {
		// Special case: first argument is not a quoted string, return as is.
		return args
	}

	// Unquote the format string
	formatStr, err := strconv.Unquote(firstArg)
	if err != nil {
		// If unquoting fails, return args as is.
		return args
	}

	// Find all format specifiers (e.g., %s, %d, %v)
	specifierPositions := findFormatSpecifiers(formatStr)
	if len(specifierPositions) == 0 {
		// No format specifiers found, return the original string as the static part.
		return []string{firstArg}
	}

	// Extract static parts and infer parameter names based on specifiers
	staticParts, inferredParamNames, err := extractStaticAndParamNames(formatStr, specifierPositions)
	if err != nil {
		// If extraction fails, return args as is.
		return args
	}

	varNames := args[1:]
	if len(varNames) != len(inferredParamNames) {
		// Mismatch between number of variables and specifiers, return args as is.
		return args
	}

	// Rebuild the static string without trailing punctuation
	staticString := rebuildStaticParts(staticParts)

	// Initialize the result with the static string
	result := []string{staticString}

	// Append parameter name and variable pairs
	for i, paramName := range inferredParamNames {
		if paramName == "" {
			// Derive parameter name from the variable name by taking the last identifier after '.'
			parts := strings.Split(varNames[i], ".")
			paramName = parts[len(parts)-1]
		}
		result = append(result, `"`+paramName+`"`, varNames[i])
	}

	return result
}

// isQuoted checks if a string is enclosed in single quotes, double quotes, or backticks.
func isQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}
	return (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') ||
		(s[0] == '`' && s[len(s)-1] == '`')
}

// findFormatSpecifiers finds all positions of format specifiers (e.g., %s, %d) in the format string.
// It returns a slice of starting indices for each specifier.
func findFormatSpecifiers(s string) []int {
	// Regex to match format specifiers like %s, %d, %v, etc.
	re := regexp.MustCompile(`%[sdv]`)
	matches := re.FindAllStringIndex(s, -1)
	specifiers := []int{}
	for _, m := range matches {
		specifiers = append(specifiers, m[0])
	}
	return specifiers
}

// extractStaticAndParamNames extracts the static parts of the format string and infers parameter names.
// It returns slices of static strings and parameter names corresponding to each format specifier.
func extractStaticAndParamNames(formatStr string, specifiers []int) ([]string, []string, error) {
	staticParts := []string{}
	paramNames := []string{}
	lastIndex := 0

	for _, specPos := range specifiers {
		// Extract static text before the specifier
		staticText := formatStr[lastIndex:specPos]
		staticParts = append(staticParts, staticText)

		// Infer parameter name from the static text
		paramName, err := guessParamName(staticText)
		if err != nil {
			// If unable to guess from static text, leave paramName empty to use variable name later
			paramName = ""
		}
		paramNames = append(paramNames, paramName)

		// Update lastIndex to the character after the specifier
		lastIndex = specPos + 2 // Assuming specifier is two characters like %s
	}

	// Add any remaining static text after the last specifier
	staticText := formatStr[lastIndex:]
	staticParts = append(staticParts, staticText)

	return staticParts, paramNames, nil
}

// guessParamName infers a parameter name based on the static text preceding a format specifier.
// It extracts the last meaningful word before the specifier and converts it to camelCase.
// If no suitable word is found, it returns an empty string.
func guessParamName(staticText string) (string, error) {
	staticText = strings.TrimRight(staticText, " ")
	if len(staticText) == 0 {
		return "", fmt.Errorf("no text before format specifier")
	}

	// Regex to capture the last word before the specifier
	// It looks for word characters possibly followed by punctuation
	re := regexp.MustCompile(`(\w+)[\s,:/-]*$`)
	matches := re.FindStringSubmatch(staticText)
	if len(matches) < 2 {
		return "", fmt.Errorf("no suitable word found before format specifier")
	}

	lastWord := matches[1]

	// Convert the last word to camelCase
	paramName := toCamelCase(lastWord)
	return paramName, nil
}

// toCamelCase converts the first character of a string to lowercase.
// For example, "Namespace" becomes "namespace".
func toCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

// rebuildStaticParts reconstructs the static parts into a single quoted string.
// It joins all static segments, removes any trailing punctuation, and wraps them in double quotes.
func rebuildStaticParts(staticParts []string) string {
	joined := strings.Join(staticParts, "")
	// Remove trailing punctuation like ":" or "/"
	joined = strings.TrimRight(joined, ":/")
	return `"` + joined + `"`
}
