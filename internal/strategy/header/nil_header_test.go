package header

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	strategy := &NilHeader{}
	params := Params{}
	result := strategy.GenerateHeaders(params)
	require.Nil(t, result)
}
