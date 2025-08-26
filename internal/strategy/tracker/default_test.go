package tracker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	tt := []struct {
		name        string
		trackingURL string
		productURL  string
		clickID     string
		want        string
	}{
		{
			name:        "GIVEN valid parameters THEN return the expected URL",
			trackingURL: "{product_url}&click_param=test&id={click_id_base64}",
			productURL:  "https://product.com/item123",
			clickID:     "abc123",
			want:        "https://product.com/item123&click_param=test&id=YWJjMTIz",
		},
		{
			name:        "GIVEN missing placeholders THEN return the expected URL",
			trackingURL: "{product_url}&click_param=test",
			productURL:  "https://product.com/item123",
			clickID:     "abc123",
			want:        "https://product.com/item123&click_param=test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Default{}
			params := Params{
				TrackingURL: tc.trackingURL,
				ProductURL:  tc.productURL,
				ClickID:     tc.clickID,
			}
			got := strategy.GenerateTrackingURL(params)
			require.Equal(t, tc.want, got)
		})
	}
}
