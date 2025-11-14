package unmarshaler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type replaceResp struct {
	RCode    string `json:"rCode"`
	RMessage string `json:"rMessage"`
	Data     result `json:"data"`
}

type result struct {
	Result []coupangResp `json:"result"`
}

type Replace struct{}

func (s *Replace) UnmarshalResponse(ctx context.Context, body []byte) ([]PartnerResp, error) {
	rResp := &replaceResp{}
	if err := json.Unmarshal(body, rResp); err != nil {
		log.WithContext(ctx).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, newInvalidFormatError(body)
	}
	if rResp.RCode != "0" {
		return nil, fmt.Errorf("resp code invalid. code: %s, msg: %s", rResp.RCode, rResp.RMessage)
	}

	res := make([]PartnerResp, 0, len(rResp.Data.Result))
	for _, item := range rResp.Data.Result {
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
