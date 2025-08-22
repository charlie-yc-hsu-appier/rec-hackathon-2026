package trackurl

import (
	"encoding/base64"
	"strings"
)

type Params struct {
	TrackingURL string
	ProductURL  string
	ClickID     string
}

type DefaultTrackingURLStrategy struct{}

func (s *DefaultTrackingURLStrategy) GenerateTrackingURL(params Params) string {
	trackingURL := params.TrackingURL
	trackingURL = strings.Replace(trackingURL, "{product_url}", params.ProductURL, 1)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(params.ClickID))
	trackingURL = strings.Replace(trackingURL, "{click_id_base64}", encoded, 1)
	return trackingURL
}
