package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeClickID(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "GIVEN valid string THEN expect base64 encoded string",
			input:    "abc123",
			expected: "YWJjMTIz",
		},
		{
			name:     "GIVEN empty string THEN expect empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range tt {
		result := EncodeClickID(tc.input)
		require.Equal(t, tc.expected, result)
	}
}
