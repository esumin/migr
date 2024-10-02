// migrator/mutator_test.go
package mutators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	matcher "mig/pkg/migrator/matcher_v2/mutators"
)

func TestMutator(t *testing.T) {
	type testCase struct {
		name                  string
		errorsPart            string
		expectedHandlerName   string
		expectedHandlerArgs   []string
		expectedHandlerOutput string
	}
	// Define test cases
	tests := []testCase{
		{
			name:                  "Simple Wrap",
			errorsPart:            `errors.Wrap(err, "Failed to get secrets")`,
			expectedHandlerName:   "Wrap",
			expectedHandlerArgs:   []string{"err", `"Failed to get secrets"`},
			expectedHandlerOutput: "aaa",
		},
		{
			name:                  "Simple Wrapf",
			errorsPart:            `errors.Wrapf(err, "Failed to process %s", data)`,
			expectedHandlerName:   "Wrapf",
			expectedHandlerArgs:   []string{"err", `"Failed to process %s"`, "data"},
			expectedHandlerOutput: "aaa",
		},
		{
			name:                  "Simple New",
			errorsPart:            `errors.New("An error occurred")`,
			expectedHandlerName:   "New",
			expectedHandlerArgs:   []string{"\"An error occurred\""},
			expectedHandlerOutput: "aaa",
		},
		{
			name:                  "Simple Errorf",
			errorsPart:            `errors.Errorf("Failed with code %d", code)`,
			expectedHandlerName:   "Errorf",
			expectedHandlerArgs:   []string{"\"Failed with code %d\"", "code"},
			expectedHandlerOutput: "aaa",
		},
		{
			name:                  "Unknown Function",
			errorsPart:            `errors.Unknown(err, "message")`,
			expectedHandlerOutput: `errors.Unknown(err, "message")`,
		},
		{
			name:                  "Errorf with Multiple Arguments",
			errorsPart:            `errors.Errorf("Error: %s, Code: %d", msg, code)`,
			expectedHandlerName:   "Errorf",
			expectedHandlerArgs:   []string{"\"Error: %s, Code: %d\"", "msg", "code"},
			expectedHandlerOutput: "aaa",
		},
		{
			name:                  "Errorf with No Arguments",
			errorsPart:            `errors.Errorf("Simple error message")`,
			expectedHandlerName:   "Errorf",
			expectedHandlerArgs:   []string{"\"Simple error message\""},
			expectedHandlerOutput: "aaa",
		},
		{
			name:                  "Wrapf with multiple named arguments 1",
			errorsPart:            `errors.Wrapf(err, "Failed to get pod from podOptions. Namespace: %s, NameFmt: %s", opts.Namespace, opts.GenerateName)`,
			expectedHandlerName:   "Wrapf",
			expectedHandlerArgs:   []string{"err", `"Failed to get pod from podOptions. Namespace: %s, NameFmt: %s"`, "opts.Namespace", "opts.GenerateName"},
			expectedHandlerOutput: "aaa",
		},
	}

	getHandler := func(invokedArgs map[string][][]string, handlerName string, output string) matcher.HandlerFunc {
		return func(args []string) string {
			invokedArgs[handlerName] = append(invokedArgs[handlerName], args)
			return output
		}
	}

	getHandlerMap := func(invokedArgs map[string][][]string, tc testCase) matcher.HandlerMap {
		return matcher.HandlerMap{
			tc.expectedHandlerName: getHandler(invokedArgs, tc.expectedHandlerName, tc.expectedHandlerOutput),
		}
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Flags to track handler invocations and store arguments
			invokedArgs := make(map[string][][]string)

			// Define the handlerMap with wrapped handlers that store arguments
			handlerMap := getHandlerMap(invokedArgs, tc)

			// Call Mutator with the errorsPart and the handlerMap
			got := matcher.Mutator(tc.errorsPart, handlerMap)
			assert.Equal(t, tc.expectedHandlerOutput, got, "Mutator returned incorrect output")
			if tc.expectedHandlerName != "" {
				assert.Equal(t, len(invokedArgs[tc.expectedHandlerName]), 1, "Handler should be invoked exactly once")
				assert.Equal(t, tc.expectedHandlerArgs, invokedArgs[tc.expectedHandlerName][0], "Handler was invoked with incorrect arguments")
			}
		})
	}
}
