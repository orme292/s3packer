package conf

import (
	"os"
	"testing"
)

func TestExpandHome(t *testing.T) {
	// Get the user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		// If the path starts with "~/", it should replace "~" with the home directory.
		{"~/folder", home + "/folder"},
		// If the path does not start with "~/", it should be unchanged.
		{"/tmp/file", "/tmp/file"},
		{"folder/subfolder", "folder/subfolder"},
		// Edge case: just "~" should not be replaced because the function checks for "~/"
		{"~", "~"},
	}

	for _, tc := range tests {
		got := expandHome(tc.input)
		if got != tc.expected {
			t.Errorf("expandHome(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "helloworld"},   // spaces removed
		{"foo+bar", "foo+bar"},          // allowed characters remain
		{"foo-bar", "foo-bar"},          // dash remains
		{"foo_bar", "foo_bar"},          // underscore remains
		{"foo!bar", "foobar"},           // exclamation removed
		{"123!@#$%^", "123"},            // only digits remain
		{"Hello, World!", "HelloWorld"}, // punctuation removed
		{"--++__", "--++__"},            // multiple allowed symbols remain
		{"", ""},                        // empty string remains empty
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := sanitizeString(tc.input)
			if result != tc.expected {
				t.Errorf("sanitizeString(%q) = %q; want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestTidyStrings(t *testing.T) {
	tests := []struct {
		input         string
		expectedLower string
		expectedUpper string
	}{
		{"  Hello World  ", "hello world", "HELLO WORLD"},
		{"GoLang", "golang", "GOLANG"},
		{"   Mixed CASE   ", "mixed case", "MIXED CASE"},
		{"  trailing space ", "trailing space", "TRAILING SPACE"},
		{"UPPERlower", "upperlower", "UPPERLOWER"},
	}

	for _, tc := range tests {
		// Test tidyLowerString
		if got := tidyLowerString(tc.input); got != tc.expectedLower {
			t.Errorf("tidyLowerString(%q) = %q; want %q", tc.input, got, tc.expectedLower)
		}

		// Test tidyUpperString
		if got := tidyUpperString(tc.input); got != tc.expectedUpper {
			t.Errorf("tidyUpperString(%q) = %q; want %q", tc.input, got, tc.expectedUpper)
		}
	}
}
