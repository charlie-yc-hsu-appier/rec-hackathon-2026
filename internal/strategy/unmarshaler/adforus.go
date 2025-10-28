package unmarshaler

import (
	"context"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type adforusResp struct {
	ProductID    string `json:"productId"`
	ProductName  string `json:"productName"`
	ProductPrice int    `json:"productPrice"`
	ProductImage string `json:"productImage"`
	ProductURL   string `json:"productUrl"`
}

type Adforus struct{}

func (s *Adforus) UnmarshalResponse(ctx context.Context, body []byte) ([]PartnerResp, error) {
	var resp []adforusResp
	if err := json.Unmarshal(body, &resp); err != nil {
		log.WithContext(ctx).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, newInvalidFormatError(body)
	}
	if len(resp) == 0 {
		return nil, ErrNoProducts
	}

	res := make([]PartnerResp, 0, len(resp))
	for _, item := range resp {
		res = append(res, PartnerResp{
			ProductID:        item.ProductID,
			ProductURL:       item.ProductURL,
			ProductSalePrice: strconv.Itoa(item.ProductPrice),
		})
	}

	return res, nil
}
