package strutil

import (
	"testing"
)

func TestTrimMultilineWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "  line1  \n  line2  \n\n  line3  ",
			expected: "line1\nline2\nline3",
		},
		{
			input:    "  \n  line1  \n  \n  line2  \n  ",
			expected: "line1\nline2",
		},
		{
			input:    "line1\nline2\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			input:    "  \n  \n  ",
			expected: "",
		},
	}

	for _, test := range tests {
		result := TrimMultilineWhitespace(test.input)
		if result != test.expected {
			t.Errorf("TrimMultilineWhitespace(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
