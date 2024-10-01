package migrator_test

import (
	"testing"

	"mig/pkg/migrator"
)

func TestReplaceErrors(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Simple wraps
		{
			`return errors.Wrap(err, "Failed to get Kopia API server secrets")`,
			`return errkit.Wrap(err, "Failed to get Kopia API server secrets")`,
		},
		// Simple error creation wraps
		{
			`return errors.New("Failed to get Kopia API server secrets"), nil`,
			`return errkit.New("Failed to get Kopia API server secrets"), nil`,
		},
		// Wraps with named parameters using fmt.Sprintf
		{
			`return nil, errors.Wrap(err, fmt.Sprintf("Error fetching secret %s from namespace %s", ref.Name, ref.Namespace))`,
			`return nil, errkit.Wrap(err, "Error fetching secret from namespace", "secret", ref.Name, "namespace", ref.Namespace)`,
		},
		// Wraps with different kinds of brackets
		{
			`return nil, nil, errors.Wrapf(err, "Failed to parse function version {%s}", version)`,
			`return nil, nil, errkit.Wrap(err, "Failed to parse function version", "version", version)`,
		},
		// Names of unnamed parameters guessed from variable names
		{
			`return errors.Wrapf(err, "Error while Pinging the database %s, %s", stderr, err)`,
			`return errkit.Wrap(err, "Error while Pinging the database", "stderr", stderr)`,
		},
		// Keeping parameters and using fmt.Sprintf
		{
			`return false, errors.Errorf("Error: ca.crt not found in the cluster credential %s-client-secret", c.chart.Release)`,
			`return false, errkit.New(fmt.Sprintf("Error: ca.crt not found in the cluster credential %s-client-secret", c.chart.Release))`,
		},
	}

	for _, tt := range tests {
		got := migrator.ReplaceErrors(tt.input)
		if got != tt.want {
			t.Errorf("ReplaceErrors(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}
