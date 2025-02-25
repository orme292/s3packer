package distlog

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected zerolog.Level
	}{
		{
			name:     "int 0 => debug",
			input:    0,
			expected: zerolog.DebugLevel, // "0" => debug
		},
		{
			name:     "int 1 => info",
			input:    1,
			expected: zerolog.InfoLevel, // "1" => info
		},
		{
			name:     "int 2 => warn",
			input:    2,
			expected: zerolog.WarnLevel, // "2" => warn
		},
		{
			name:     "int 3 => error",
			input:    3,
			expected: zerolog.ErrorLevel, // "3" => error
		},
		{
			name:     "bool true => warn (parses as '2')",
			input:    true,
			expected: zerolog.WarnLevel, // "2" => warn
		},
		{
			name:     "bool false => warn (parses as '2')",
			input:    false,
			expected: zerolog.WarnLevel, // "2" => warn
		},
		{
			name:     "any string => returns warn (code sets str = \"\" and fails parse)",
			input:    "",
			expected: zerolog.WarnLevel, // because it never sets str=v, so parseLevel("") => error => default to warn
		},
		{
			name:     "invalid type => warn (sets '2')",
			input:    struct{}{},
			expected: zerolog.WarnLevel, // "2" => warn
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLevel(tt.input)
			require.Equal(t, tt.expected, result, "ParseLevel(%v) did not return expected level", tt.input)
		})
	}
}

func TestZerologNumericLevels(t *testing.T) {
	require.Equal(t, -1, int(zerolog.TraceLevel), "Trace level should be -1")
	require.Equal(t, 0, int(zerolog.DebugLevel), "Debug level should be 0")
	require.Equal(t, 1, int(zerolog.InfoLevel), "Info level should be 1")
	require.Equal(t, 2, int(zerolog.WarnLevel), "Warn level should be 2")
	require.Equal(t, 3, int(zerolog.ErrorLevel), "Error level should be 3")
	require.Equal(t, 4, int(zerolog.FatalLevel), "Fatal level should be 4")
	require.Equal(t, 5, int(zerolog.PanicLevel), "Panic level should be 5")
}
