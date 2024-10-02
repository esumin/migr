package param_matcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mig/pkg/migrator/matcher_v2/mutators/matcher"
)

func TestMatchTwoVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "Valid input with two variables without colons",
			input: []string{
				`"Failed to create pod. Failed to override pod specs. Namespace %s, NameFmt %s"`,
				`opts.Namespace`,
				`opts.GenerateName`,
			},
			expected: []string{
				`"Failed to create pod. Failed to override pod specs."`,
				`"namespace"`,
				`opts.Namespace`,
				`"nameFmt"`,
				`opts.GenerateName`,
			},
		},
		{
			name: "Valid input with two variables with colons",
			input: []string{
				`"Failed to create pod. Failed to override pod specs. Namespace: %s, NameFmt: %s"`,
				`opts.Namespace`,
				`opts.GenerateName`,
			},
			expected: []string{
				`"Failed to create pod. Failed to override pod specs."`,
				`"namespace"`,
				`opts.Namespace`,
				`"nameFmt"`,
				`opts.GenerateName`,
			},
		},
		{
			name: "Valid input with mixed colon usage",
			input: []string{
				`"Error encountered. Namespace: %s, ServiceName %s"`,
				`opts.Namespace`,
				`opts.ServiceName`,
			},
			expected: []string{
				`"Error encountered."`,
				`"namespace"`,
				`opts.Namespace`,
				`"serviceName"`,
				`opts.ServiceName`,
			},
		},
		{
			name: "Invalid template with one variable",
			input: []string{
				`"Error encountered. Namespace: %s"`,
				`opts.Namespace`,
				`opts.ServiceName`,
			},
			expected: nil,
		},
		{
			name: "Invalid template with no variables",
			input: []string{
				`"Error encountered."`,
				`opts.Namespace`,
				`opts.ServiceName`,
			},
			expected: nil,
		},
		{
			name: "Extra parameters",
			input: []string{
				`"Error encountered. Namespace: %s, ServiceName: %s"`,
				`opts.Namespace`,
				`opts.ServiceName`,
				`extraParam`,
			},
			expected: nil,
		},
		{
			name: "Missing parameters",
			input: []string{
				`"Error encountered. Namespace: %s, ServiceName: %s"`,
				`opts.Namespace`,
			},
			expected: nil,
		},
		{
			name: "Different placeholder types",
			input: []string{
				`"Error encountered. Namespace: %s, Replicas: %d"`,
				`opts.Namespace`,
				`opts.Replicas`,
			},
			expected: nil, // Because the second placeholder is not %s
		},
		{
			name: "Whitespace variations without colons",
			input: []string{
				`"Error encountered. Namespace%s, ServiceName %s"`,
				`opts.Namespace`,
				`opts.ServiceName`,
			},
			expected: nil, // Because missing space before %s in first variable
		},
		{
			name: "Valid input with different variable names",
			input: []string{
				`"Deployment failed. Region: %s, ClusterName: %s"`,
				`opts.Region`,
				`opts.ClusterName`,
			},
			expected: []string{
				`"Deployment failed."`,
				`"region"`,
				`opts.Region`,
				`"clusterName"`,
				`opts.ClusterName`,
			},
		},
		{
			name: "Valid input with additional spaces and colons",
			input: []string{
				`"Operation failed. Zone:   %s, InstanceName:    %s"`,
				`opts.Zone`,
				`opts.InstanceName`,
			},
			expected: []string{
				`"Operation failed."`,
				`"zone"`,
				`opts.Zone`,
				`"instanceName"`,
				`opts.InstanceName`,
			},
		},
		{
			name: "Valid input with one colon missing",
			input: []string{
				`"Operation failed. Zone %s, InstanceName: %s"`,
				`opts.Zone`,
				`opts.InstanceName`,
			},
			expected: []string{
				`"Operation failed."`,
				`"zone"`,
				`opts.Zone`,
				`"instanceName"`,
				`opts.InstanceName`,
			},
		},
		{
			name: "Valid input with both colons missing",
			input: []string{
				`"Operation failed. Zone %s, InstanceName %s"`,
				`opts.Zone`,
				`opts.InstanceName`,
			},
			expected: []string{
				`"Operation failed."`,
				`"zone"`,
				`opts.Zone`,
				`"instanceName"`,
				`opts.InstanceName`,
			},
		},
		{
			name: "Valid input with multiple spaces",
			args: []string{`"Failed to create content, Volumesnapshot: %s, Error: %v"`, "snap.GetName()", "err"},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := param_matcher.MatchTwoVariablesSimple(tt.input)
			assert.Equal(t, tt.expected, result)
			//if !reflect.DeepEqual(result, tt.expected) {
			//	t.Errorf("MatchTwoVariablesSimple(%v) = %v; want %v", tt.input, result, tt.expected)
			//}
		})
	}
}
