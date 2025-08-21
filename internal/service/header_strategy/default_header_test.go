package header

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoHeaderStrategy(t *testing.T) {
	strategy := &NoHeaderStrategy{}
	params := Params{UserID: "user1", ClickID: "click1"}
	result := strategy.GenerateHeaders(params)
	require.Nil(t, result)
}
