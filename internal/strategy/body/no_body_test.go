package body

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoBody(t *testing.T) {
	tt := []struct {
		name   string
		params Params
	}{
		{
			name:   "GIVEN empty params THEN return nil",
			params: Params{},
		},
		{
			name: "GIVEN valid params THEN return nil",
			params: Params{
				UserID:    "TestUser123",
				ClickID:   "click-id-with-special@chars#123",
				ImgWidth:  1200,
				ImgHeight: 627,
				BundleID:  "com.example.app",
				SubID:     "sub-id-456",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &NoBody{}
			result := strategy.GenerateBody(tc.params)
			require.Nil(t, result)
		})
	}
}
