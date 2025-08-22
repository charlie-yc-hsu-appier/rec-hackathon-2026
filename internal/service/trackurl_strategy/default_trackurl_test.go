package trackurl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultTrackingURLStrategy(t *testing.T) {
	tt := []struct {
		name        string
		trackingURL string
		productURL  string
		clickID     string
		want        string
	}{
		{
			name:        "GIVEN valid parameters THEN return the expected URL",
			trackingURL: "https://track.com?url={product_url}&id={click_id_base64}",
			productURL:  "https://product.com/item123",
			clickID:     "abc123",
			want:        "https://track.com?url=https://product.com/item123&id=YWJjMTIz",
		},
		{
			name:        "GIVEN missing placeholders THEN return the expected URL",
			trackingURL: "https://track.com?url={product_url}",
			productURL:  "https://product.com/item123",
			clickID:     "abc123",
			want:        "https://track.com?url=https://product.com/item123",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &DefaultTrackingURLStrategy{}
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
