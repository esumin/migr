package mutators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mig/pkg/migrator/matcher_v2/mutators"
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
			name:     "Wrapf with multiple variables - 1",
			args:     []string{"err", `"Could not update Deployment{Namespace %s, Name: %s}"`, "namespace", "name"},
			expected: `errkit.Wrap(err, "Could not update Deployment", "namespace", namespace, "name", name)`,
		},
		{
			name:     "Wrapf with multiple variables - 2",
			args:     []string{"err", `"Failed to get pod from podOptions. Namespace: %s, NameFmt: %s"`, "opts.Namespace", "opts.GenerateName"},
			expected: `errkit.Wrap(err, "Failed to get pod from podOptions.", "namespace", opts.Namespace, "nameFmt", opts.GenerateName)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mutators.HandleWrapf(tt.args)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestHandleErrorf(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Simple Errorf",
			args:     []string{`"Failed to get source"`},
			expected: `errkit.New("Failed to get source")`,
		},
		{
			name:     "Errorf with two parameters - 1",
			args:     []string{`"Pod %s failed. Pod status: %s"`, "name", "p.Status.String()"},
			expected: `errkit.New(fmt.Sprintf("Pod %s failed. Pod status: %s", name, p.Status.String()))`,
		},
		{
			name:     "Errorf with two parameters - 2",
			args:     []string{`"Failed to create content, Volumesnapshot: %s, Error: %v"`, "snap.GetName()", "err"},
			expected: `errkit.New(fmt.Sprintf("Failed to create content, Volumesnapshot: %s, Error: %v", snap.GetName(), err))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mutators.HandleErrorf(tt.args)
			assert.Equal(t, tt.expected, got)
		})
	}
}
