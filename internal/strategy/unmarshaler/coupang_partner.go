package unmarshaler

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type coupangResp struct {
	ProductID    int    `json:"productId"`
	ProductURL   string `json:"productUrl"`
	ProductImage string `json:"productImage"`
}

type CoupangPartner struct{}

func (s *CoupangPartner) UnmarshalResponse(body []byte) ([]PartnerResp, error) {
	var resp []coupangResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("invalid format. body: %v", string(body))
	}
	res := make([]PartnerResp, 0, len(resp))
	for _, item := range resp {
		res = append(res, PartnerResp{
			ProductID:    strconv.Itoa(item.ProductID),
			ProductURL:   item.ProductURL,
			ProductImage: item.ProductImage,
		})
	}
	return res, nil
}
