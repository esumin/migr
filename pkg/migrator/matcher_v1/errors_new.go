package matcher_v1

import "strings"

func MatchSimpleErrorsNew(line string) string {
	if !strings.Contains(line, "errors.New") {
		return ""
	}

	replaced := strings.Replace(line, "errors.New", "errkit.New", -1)

	if strings.Contains(line, "fmt.Sprintf") {
		replaced = replaced + " // TODO: Fixme"
	}

	return replaced
}
