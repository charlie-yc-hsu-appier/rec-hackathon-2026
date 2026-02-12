package unmarshaler

import (
	"context"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type vendorPostResp struct {
	RCode    string           `json:"rCode"`
	RMessage string           `json:"rMessage"`
	Data     vendorPostResult `json:"data"`
}

type vendorPostResult struct {
	Result []coupangResp `json:"result"`
}

type VendorPost struct{}

func (s *VendorPost) UnmarshalResponse(ctx context.Context, body []byte) ([]PartnerResp, error) {
	rResp := &vendorPostResp{}
	if err := json.Unmarshal(body, rResp); err != nil {
		log.WithContext(ctx).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, newInvalidFormatError(body)
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
