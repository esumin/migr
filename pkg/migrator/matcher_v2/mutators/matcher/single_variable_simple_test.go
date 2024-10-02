package param_matcher_test

import (
	"reflect"
	"testing"

	parammatcher "mig/pkg/migrator/matcher_v2/mutators/matcher"
)

func TestMatchOneVariableSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "Single variable without object field",
			input: []string{
				`"Unable to parse sizeFormat %s"`,
				`sizeFmt`,
			},
			expected: []string{
				`"Unable to parse sizeFormat"`,
				`"sizeFormat"`,
				`sizeFmt`,
			},
		},
		{
			name: "Single variable with different placeholder type",
			input: []string{
				`"Unable to create PVC %v"`,
				`pvc`,
			},
			expected: []string{
				`"Unable to create PVC"`,
				`"PVC"`,
				`pvc`,
			},
		},
		{
			name: "Single variable with object field",
			input: []string{
				`"Failed to create job %s"`,
				`job.name`,
			},
			expected: []string{
				`"Failed to create job"`,
				`"job"`,
				`job.name`,
			},
		},
		{
			name: "Single variable with acronym",
			input: []string{
				`"Failed to get PV %s"`,
				`pvName`,
			},
			expected: []string{
				`"Failed to get PV"`,
				`"PV"`,
				`pvName`,
			},
		},
		{
			name: "Single variable with contextual name",
			input: []string{
				`"Unable to create PV for volume %v"`,
				`pv`,
			},
			expected: []string{
				`"Unable to create PV for volume"`,
				`"volume"`,
				`pv`,
			},
		},
		{
			name: "Corner case with multiple placeholders",
			input: []string{
				`"%s %s"`,
				`errAccessingNode`,
				`n[0]`,
			},
			expected: nil, // Expecting no match
		},
		{
			name: "Parameter is in the middle of text, name should be taken as usually from word before placeholder",
			input: []string{
				`"Error waiting for application %s to be ready to reset it"`,
				"c.name",
			},
			expected: []string{
				`"Error waiting for application to be ready to reset it"`,
				`"application"`,
				`c.name`,
			},
		},
		{
			name:  "Parameter name blacklist, should not be taken as variable name",
			input: []string{`"Failed to uninstall %s helm release"`, "cb.chart.Release"},
			expected: []string{
				`"Failed to uninstall helm release"`,
				`"release"`,
				`cb.chart.Release`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := parammatcher.MatchOneVariableSimple(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MatchOneVariableSimple(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
