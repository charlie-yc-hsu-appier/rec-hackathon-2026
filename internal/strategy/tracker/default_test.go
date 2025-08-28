package tracker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	tt := []struct {
		name   string
		params Params
		want   string
	}{
		{
			name: "GIVEN valid parameters THEN return the expected URL",
			params: Params{
				TrackingURL: "{product_url}&click_param=test&id={click_id_base64}",
				ProductURL:  "https://product.com/item123",
				ClickID:     "abc123",
			},
			want: "https://product.com/item123&click_param=test&id=YWJjMTIz",
		},
		{
			name: "GIVEN missing placeholders THEN return the expected URL",
			params: Params{
				TrackingURL: "{product_url}&click_param=test",
				ProductURL:  "https://product.com/item123",
				ClickID:     "abc123",
			},
			want: "https://product.com/item123&click_param=test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Default{}
			got := strategy.GenerateTrackingURL(tc.params)
			require.Equal(t, tc.want, got)
		})
	}
}
