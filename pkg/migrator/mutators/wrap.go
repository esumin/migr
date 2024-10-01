package mutators

import (
	"fmt"
	"strings"

	parammatcher "mig/pkg/migrator/mutators/matcher"
)

var (
	wrapMatchers = []func([]string) []string{
		parammatcher.MatchSingleVariableAppend,
	}
	wrapfMatchers = []func([]string) []string{
		parammatcher.MatchSingleVariableAppend,
		parammatcher.MatchOneVariableSimple,
		parammatcher.MatchTwoVariablesSimple,
		parammatcher.MatchTwoVariablesNoName,
		parammatcher.MatchCurlyBracedTwoVariables,
	}
)

// HandleWrap takes a slice of arguments and applies matchers to format the elements.
// It returns the formatted error wrapping string.
var HandleWrap = getHandler(wrapMatchers)
var HandleWrapf = getHandler(wrapfMatchers)

func getHandler(matchers []func([]string) []string) func([]string) string {
	return func(args []string) string {
		// If no arguments are provided, return an empty string
		if len(args) == 0 {
			return ""
		}

		// If only err arg is provided, return an empty string
		if len(args) == 1 {
			return ""
		}

		// Extract the first argument, which is the error variable (e.g., "err")
		errVar := args[0]
		// Extract the remaining arguments, which are the error message and variables
		params := args[1:]

		// Try matching the argument using available matchers
		for _, matcher := range matchers {
			// Each matcher expects a slice as input, so we wrap the current arg in a slice
			matchedResult := matcher(params)
			if matchedResult != nil {
				return fmt.Sprintf("errkit.Wrap(%s, %s)", errVar, strings.Join(matchedResult, ", "))
			}
		}

		// If two arguments are provided, return the formatted error wrapping string
		if len(args) == 2 {
			return fmt.Sprintf("errkit.Wrap(%s, %s)", args[0], args[1])
		}

		return ""
	}
}
