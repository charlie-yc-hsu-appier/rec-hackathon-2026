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
	rResp := &adpackerResp{}
	if err := json.Unmarshal(body, rResp); err != nil {
		log.WithError(err).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, ErrInvalidFormat
	}

	res := make([]PartnerResp, 0, len(rResp.Data))
	for _, item := range rResp.Data {
		res = append(res, PartnerResp{
			ProductID:    strconv.Itoa(item.ProductID),
			ProductURL:   item.ProductURL,
			ProductImage: item.ProductImage,
		})
	}
	return res, nil
}
