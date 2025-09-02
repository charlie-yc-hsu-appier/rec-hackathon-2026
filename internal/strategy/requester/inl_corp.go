package requester

import (
	"fmt"
	"strings"
)

type InlCorp struct {
	SizeCodes map[string]string
}

func (s *InlCorp) GenerateRequestURL(params Params) (string, error) {
	defaultRequester := &Default{}
	url, err := defaultRequester.GenerateRequestURL(params)
	if err != nil {
		return "", err
	}

	if s.SizeCodes != nil {
		sizeCode, err := s.getSizeCode(params.ImgWidth, params.ImgHeight)
		if err != nil {
			return "", err
		}
		url = strings.Replace(url, "{size_code}", sizeCode, 1)
	}

	return url, nil
}

func (s *InlCorp) getSizeCode(width, height int) (string, error) {
	size := fmt.Sprintf("%dx%d", width, height)
	if c, ok := s.SizeCodes[size]; ok {
		return c, nil
	}
	return "", fmt.Errorf("not supported size: %dx%d", width, height)
}
