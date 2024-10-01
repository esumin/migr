package migrator

import (
	"mig/pkg/migrator/parser"
)

// ReplaceErrors receives a line of code and returns a transformed line.
// It replaces usages of the `errors` module with the `errkit` module.
func ReplaceErrors(line string) string {
	// Use parser.ParseLine to find and split the line
	prefix, errorsPart, suffix, err := parser.ParseLine(line)
	if err != nil {
		// Handle unexpected error, return the original line
		return line
	}
	if errorsPart == "" {
		// No 'errors' invocation found, return the line as is
		return line
	}

	// Pass errorsPart to Mutator
	mutatedErrorsPart := Mutator(errorsPart, HandlerMap{})

	if mutatedErrorsPart == errorsPart {
		// No modification made by Mutator, return original line
		return line
	}

	// Return the concatenated prefix, mutatedErrorsPart, and suffix
	return prefix + mutatedErrorsPart + suffix
}
