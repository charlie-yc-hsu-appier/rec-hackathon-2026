package requester

import (
	"rec-vendor-api/internal/strategy/utils"
	"strconv"
	"strings"
)

type Default struct {
	SizeCodes map[string]string
}

func (s *Default) GenerateRequestURL(params Params) (string, error) {
	url := params.RequestURL
	url = strings.Replace(url, "{width}", strconv.Itoa(params.ImgWidth), 1)
	url = strings.Replace(url, "{height}", strconv.Itoa(params.ImgHeight), 1)
	url = strings.Replace(url, "{user_id_lower}", strings.ToLower(params.UserID), 1)
	url = strings.Replace(url, "{click_id_base64}", utils.EncodeClickID(params.ClickID), 1)
	return url, nil
}
