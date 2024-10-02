package matcher_v2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	matcher "mig/pkg/migrator/matcher_v2"
)

func TestHandleLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Wrap to Wrap without parameters",
			input:    `return nil, errors.Wrap(err, "Failed to get controller namespace")`,
			expected: `return nil, errkit.Wrap(err, "Failed to get controller namespace")`,
		},
		{
			name:     "Wrapf to Wrap without parameters",
			input:    `return nil, errors.Wrapf(err, "Failed to get controller namespace")`,
			expected: `return nil, errkit.Wrap(err, "Failed to get controller namespace")`,
		},
		{
			name:     "Wrapf to Wrap with single parameter 1",
			input:    `return errors.Wrapf(err, "Failed to get PVC %s", pvcName)`,
			expected: `return errkit.Wrap(err, "Failed to get PVC", "PVC", pvcName)`,
		},
		{
			name:     "Wrapf to Wrap with single parameter 2",
			input:    `return errors.Wrapf(err, "Failed to get PV %s", pvName)`,
			expected: `return errkit.Wrap(err, "Failed to get PV", "PV", pvName)`,
		},
		{
			name:     "Wrapf to Wrap with single parameter 3",
			input:    `return errors.Wrapf(err, "Failed to create volume spec for job %s", job.name)`,
			expected: `return errkit.Wrap(err, "Failed to create volume spec for job", "job", job.name)`,
		},
		{
			name:     "Wrapf to Wrap with multiple variables 1",
			input:    `return errors.Wrapf(err, "Could not get Statefulset{Namespace %s, Name: %s}", namespace, name)`,
			expected: `return errkit.Wrap(err, "Could not get Statefulset", "namespace", namespace, "name", name)`,
		},
		{
			name:     "Wrapf to Wrap with multiple variables 2",
			input:    `return nil, errors.Wrapf(err, "Failed to get pod from podOptions. Namespace: %s, NameFmt: %s", opts.Namespace, opts.GenerateName)`,
			expected: `return nil, errkit.Wrap(err, "Failed to get pod from podOptions.", "namespace", opts.Namespace, "nameFmt", opts.GenerateName)`,
		},
		{
			name:     "Corner Case Wrapf - keep as is",
			input:    `return errors.Wrapf(err, "%s %s", errAccessingNode, n[0])`,
			expected: `return errors.Wrapf(err, "%s %s", errAccessingNode, n[0]) // TODO: migrate manually`,
		},
		{
			name:     "Errorf to New without parameters",
			input:    `return errors.Errorf("Failed to get source")`,
			expected: `return errkit.New("Failed to get source")`,
		},
		{
			name:     "Errorf to New with two parameters - keep first text, remove second",
			input:    `return errors.Errorf("Pod %s failed. Pod status: %s", name, p.Status.String())`,
			expected: `return errkit.New(fmt.Sprintf("Pod %s failed. Pod status: %s", name, p.Status.String()))`,
		},
		{
			name:     "Errorf to New with two parameters",
			input:    `return errors.Errorf("Failed to create content, Volumesnapshot: %s, Error: %v", snap.GetName(), err)`,
			expected: `return errkit.New(fmt.Sprintf("Failed to create content, Volumesnapshot: %s, Error: %v", snap.GetName(), err))`,
		},
		{
			name:     "Errorf to New with single parameter parameter",
			input:    `return errors.Errorf("Invalid secret name %s, it should not be of the form namespace/name )", repositoryPassword)`,
			expected: `return errkit.New(fmt.Sprintf("Invalid secret name %s, it should not be of the form namespace/name )", repositoryPassword))`,
		},
		{
			name:     "New to New",
			input:    `		return errors.New(PasswordIncorrect)`,
			expected: `		return errkit.New(PasswordIncorrect)`,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.HandleLine(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
