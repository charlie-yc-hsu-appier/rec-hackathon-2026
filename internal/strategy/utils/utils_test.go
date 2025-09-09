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

func TestSortedStringKeys(t *testing.T) {
	cases := []struct {
		name  string
		input map[string]string
		want  []string
	}{
		{
			name:  "GIVEN unordered keys THEN returns sorted keys",
			input: map[string]string{"b": "2", "a": "1", "c": "3"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "GIVEN empty map THEN returns empty slice",
			input: map[string]string{},
			want:  []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetSortedStringKeys(tc.input)
			require.Equal(t, tc.want, got)
		})
	}
}
