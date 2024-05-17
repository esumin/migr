package matcher

type Matcher func(line string) string

// Matchers should go from more complex to less complex
var AllMatchers = []Matcher{
	MatchWrapfStderr,
	MatchSimpleWraps,
	MatchSimpleErrorsNew,
	MatchImport,
}
