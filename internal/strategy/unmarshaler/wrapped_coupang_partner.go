package unmarshaler

import (
	"encoding/json"
	"fmt"
	"strconv"
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
		return nil, fmt.Errorf("invalid format. body: %v", string(body))
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
