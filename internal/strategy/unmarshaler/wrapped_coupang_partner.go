package unmarshaler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type wrappedResp struct {
	RCode    string        `json:"rCode"`
	RMessage string        `json:"rMessage"`
	Data     []coupangResp `json:"data"`
}

type WrappedCoupangPartner struct{}

func (s *WrappedCoupangPartner) UnmarshalResponse(body []byte) ([]PartnerResp, error) {
	rResp := &wrappedResp{}
	if err := json.Unmarshal(body, rResp); err != nil {
		log.WithError(err).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, errors.New("invalid format")
	}

	if rResp.RCode != "0" {
		return nil, fmt.Errorf("resp code invalid. code: %s, msg: %s", rResp.RCode, rResp.RMessage)
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
