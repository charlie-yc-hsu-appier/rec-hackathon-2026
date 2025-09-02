package requester

import (
	"fmt"
	"rec-vendor-api/internal/strategy/utils"
	"strconv"
	"strings"
)

type Default struct {
	SizeCodes map[string]string
}

func (s *Default) GenerateRequestURL(params Params) string {
	url := params.RequestURL
	url = strings.Replace(url, "{width}", strconv.Itoa(params.ImgWidth), 1)
	url = strings.Replace(url, "{height}", strconv.Itoa(params.ImgHeight), 1)
	url = strings.Replace(url, "{user_id_lower}", strings.ToLower(params.UserID), 1)
	url = strings.Replace(url, "{click_id_base64}", utils.EncodeClickID(params.ClickID), 1)
	url = strings.Replace(url, "{size_code}", s.getSizeCode(params.ImgWidth, params.ImgHeight), 1)
	return url
}

func (s *Default) getSizeCode(width, height int) string {
	size := fmt.Sprintf("%dx%d", width, height)
	if c, ok := s.SizeCodes[size]; ok {
		return c
	}
	return s.SizeCodes["300x300"]
}
