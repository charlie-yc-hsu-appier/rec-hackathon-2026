package url

import (
	"fmt"
	urlpkg "net/url"
	"rec-vendor-api/internal/config"
	"rec-vendor-api/internal/strategy/utils"
	"regexp"
	"strconv"
	"strings"
)

var (
	MacroRegExp = regexp.MustCompile(`\{[^}]*\}`)
)

// Default Strategy: Replace macros in URL and query values with values from Params
type Default struct{}

func (s *Default) GenerateURL(urlPattern config.URLPattern, params Params) (string, error) {
	url := urlPattern.URL
	queries := urlPattern.Queries

	url, err := s.replaceMacros(url, params)
	if err != nil {
		return "", err
	}

	// by parsing the url, we can handle
	//   - existing query parameters in the url
	//   - url encoding of replaced macro in url/path
	parsedURL, err := urlpkg.Parse(url)
	if err != nil {
		return "", err
	}

	q := parsedURL.Query()
	// note that only the values of the queries can have macros
	for _, query := range queries {
		value, err := s.replaceMacros(query.Value, params)
		if err != nil {
			return "", err
		}
		q.Set(query.Key, value)
	}
	parsedURL.RawQuery = q.Encode()

	return parsedURL.String(), nil
}

func (s *Default) replaceMacros(str string, params Params) (string, error) {
	matches := MacroRegExp.FindAllString(str, -1)
	for _, macro := range matches {
		value, err := s.getMacroValue(macro, params)
		if err != nil {
			return "", err
		}
		str = strings.Replace(str, macro, value, 1)
	}
	return str, nil
}

// TODO replace error with dedicated error type so that controller can render 400
func (s *Default) getMacroValue(macro string, params Params) (string, error) {
	switch macro {
	case "{width}":
		if params.ImgWidth == 0 {
			return "", fmt.Errorf("ImgWidth not provided")
		}
		return strconv.Itoa(params.ImgWidth), nil
	case "{height}":
		if params.ImgHeight == 0 {
			return "", fmt.Errorf("ImgHeight not provided")
		}
		return strconv.Itoa(params.ImgHeight), nil
	case "{user_id_lower}":
		if params.UserID == "" {
			return "", fmt.Errorf("UserID not provided")
		}
		return strings.ToLower(params.UserID), nil
	case "{click_id_base64}":
		if params.ClickID == "" {
			return "", fmt.Errorf("ClickID not provided")
		}
		return utils.EncodeClickID(params.ClickID), nil
	case "{web_host}":
		return params.WebHost, nil
	case "{bundle_id}":
		return params.BundleID, nil
	case "{adtype}":
		if params.AdType == 0 {
			return "", fmt.Errorf("AdType not provided")
		}
		return strconv.Itoa(params.AdType), nil
	case "{partner_id}":
		if params.PartnerID == "" {
			return "", fmt.Errorf("PartnerID not provided")
		}
		return params.PartnerID, nil
	case "{subid}":
		if params.SubID == "" {
			return "", fmt.Errorf("subID not provided")
		}
		return params.SubID, nil
	case "{product_url}":
		if params.ProductURL == "" {
			return "", fmt.Errorf("ProductURL not provided")
		}
		return params.ProductURL, nil
	default:
		return "", fmt.Errorf("unknown macro: %s", macro)
	}
}
