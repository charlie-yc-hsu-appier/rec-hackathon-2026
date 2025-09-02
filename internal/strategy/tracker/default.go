package tracker

import (
	"net/url"
	"rec-vendor-api/internal/strategy/utils"
	"strings"
)

type Default struct{}

func (s *Default) GenerateTrackingURL(params Params) string {
	trackingURL := params.TrackingURL
	trackingURL = strings.Replace(trackingURL, "{product_url}", params.ProductURL, 1)
	trackingURL = strings.Replace(trackingURL, "{encoded_product_url}", url.QueryEscape(params.ProductURL), 1)
	trackingURL = strings.Replace(trackingURL, "{click_id_base64}", utils.EncodeClickID(params.ClickID), 1)
	trackingURL = strings.Replace(trackingURL, "{user_id_lower}", strings.ToLower(params.UserID), 1)
	return trackingURL
}
