package unmarshaler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type adpackerResp struct {
	Data []coupangResp `json:"data"`
}

type Adpacker struct{}

func (s *Adpacker) UnmarshalResponse(body []byte) ([]PartnerResp, error) {
	resp := &adpackerResp{}
	if err := json.Unmarshal(body, resp); err != nil {
		log.WithError(err).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, ErrInvalidFormat
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
	return res, nil
}
