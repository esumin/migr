// param_matcher/single_variable_append_test.go
package param_matcher_test

import (
	"reflect"
	"testing"

	parammatcher "mig/pkg/migrator/matcher_v2/mutators/matcher"
)

func TestMatchSingleVariableAppend(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "Match with single variable appended",
			input: []string{
				`"failed to read env from dir:" + dir`,
			},
			expected: []string{
				`"failed to read env from dir"`,
				`"dir"`,
				`dir`,
			},
		},
		{
			name: "Match - more flexible format format",
			input: []string{
				`"failed to read env from dir: " + dirName`,
			},
			expected: []string{
				`"failed to read env from dir"`,
				`"dir"`,
				`dirName`,
			},
		},
		{
			name: "No match - no variable appended",
			input: []string{
				`"failed to read env from dir"`,
			},
			expected: nil,
		},
		{
			name: "No match - multiple variables",
			input: []string{
				`"dir: " + dir + "path: " + path`,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parammatcher.MatchSingleVariableAppend(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MatchSingleVariableAppend(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
