package requester

import (
	"strconv"
	"strings"
)

type Default struct{}

func (s *Default) GenerateRequestURL(params Params) string {
	url := params.RequestURL
	url = strings.Replace(url, "{width}", strconv.Itoa(params.ImgWidth), 1)
	url = strings.Replace(url, "{height}", strconv.Itoa(params.ImgHeight), 1)
	url = strings.Replace(url, "{user_id_lower}", strings.ToLower(params.UserID), 1)
	return url
}
