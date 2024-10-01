package mutators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mig/pkg/migrator/mutators"
)

func TestHandleWrap(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Simple Wrap",
			args:     []string{"err", `"Failed to create session"`},
			expected: `errkit.Wrap(err, "Failed to create session")`,
		},
		{
			name:     "Wrap with concatenated string",
			args:     []string{"err", `"failed to read env from dir:"+dir`},
			expected: `errkit.Wrap(err, "failed to read env from dir", "dir", dir)`,
		},
		{
			name:     "Wrap with another concatenated string",
			args:     []string{"err", `"Invalid log level: "+v`},
			expected: `errkit.Wrap(err, "Invalid log level", "level", v)`,
		},
		{
			name:     "Wrap with fmt.Sprintf",
			args:     []string{"err", `fmt.Sprintf("unable to convert parsed count value %s", countStr)`},
			expected: `errkit.Wrap(err, fmt.Sprintf("unable to convert parsed count value %s", countStr))`,
		},
		{
			name:     "Wrap with variable message",
			args:     []string{"err", "message"},
			expected: `errkit.Wrap(err, message)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mutators.HandleWrap(tt.args)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestHandleWrapf(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Simple Wrapf",
			args:     []string{"err", `"Failed to create session"`},
			expected: `errkit.Wrap(err, "Failed to create session")`,
		},
		{
			name:     "Wrapf with concatenated string",
			args:     []string{"err", `"failed to read env from dir:"+dir`},
			expected: `errkit.Wrap(err, "failed to read env from dir", "dir", dir)`,
		},
		{
			name:     "Wrapf with another concatenated string",
			args:     []string{"err", `"Invalid log level: "+v`},
			expected: `errkit.Wrap(err, "Invalid log level", "level", v)`,
		},
		{
			name:     "Wrapf with variable message",
			args:     []string{"err", "message"},
			expected: `errkit.Wrap(err, message)`,
		},
		{
			name:     "Wrapf with multiple variables",
			args:     []string{"err", `"Could not update Deployment{Namespace %s, Name: %s}", namespace, name`},
			expected: `errkit.Wrap(err, "Could not update Deployment", "namespace", namespace, "name", name)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mutators.HandleWrapf(tt.args)
			assert.Equal(t, tt.expected, got)
		})
	}
}
