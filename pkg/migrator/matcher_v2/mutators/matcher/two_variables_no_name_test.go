package param_matcher_test

import (
	"reflect"
	"testing"

	parammatcher "mig/pkg/migrator/matcher_v2/mutators/matcher"
)

func TestMatchTwoVariablesNoName(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "Valid input with two variables without colons",
			input: []string{
				`"Failed to validate if PVC %s:%s exists"`,
				`namespace`,
				`claimName`,
			},
			expected: []string{
				`"Failed to validate if PVC exists"`,
				`"namespace"`,
				`namespace`,
				`"claimName"`,
				`claimName`,
			},
		},
		{
			name: "Valid input with two variables with colons",
			input: []string{
				`"Failed to validate if PVC %s:%s exists"`,
				`opts.Namespace`,
				`opts.ClaimName`,
			},
			expected: []string{
				`"Failed to validate if PVC exists"`,
				`"namespace"`,
				`opts.Namespace`,
				`"claimName"`,
				`opts.ClaimName`,
			},
		},
		{
			name: "Valid input with different separators",
			input: []string{
				`"Failed to find VolumeSnapshot: %s/%s"`,
				`opts.Namespace`,
				`opts.GenerateName`,
			},
			expected: []string{
				`"Failed to find VolumeSnapshot"`,
				`"namespace"`,
				`opts.Namespace`,
				`"generateName"`,
				`opts.GenerateName`,
			},
		},
		{
			name: "Valid input with mixed separators",
			input: []string{
				`"Operation failed. PVC %s/%s not found"`,
				`opts.Namespace`,
				`opts.PVCName`,
			},
			expected: []string{
				`"Operation failed. PVC not found"`,
				`"namespace"`,
				`opts.Namespace`,
				`"PVCName"`,
				`opts.PVCName`,
			},
		},
		{
			name: "Invalid template with one placeholder",
			input: []string{
				`"Failed to validate if PVC %s exists"`,
				`namespace`,
				`claimName`,
			},
			expected: nil,
		},
		{
			name: "Invalid template with no placeholders",
			input: []string{
				`"Failed to validate if PVC exists"`,
				`namespace`,
				`claimName`,
			},
			expected: nil,
		},
		{
			name: "Extra parameters",
			input: []string{
				`"Failed to validate if PVC %s:%s exists"`,
				`namespace`,
				`claimName`,
				`extraParam`,
			},
			expected: nil,
		},
		{
			name: "Missing parameters",
			input: []string{
				`"Failed to validate if PVC %s:%s exists"`,
				`namespace`,
			},
			expected: nil,
		},
		{
			name: "Different placeholder types",
			input: []string{
				`"Failed to validate if PVC %d:%s exists"`,
				`namespace`,
				`claimName`,
			},
			expected: nil, // Because the first placeholder is not %s
		},
		{
			name: "Whitespace variations without separators",
			input: []string{
				`"Failed to validate if PVC%s%s exists"`,
				`namespace`,
				`claimName`,
			},
			expected: nil, // Because missing separators between %s
		},
		{
			name: "Valid input with additional spaces and different separators",
			input: []string{
				`"Failed to validate if PVC  %s : %s exists"`,
				`opts.Namespace`,
				`opts.ClaimName`,
			},
			expected: []string{
				`"Failed to validate if PVC exists"`,
				`"namespace"`,
				`opts.Namespace`,
				`"claimName"`,
				`opts.ClaimName`,
			},
		},
		{
			name: "Valid input with object fields and mixed cases",
			input: []string{
				`"Failed to validate if PVC %s:%s exists"`,
				`opts.Namespace`,
				`opts.ClaimName`,
			},
			expected: []string{
				`"Failed to validate if PVC exists"`,
				`"namespace"`,
				`opts.Namespace`,
				`"claimName"`,
				`opts.ClaimName`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := parammatcher.MatchTwoVariablesNoName(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MatchTwoVariablesNoName(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
