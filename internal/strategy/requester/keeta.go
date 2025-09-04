package requester

import (
	"fmt"
	urlpkg "net/url"
	"rec-vendor-api/internal/strategy/utils"
	"strings"
)

type Keeta struct {
	SceneType    string
	Ver          string
	ChannelToken string
}

// GenerateRequestURL generates the request URL for keeta vendor.
//
// AI-26810: The parameters in the URL query string must be added in the order of their dictionary (alphabetical) keys,
// to facilitate subsequent signature parameter processing.
func (s *Keeta) GenerateRequestURL(params Params) (string, error) {
	queryParametersMap := map[string]string{
		"reqId":        urlpkg.QueryEscape(params.ClickID),
		"ip":           urlpkg.QueryEscape(params.ClientIP),
		"campaignId":   urlpkg.QueryEscape(params.KeetaCampaignID),
		"lat":          urlpkg.QueryEscape(params.Latitude),
		"lon":          urlpkg.QueryEscape(params.Longitude),
		"SceneType":    s.SceneType,
		"Ver":          s.Ver,
		"ChannelToken": s.ChannelToken,
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
	url := params.RequestURL
	if queryString != "" {
		url = url + "?" + queryString
	}
	return url, nil
}
