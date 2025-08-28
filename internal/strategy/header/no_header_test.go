package header

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoHeader(t *testing.T) {
	strategy := &NoHeader{}
	params := Params{}
	result := strategy.GenerateHeaders(params)
	require.Equal(t, map[string]string{}, result)
}
