package mutators

import (
	"fmt"
	"strings"

	"mig/pkg/migrator/matcher_v2/mutators/matcher"
	"mig/pkg/migrator/matcher_v2/mutators/sanitizer"
)

type MatcherFn func([]string) []string

var (
	wrapMatchers = []MatcherFn{
		param_matcher.MatchSingleVariableAppend,
	}
	wrapfMatchers = []MatcherFn{
		param_matcher.MatchSingleVariableAppend,
		param_matcher.MatchOneVariableSimple,
		param_matcher.MatchTwoVariablesSimple,
		param_matcher.MatchTwoVariablesNoName,
		param_matcher.MatchCurlyBracedTwoVariables,
	}
	errorfMatchers = []MatcherFn{
		param_matcher.MatchSingleVariableAppend,
		param_matcher.MatchOneVariableSimple,
		param_matcher.MatchTwoVariablesSimple,
		param_matcher.MatchTwoVariablesNoName,
		param_matcher.MatchCurlyBracedTwoVariables,
	}
)

// HandleWrap takes a slice of arguments and applies matchers to format the elements.
// It returns the formatted error wrapping string.
var HandleWrap = getWrapHandler(wrapMatchers)
var HandleWrapf = getWrapHandler(wrapfMatchers)
var HandleErrorf = getNewHandler(errorfMatchers)
var HandleNew = getNewHandler([]MatcherFn{})

func getNewHandler(matchers []MatcherFn) func([]string) string {
	return func(args []string) string {
		// If no arguments are provided, return an empty string
		if len(args) == 0 {
			return ""
		}

		// If single argument is provided, return the formatted error wrapping string
		if len(args) == 1 {
			return fmt.Sprintf("errkit.New(%s)", sanitizer.SanitizeString(args[0]))
		}

		return fmt.Sprintf(`errkit.New(fmt.Sprintf(%s, %s))`, args[0], sanitizer.SanitizeString(strings.Join(args[1:], ", ")))
	}
}

func getWrapHandler(matchers []MatcherFn) func([]string) string {
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
				return fmt.Sprintf("errkit.Wrap(%s, %s)", errVar, sanitizer.SanitizeString(strings.Join(matchedResult, ", ")))
			}
		}

		// If two arguments are provided, return the formatted error wrapping string
		if len(args) == 2 {
			return fmt.Sprintf("errkit.Wrap(%s, %s)", args[0], sanitizer.SanitizeString(args[1]))
		}

		return ""
	}
}
