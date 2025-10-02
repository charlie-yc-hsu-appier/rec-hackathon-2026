package unmarshaler

import (
	"context"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type coupangResp struct {
	ProductID    int    `json:"productId"`
	ProductURL   string `json:"productUrl"`
	ProductImage string `json:"productImage"`
}

type CoupangPartner struct{}

func (s *CoupangPartner) UnmarshalResponse(ctx context.Context, body []byte) ([]PartnerResp, error) {
	var resp []coupangResp
	if err := json.Unmarshal(body, &resp); err != nil {
		log.WithContext(ctx).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, newInvalidFormatError(body)
	}

	res := make([]PartnerResp, 0, len(resp))
	for _, item := range resp {
		res = append(res, PartnerResp{
			ProductID:    strconv.Itoa(item.ProductID),
			ProductURL:   item.ProductURL,
			ProductImage: item.ProductImage,
		})
	}
	if len(res) == 1 && res[0].ProductID == "0" {
		return nil, ErrInvalidProductID
	}
	return res, nil
}
