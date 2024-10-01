package mutators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mig/pkg/migrator/mutators"
)

func TestHandleErrorParams(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		// Case #1
		{
			name:     "Case #1: Simple string without parameters",
			args:     []string{`"Failed to create UUID"`},
			expected: []string{`"Failed to create UUID"`},
		},
		// Case #2
		{
			name: "Case #2: String with multiple parameters",
			args: []string{`"Could not update Deployment{Namespace %s, Name: %s}"`, `namespace`, `name`},
			expected: []string{
				`"Could not update Deployment"`,
				`"Namespace"`, `namespace`,
				`"Name"`, `name`,
			},
		},
		// Case #3
		{
			name: "Case #3: Complex string with multiple parameters",
			args: []string{`"Failed to create pod. Failed to override pod specs. Namespace: %s, NameFmt: %s"`, `opts.Namespace`, `opts.GenerateName`},
			expected: []string{
				`"Failed to create pod. Failed to override pod specs."`,
				`"Namespace"`, `opts.Namespace`,
				`"NameFmt"`, `opts.GenerateName`,
			},
		},
		// Case #4
		{
			name:     "Case #4: Guess parameter name from string template",
			args:     []string{`"Failed to create job %s"`, `job.name`},
			expected: []string{`"Failed to create job"`, `"job"`, `job.name`},
		},
		// Case #5
		{
			name: "Case #5: Guess parameter names from variable names",
			args: []string{`"Failed to validate if PVC %s:%s exists"`, `namespace`, `claimName`},
			expected: []string{
				`"Failed to validate if PVC exists"`,
				`"namespace"`, `namespace`,
				`"claimName"`, `claimName`,
			},
		},
		// Case #6
		{
			name: "Case #6: Guess parameter names from variable names",
			args: []string{`"Failed to find VolumeSnapshot: %s/%s"`, `namespace`, `name`},
			expected: []string{
				`"Failed to find VolumeSnapshot"`,
				`"namespace"`, `namespace`,
				`"name"`, `name`,
			},
		},
		// Case #7
		{
			name: "Case #7: Guess parameter names from string template and variable names",
			args: []string{`"Failed to query PVC %s, Namespace %s"`, `volumeName`, `namespace`},
			expected: []string{
				`"Failed to query PVC"`,
				`"PVC"`, `volumeName`,
				`"namespace"`, `namespace`,
			},
		},
		// Case #8
		{
			name:     "Case #8: Special case with non-quoted first argument",
			args:     []string{`notFoundTmpl`, `"blueprint"`, `p.Blueprint`, `p.Namespace`},
			expected: []string{`notFoundTmpl`, `"blueprint"`, `p.Blueprint`, `p.Namespace`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mutators.HandleErrorParams(tt.args)
			assert.Equal(t, tt.expected, got, "ExtractErrorComponents did not return expected output")
		})
	}
}
