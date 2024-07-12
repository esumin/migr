package util

import (
	"regexp"
	"strings"
	"unicode"
)

func decapitalizeFirst(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

func decapitalizeAll(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	all := true
	for i, r := range runes {
		if !unicode.IsUpper(r) {
			all = false
		}
		runes[i] = unicode.ToLower(r)
	}

	if all {
		return string(runes)
	}

	return s
}

func removeBrackets(s string) string {
	if len(s) < 2 {
		return s
	}

	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return s[1 : len(s)-1]
	}

	return s
}

func tryNameWithID(s string) string {
	regex1 := regexp.MustCompile(`(.*)\s+(.*)\s+with\s+ID`)
	match := regex1.FindStringSubmatch(s)
	if match == nil {
		regex2 := regexp.MustCompile(`(.*)\s+(.*)\s+ID`)
		match = regex2.FindStringSubmatch(s)
	}

	if match == nil {
		return ""
	}

	return match[2] + "ID"
}

func tryLastWordOrParamName(s string, s1 string) string {
	spl := strings.Split(s, " ")
	paramName := spl[len(spl)-1]

	if paramName == "type" {
		paramName = spl[len(spl)-2] + "Type"
	}

	if paramName == "create" {
		paramName = spl[len(spl)-2]
	}

	if paramName == "got" {
		paramName = s1
	}

	return paramName
}

func GuessParamName(s string, s1 string) string {
	paramName := tryNameWithID(s)
	if paramName == "" {
		paramName = tryLastWordOrParamName(s, s1)
	}

	return decapitalizeFirst(decapitalizeAll(removeBrackets(paramName)))
}
