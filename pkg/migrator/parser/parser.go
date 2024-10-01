package parser

import (
	"fmt"
	"strings"
	"unicode"
)

// ParseLine locates the 'errors' function call in the line and splits it into prefix, errorsPart, and suffix.
// Returns an error if any unexpected parsing error occurs.
// If 'errors' not found, returns the entire line as prefix, and errorsPart and suffix as empty strings.
func ParseLine(line string) (prefix string, errorsPart string, suffix string, err error) {
	// Find the 'errors' invocation in the line
	start, end, err := FindErrorsInvocation(line)
	if err != nil {
		return "", "", "", err
	}
	if start == -1 || end == -1 {
		// No 'errors' invocation found, return the line as prefix
		return line, "", "", nil
	}

	// Validate indices
	if start < 0 || end > len(line) || start >= end {
		return "", "", "", fmt.Errorf("invalid indices for splitting the line")
	}

	// Split the line into prefix, errorsPart, and suffix
	prefix = line[:start]
	errorsPart = line[start:end]
	suffix = line[end:]

	return prefix, errorsPart, suffix, nil
}

// FindErrorsInvocation locates the 'errors' function call in the line.
// Returns the start and end indices of the invocation and an error if parsing fails.
func FindErrorsInvocation(line string) (start int, end int, err error) {
	// Find the first occurrence of 'errors.' in the line
	idx := strings.Index(line, "errors.")
	if idx == -1 {
		return -1, -1, nil // 'errors.' not found
	}

	// Move past 'errors.' to get to the function name
	funcNameStart := idx + len("errors.")
	funcNameEnd := funcNameStart

	// Find the end of the function name
	for funcNameEnd < len(line) && (unicode.IsLetter(rune(line[funcNameEnd])) || unicode.IsDigit(rune(line[funcNameEnd]))) {
		funcNameEnd++
	}

	// Find the '(' after the function name
	remainingLine := line[funcNameEnd:]
	parenIdx := strings.Index(remainingLine, "(")
	if parenIdx == -1 {
		return -1, -1, nil // No '(' after function name
	}
	parenIdx += funcNameEnd // Adjust index to absolute

	// Now, find the matching closing parenthesis ')'
	depth := 0
	i := parenIdx
	for i < len(line) {
		c := line[i]
		if c == '(' {
			depth++
		} else if c == ')' {
			depth--
			if depth == 0 {
				// Found matching ')'
				return idx, i + 1, nil // end index is exclusive
			}
		} else if c == '"' || c == '\'' || c == '`' {
			// Skip over string literals
			quote := c
			i++
			escaped := false
			for i < len(line) {
				if line[i] == '\\' && !escaped {
					escaped = true
					i++
					continue
				}
				if line[i] == quote && !escaped {
					break
				}
				escaped = false
				i++
			}
			if i >= len(line) {
				return -1, -1, fmt.Errorf("unclosed string literal starting at position %d", i)
			}
		}
		i++
	}

	if depth != 0 {
		return -1, -1, fmt.Errorf("unmatched parentheses in line")
	}

	// No matching ')' found
	return -1, -1, fmt.Errorf("could not find matching closing parenthesis")
}
