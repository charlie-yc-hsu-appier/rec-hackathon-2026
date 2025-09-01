package tracker

import (
	"encoding/base64"
	"net/url"
	"strings"
)

type Default struct{}

func (s *Default) GenerateTrackingURL(params Params) string {
	trackingURL := params.TrackingURL
	trackingURL = strings.Replace(trackingURL, "{product_url}", params.ProductURL, 1)
	trackingURL = strings.Replace(trackingURL, "{encoded_product_url}", url.QueryEscape(params.ProductURL), 1)
	trackingURL = strings.Replace(trackingURL, "{click_id_base64}", encodeClickID(params.ClickID), 1)
	trackingURL = strings.Replace(trackingURL, "{user_id_lower}", strings.ToLower(params.UserID), 1)
	return trackingURL
}

func encodeClickID(clickID string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(clickID))
	return strings.TrimRight(encoded, "=")
}
