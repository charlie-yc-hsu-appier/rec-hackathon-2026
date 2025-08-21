package requester

import (
	"fmt"
	"strings"
)

type Params struct {
	RequestURL string
	UserID     string
	ImgWidth   int
	ImgHeight  int
}

type DefaultRequestURLStrategy struct{}

func (s *DefaultRequestURLStrategy) GenerateRequestURL(params Params) string {
	url := params.RequestURL
	url = strings.Replace(url, "{width}", fmt.Sprintf("%d", params.ImgWidth), 1)
	url = strings.Replace(url, "{height}", fmt.Sprintf("%d", params.ImgHeight), 1)
	url = strings.Replace(url, "{user_id}", strings.ToLower(params.UserID), 1)
	return url
}
