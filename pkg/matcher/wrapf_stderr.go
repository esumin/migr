package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

// example
// return errors.Wrapf(err, "Error %s, resetting the mongodb application. stdout is %s", stderr, stdout)
func matchWrapfStderrStdout(line string) string {
	wrap := regexp.MustCompile(`(.*)errors.Wrapf\((.*?)\"Error\s+\%s(?:[:,]?)\s*(.*?)\s*stdout is \%s\",\s*stderr,\s*stdout\)(.*)`)
	match := wrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[3], "%") || strings.Contains(match[3], "fmt.Sprintf") {
		return ""
	}

	return fmt.Sprintf("%serrkit.Wrap(%s\"Error %s\", \"stdout\", stdout, \"stderr\", stderr)%s", match[1], match[2], match[3], match[4])
}

// example
// return errors.Wrapf(err, "Error while Pinging the database: %s, app: %s", stderr, a.name)
// return errors.Wrapf(err, "Failed to delete documents from default bucket. %s app=%s", stderr, cb.name)
func matchWrapfStderr3(line string) string {
	simpleWrap := regexp.MustCompile(`(.*)errors.Wrapf\((.*?)[:.]\s*%s,?\s*app[:=]\s*%s\",\s*stderr\,\s*(.*)\)(.*)`)
	match := simpleWrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[2], "%") || strings.Contains(match[2], "fmt.Sprintf") {
		return ""
	}

	replaced := fmt.Sprintf("%serrkit.Wrap(%s\", \"stderr\", stderr, \"app\", %s)%s", match[1], match[2], match[3], match[4])
	return replaced

}

// example
// return errors.Wrapf(err, "Error %s: Resetting the application.", stderr)
// return errors.Wrapf(err, "Error %s while pinging the database.", stderr)
// return errors.Wrapf(err, "Error %s, deleting resources while reseting application.", stderr)

func matchWrapfStderr2(line string) string {
	simpleWrap := regexp.MustCompile(`(.*)errors.Wrapf\((.*?)\"Error\s+\%s(?:[:,]?)(.*?)\",\s*stderr\)(.*)`)
	match := simpleWrap.FindStringSubmatch(line)
	if match == nil ||
		strings.Contains(match[2], "%") ||
		strings.Contains(match[2], "fmt.Sprintf") ||
		strings.Contains(match[3], "%") ||
		strings.Contains(match[3], "fmt.Sprintf") {
		return ""
	}

	replaced := fmt.Sprintf("%serrkit.Wrap(%s\"Error%s\", \"stderr\", stderr)%s", match[1], match[2], match[3], match[4])
	return replaced
}

// examples
//
// return errors.Wrapf(err, "Failed to ping postgresql DB. %s", stderr)
// return errors.Wrapf(err, "Failed to ping the application. Error:%s", stderr)
func matchWrapfStderr1(line string) string {
	simpleWrap := regexp.MustCompile(`(.*)errors.Wrapf\((.*?)(?:\.\s*Error:|Error:|:|\.)?\s*\%s\s*\"\,\s*stderr\)(.*)`)
	match := simpleWrap.FindStringSubmatch(line)
	if match == nil || strings.Contains(match[2], "%") || strings.Contains(match[2], "fmt.Sprintf") {
		return ""
	}

	replaced := fmt.Sprintf("%serrkit.Wrap(%s\", \"stderr\", stderr)%s", match[1], match[2], match[3])
	return replaced
}

func MatchWrapfStderr(line string) string {
	replaced := matchWrapfStderr1(line)
	if replaced != "" {
		return replaced
	}

	replaced = matchWrapfStderr2(line)
	if replaced != "" {
		return replaced
	}

	replaced = matchWrapfStderr3(line)
	if replaced != "" {
		return replaced
	}

	replaced = matchWrapfStderrStdout(line)
	if replaced != "" {
		return replaced
	}

	return ""
}
