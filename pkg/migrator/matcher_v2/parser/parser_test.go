package parser_test

import (
	"testing"

	"mig/pkg/migrator/matcher_v2/parser"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		input          string
		wantPrefix     string
		wantErrorsPart string
		wantSuffix     string
		expectError    bool
	}{
		// Simple wraps
		{
			input:          `return errors.Wrap(err, "Failed to get Kopia API server secrets")`,
			wantPrefix:     `return `,
			wantErrorsPart: `errors.Wrap(err, "Failed to get Kopia API server secrets")`,
			wantSuffix:     ``,
			expectError:    false,
		},
		// Simple error creation wraps
		{
			input:          `return errors.New("Failed to get Kopia API server secrets"), nil`,
			wantPrefix:     `return `,
			wantErrorsPart: `errors.New("Failed to get Kopia API server secrets")`,
			wantSuffix:     `, nil`,
			expectError:    false,
		},
		// Wraps with named parameters using fmt.Sprintf
		{
			input:          `return nil, errors.Wrap(err, fmt.Sprintf("Error fetching secret %s from namespace %s", ref.Name, ref.Namespace))`,
			wantPrefix:     `return nil, `,
			wantErrorsPart: `errors.Wrap(err, fmt.Sprintf("Error fetching secret %s from namespace %s", ref.Name, ref.Namespace))`,
			wantSuffix:     ``,
			expectError:    false,
		},
		// Wraps with different kinds of brackets
		{
			input:          `return nil, nil, errors.Wrapf(err, "Failed to parse function version {%s}", version)`,
			wantPrefix:     `return nil, nil, `,
			wantErrorsPart: `errors.Wrapf(err, "Failed to parse function version {%s}", version)`,
			wantSuffix:     ``,
			expectError:    false,
		},
		// Names of unnamed parameters guessed from variable names
		{
			input:          `return errors.Wrapf(err, "Error while Pinging the database %s, %s", stderr, err)`,
			wantPrefix:     `return `,
			wantErrorsPart: `errors.Wrapf(err, "Error while Pinging the database %s, %s", stderr, err)`,
			wantSuffix:     ``,
			expectError:    false,
		},
		// Keeping parameters and using fmt.Sprintf
		{
			input:          `return false, errors.Errorf("Error: ca.crt not found in the cluster credential %s-client-secret", c.chart.Release)`,
			wantPrefix:     `return false, `,
			wantErrorsPart: `errors.Errorf("Error: ca.crt not found in the cluster credential %s-client-secret", c.chart.Release)`,
			wantSuffix:     ``,
			expectError:    false,
		},
		// Line without errors
		{
			input:          `return nil, fmt.Errorf("An error occurred")`,
			wantPrefix:     `return nil, fmt.Errorf("An error occurred")`,
			wantErrorsPart: ``,
			wantSuffix:     ``,
			expectError:    false,
		},
		// Complex line
		{
			input:          `if err != nil { return errors.Wrap(err, "Failed to process data") }`,
			wantPrefix:     `if err != nil { return `,
			wantErrorsPart: `errors.Wrap(err, "Failed to process data")`,
			wantSuffix:     ` }`,
			expectError:    false,
		},
		// Line with no errors invocation
		{
			input:          `fmt.Println("Hello, world!")`,
			wantPrefix:     `fmt.Println("Hello, world!")`,
			wantErrorsPart: ``,
			wantSuffix:     ``,
			expectError:    false,
		},
		// Line with invalid syntax
		{
			input:          `return errors.Wrap(err, "Unclosed string literal)`,
			wantPrefix:     ``,
			wantErrorsPart: ``,
			wantSuffix:     ``,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			prefix, errorsPart, suffix, err := parser.ParseLine(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseLine(%q) error = %v, expectError = %v", tt.input, err, tt.expectError)
				return
			}
			if prefix != tt.wantPrefix || errorsPart != tt.wantErrorsPart || suffix != tt.wantSuffix {
				t.Errorf("ParseLine(%q) =\n  prefix      %q\n  errorsPart  %q\n  suffix      %q\nwant:\n  prefix      %q\n  errorsPart  %q\n  suffix      %q",
					tt.input, prefix, errorsPart, suffix, tt.wantPrefix, tt.wantErrorsPart, tt.wantSuffix)
			}
		})
	}
}
