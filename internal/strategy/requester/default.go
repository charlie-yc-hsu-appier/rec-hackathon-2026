package requester

import (
	"fmt"
	urlpkg "net/url"
	"rec-vendor-api/internal/strategy/utils"
	"strconv"
	"strings"
)

type Default struct{}

// TODO: all the query parameters should be URL-encoded
func (s *Default) GenerateRequestURL(params Params) (string, error) {
	url := params.RequestURL
	url = strings.Replace(url, "{width}", strconv.Itoa(params.ImgWidth), 1)
	url = strings.Replace(url, "{height}", strconv.Itoa(params.ImgHeight), 1)
	url = strings.Replace(url, "{user_id_lower}", strings.ToLower(params.UserID), 1)
	url = strings.Replace(url, "{click_id_base64}", utils.EncodeClickID(params.ClickID), 1)
	url = strings.Replace(url, "{web_host}", urlpkg.QueryEscape(params.WebHost), 1)
	url = strings.Replace(url, "{bundle_id}", urlpkg.QueryEscape(params.BundleID), 1)
	url = strings.Replace(url, "{adtype}", strconv.Itoa(params.AdType), 1)
	url = strings.Replace(url, "{partner_id}", params.PartnerID, 1)

	// TODO(AI-28426): Consolidate error handling for all missing required parameters
	// Check if URL contains {subid} placeholder but SubID is not provided
	if strings.Contains(url, "{subid}") && params.SubID == "" {
		return "", fmt.Errorf("subID not provided (image size: %dx%d)", params.ImgWidth, params.ImgHeight)
	}
	url = strings.Replace(url, "{subid}", params.SubID, 1)
	return url, nil
}
