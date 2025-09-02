package unmarshaler

import (
	"encoding/json"
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"
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
		log.WithError(err).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, errors.New("invalid format")
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
