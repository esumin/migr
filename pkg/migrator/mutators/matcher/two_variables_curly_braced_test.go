package param_matcher_test

import (
	"reflect"
	"testing"

	parammatcher "mig/pkg/migrator/mutators/matcher"
)

func TestMatchCurlyBracedTwoVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "Valid input with two variables without colons",
			input: []string{
				`"Could not update Deployment{Namespace %s, Name %s}"`,
				"namespace",
				"name",
			},
			expected: []string{
				`"Could not update Deployment"`,
				`"namespace"`,
				"namespace",
				`"name"`,
				"name",
			},
		},
		{
			name: "Valid input with two variables with colons",
			input: []string{
				`"Could not update Deployment{Namespace: %s, Name: %s}"`,
				"namespace",
				"name",
			},
			expected: []string{
				`"Could not update Deployment"`,
				`"namespace"`,
				"namespace",
				`"name"`,
				"name",
			},
		},
		{
			name: "Valid input with mixed colon usage",
			input: []string{
				`"Could not update Deployment{Namespace %s, Name: %s}"`,
				"namespace",
				"name",
			},
			expected: []string{
				`"Could not update Deployment"`,
				`"namespace"`,
				"namespace",
				`"name"`,
				"name",
			},
		},
		{
			name: "Invalid template with one variable",
			input: []string{
				`"Could not update Deployment{Namespace %s}"`,
				"namespace",
				"name",
			},
			expected: nil,
		},
		{
			name: "Invalid template with no variables",
			input: []string{
				`"Could not update Deployment"`,
				"namespace",
				"name",
			},
			expected: nil,
		},
		{
			name: "Extra parameters",
			input: []string{
				`"Could not update Deployment{Namespace %s, Name %s}"`,
				"namespace",
				"name",
				"extra",
			},
			expected: nil,
		},
		{
			name: "Missing parameters",
			input: []string{
				`"Could not update Deployment{Namespace %s, Name %s}"`,
				"namespace",
			},
			expected: nil,
		},
		{
			name: "Different placeholder types",
			input: []string{
				`"Could not update Deployment{Namespace %s, Replicas %d}"`,
				"namespace",
				"replicas",
			},
			expected: nil, // Because the second placeholder is not %s
		},
		{
			name: "Whitespace variations without colons",
			input: []string{
				`"Could not update Deployment{Namespace%s,Name %s}"`,
				"namespace",
				"name",
			},
			expected: nil, // Because missing space before %s in first variable
		},
		{
			name: "Valid input with different variable names",
			input: []string{
				`"Error updating Service{Region %s, ServiceName %s}"`,
				"region",
				"serviceName",
			},
			expected: []string{
				`"Error updating Service"`,
				`"region"`,
				"region",
				`"serviceName"`,
				"serviceName",
			},
		},
		{
			name: "Valid input with additional spaces and colons",
			input: []string{
				`"Failed to modify Pod{Cluster: %s, PodName: %s}"`,
				"cluster",
				"podName",
			},
			expected: []string{
				`"Failed to modify Pod"`,
				`"cluster"`,
				"cluster",
				`"podName"`,
				"podName",
			},
		},
		{
			name: "Valid input with one colon missing",
			input: []string{
				`"Could not update Deployment{Namespace: %s, Name %s}"`,
				"namespace",
				"name",
			},
			expected: []string{
				`"Could not update Deployment"`,
				`"namespace"`,
				"namespace",
				`"name"`,
				"name",
			},
		},
		{
			name: "Valid input with both colons missing",
			input: []string{
				`"Could not update Deployment{Namespace %s, Name %s}"`,
				"namespace",
				"name",
			},
			expected: []string{
				`"Could not update Deployment"`,
				`"namespace"`,
				"namespace",
				`"name"`,
				"name",
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := parammatcher.MatchCurlyBracedTwoVariables(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MatchCurlyBracedTwoVariables(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
