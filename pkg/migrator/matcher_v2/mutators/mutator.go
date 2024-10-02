package mutators

import (
	"fmt"
	"strings"
)

// HandlerFunc is the handler function type
type HandlerFunc func(args []string) string

// HandlerMap is a map from function names to their corresponding handler functions.
type HandlerMap map[string]HandlerFunc

// Mutator takes the errorsPart and a handlerMap, parses it to extract the function name and arguments,
// and dispatches it to the appropriate handler function using the handlerMap.
func Mutator(errorsPart string, handlerMap HandlerMap) string {
	// Step 1: Strip 'errors.' prefix
	if !strings.HasPrefix(errorsPart, "errors.") {
		// Not an errors call, return as is
		return errorsPart
	}
	rest := errorsPart[len("errors."):] // Remove 'errors.' prefix

	// Step 2: Extract function name and arguments
	funcName, args, err := parseFunctionCall(rest)
	if err != nil {
		// Parsing failed, return original
		return errorsPart
	}

	// Step 3: Dispatch to the appropriate handler using the handlerMap
	if handler, exists := handlerMap[funcName]; exists {
		transformed := handler(args)
		if transformed != "" {
			return transformed
		}
		// If the handler returns an empty string, return the original errorsPart
		return errorsPart
	}

	// Unknown function, return original
	return errorsPart
}

// parseFunctionCall parses a function call string to extract the function name and arguments.
// For example, given "Wrap(err, \"message\")", it returns "Wrap" and ["err", "\"message\""].
func parseFunctionCall(callStr string) (funcName string, args []string, err error) {
	// Find the function name by locating the first '('
	idx := strings.Index(callStr, "(")
	if idx == -1 {
		return "", nil, fmt.Errorf("no opening parenthesis found in function call")
	}
	funcName = strings.TrimSpace(callStr[:idx])
	argsStr := callStr[idx+1:]

	// Find the matching closing parenthesis
	depth := 1
	endIdx := idx + 1
	for i := idx + 1; i < len(callStr); i++ {
		c := callStr[i]
		if c == '(' {
			depth++
		} else if c == ')' {
			depth--
			if depth == 0 {
				endIdx = i
				break
			}
		} else if c == '"' || c == '\'' || c == '`' {
			// Skip over string literals
			quote := c
			i++
			for i < len(callStr) {
				if callStr[i] == '\\' {
					i += 2 // Skip escaped character
				} else if callStr[i] == quote {
					break
				} else {
					i++
				}
			}
		}
	}
	if depth != 0 {
		return "", nil, fmt.Errorf("unbalanced parentheses in function call")
	}

	// Extract arguments string
	argsStr = callStr[idx+1 : endIdx]
	args, err = splitArguments(argsStr)
	if err != nil {
		return "", nil, err
	}
	return funcName, args, nil
}

// splitArguments splits the arguments string into individual arguments, handling nested parentheses and string literals.
func splitArguments(argsStr string) ([]string, error) {
	args := []string{}
	start := 0
	depth := 0
	inString := false
	var stringChar rune

	for i := 0; i < len(argsStr); i++ {
		c := argsStr[i]
		if inString {
			if c == '\\' {
				i++ // Skip escaped character
			} else if c == byte(stringChar) {
				inString = false
			}
		} else {
			if c == '(' || c == '[' || c == '{' {
				depth++
			} else if c == ')' || c == ']' || c == '}' {
				depth--
			} else if c == '"' || c == '\'' || c == '`' {
				inString = true
				stringChar = rune(c)
			} else if c == ',' && depth == 0 {
				arg := strings.TrimSpace(argsStr[start:i])
				args = append(args, arg)
				start = i + 1
			}
		}
	}
	if inString {
		return nil, fmt.Errorf("unclosed string literal in arguments")
	}
	arg := strings.TrimSpace(argsStr[start:])
	if arg != "" {
		args = append(args, arg)
	}
	return args, nil
}

// Handler functions with signatures only

// handleWrapf handles the errors.Wrapf function.
// For example, errors.Wrapf(err, "format", args...) => errkit.Wrap(err, "message", "param", value)
func handleWrapf(args []string) string {
	// Implementation will be added later
	return ""
}

// handleNew handles the errors.New function.
// For example, errors.New("message") => errkit.New("message")
func handleNew(args []string) string {
	// Implementation will be added later
	return ""
}

// handleErrorf handles the errors.Errorf function.
// For example, errors.Errorf("format", args...) => errkit.New(fmt.Sprintf("format", args...))
func handleErrorf(args []string) string {
	// Implementation will be added later
	return ""
}
