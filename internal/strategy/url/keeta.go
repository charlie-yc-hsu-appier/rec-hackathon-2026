package url

import (
	"fmt"
	urlpkg "net/url"
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy/utils"
	"strings"
)

type KeetaRequest struct {
	SceneType    string
	Ver          string
	ChannelToken string
}

// GenerateRequestURL generates the request URL for keeta vendor.
//
// AI-26810: The parameters in the URL query string must be added in the order of their dictionary (alphabetical) keys,
// to facilitate subsequent signature parameter processing.
func (s *KeetaRequest) GenerateURL(urlPattern config.URLPattern, params Params) (string, map[string]string, error) {
	queryParametersMap := map[string]string{
		"reqId":        urlpkg.QueryEscape(params.ClickID),
		"ip":           urlpkg.QueryEscape(params.ClientIP),
		"campaignId":   urlpkg.QueryEscape(params.KeetaCampaignID),
		"lat":          urlpkg.QueryEscape(params.Latitude),
		"lon":          urlpkg.QueryEscape(params.Longitude),
		"sceneType":    s.SceneType,
		"ver":          s.Ver,
		"channelToken": s.ChannelToken,
		"bizType":      "bType",
	}

	queryParameters := []string{}
	for _, k := range utils.GetSortedStringKeys(queryParametersMap) {
		v := queryParametersMap[k]
		// handle parameter with empty value as instructed by vendor API docs
		//
		// Note that the rest client still appends an "=" character when sending the request (like so ?empty_param=&...)
		// but it is acceptable as long as the query string used to generate the signature is correctly handled
		if v == "" {
			queryParameters = append(queryParameters, k)
		} else {
			queryParameters = append(queryParameters, fmt.Sprintf("%s=%s", k, v))
		}
	}
	queryString := strings.Join(queryParameters, "&")

	// handle empty query string as instructed by vendor API docs
	url := urlPattern.URL
	if queryString != "" {
		url = url + "?" + queryString
	}
	return url, nil, nil
}
