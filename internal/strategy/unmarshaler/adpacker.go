package unmarshaler

import (
	"context"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type adpackerResp struct {
	Data []coupangResp `json:"data"`
}

type Adpacker struct{}

func (s *Adpacker) UnmarshalResponse(ctx context.Context, body []byte) ([]PartnerResp, error) {
	resp := &adpackerResp{}
	if err := json.Unmarshal(body, resp); err != nil {
		log.WithContext(ctx).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, newInvalidFormatError(body)
	}
	if len(resp.Data) == 0 {
		return nil, ErrNoProducts
	}

	res := make([]PartnerResp, 0, len(resp.Data))
	for _, item := range resp.Data {
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
