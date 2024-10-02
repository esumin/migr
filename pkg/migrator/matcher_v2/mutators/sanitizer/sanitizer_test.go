package sanitizer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mig/pkg/migrator/matcher_v2/mutators/sanitizer"
)

// TestSanitizeString tests the SanitizeString function with various inputs.
func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Example with extra spaces and space before dot",
			input:    "Error  getting pod and container name .",
			expected: "Error getting pod and container name.",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "   Hello World   ",
			expected: "Hello World",
		},
		{
			name:     "No unnecessary spaces",
			input:    "Clean sentence.",
			expected: "Clean sentence.",
		},
		{
			name:     "Only spaces",
			input:    "     ",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Tabs and newlines",
			input:    "\tThis is\t a test.\n",
			expected: "\tThis is\t a test.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
