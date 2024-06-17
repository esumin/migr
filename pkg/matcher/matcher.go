package matcher

type Matcher func(line string) string

// Matchers should go from more complex to less complex
var AllMatchers = []Matcher{
	MatchErrorfWithNamedParams,
	MatchWrapfWithNamedParams,
	MatchWrapfStderr,
	MatchSimpleWraps,
	MatchSimpleErrorsNew,
	MatchImport,
}

func MatchSequentially(matchers []Matcher, line string) string {
	for _, m := range matchers {
		r := m(line)
		if r != "" {
			return r
		}
	}

	return ""
}

func maybeTrimComma(msg string) string {
	l := len(msg) - 1
	if msg[l:] == "," || msg[l:] == ":" {
		return msg[:l]
	}

	return msg
}
