package matcher_v1

import (
	"fmt"
	"regexp"
	"strings"
)

func matchWrapf(line string) string {
	simpleWrap := regexp.MustCompile(`(.*)errors.Wrapf\((.*)\)(.*)`)
	match := simpleWrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[2], "%") {
		return ""
	}

	replaced := fmt.Sprintf("%serrkit.Wrap(%s)%s", match[1], match[2], match[3])
	if strings.Contains(match[2], "fmt.Sprintf") {
		replaced = replaced + " // TODO: Fixme"
	}

	return replaced
}

func matchWrap(line string) string {
	simpleWrap := regexp.MustCompile(`(.*)errors.Wrap\((.*)\)(.*)`)

	match := simpleWrap.FindStringSubmatch(line)
	replaced := ""
	if match != nil {
		replaced = fmt.Sprintf("%serrkit.Wrap(%s)%s", match[1], match[2], match[3])
		if strings.Contains(match[2], "fmt.Sprintf") {
			replaced = replaced + " // TODO: Fixme"
		}
	}

	return replaced
}

func MatchSimpleWraps(line string) string {
	replaced := matchWrap(line)
	if replaced != "" {
		return replaced
	}

	return matchWrapf(line)
}
