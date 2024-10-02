package matcher_v2

import (
	"mig/pkg/migrator/matcher_v2/mutators"
	"mig/pkg/migrator/matcher_v2/parser"
)

// Define the handler map using the HandlerFunc type
var handlerMap = map[string]mutators.HandlerFunc{
	"Wrap":   mutators.HandleWrap,
	"Wrapf":  mutators.HandleWrapf,
	"Errorf": mutators.HandleErrorf,
	"New":    mutators.HandleNew,
}

// HandleLine receives a line of code and returns a transformed line.
// If `errors` package invocation found, it will either be replaced with `errkit` package invocation
// or marked as to be migrated manually.
func HandleLine(line string) string {
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
	mutatedErrorsPart := mutators.Mutator(errorsPart, handlerMap)

	if mutatedErrorsPart == errorsPart {
		// No modification made by Mutator, mark this line as to be migrated manually
		return line + " // TODO: migrate manually"
	}

	// Return the concatenated prefix, mutatedErrorsPart, and suffix
	return prefix + mutatedErrorsPart + suffix
}
