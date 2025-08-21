package unmarshaler

import (
	"encoding/json"
	"fmt"
)

type CoupangPartnerResp struct {
	ProductID    int    `json:"productId"`
	ProductUrl   string `json:"productUrl"`
	ProductImage string `json:"productImage"`
}

type DefaultUnmarshalStrategy struct{}

func (s *DefaultUnmarshalStrategy) UnmarshalResponse(body []byte) (*[]CoupangPartnerResp, error) {
	var resp []CoupangPartnerResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("invalid format. body: %v", string(body))
	}
	return &resp, nil
}
