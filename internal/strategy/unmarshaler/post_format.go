package unmarshaler

import (
	"context"
	"encoding/json"
	"strconv"
)

type PostFormat struct{}

type postFormatResponseBody struct {
	RCode    string `json:"rCode"`
	RMessage string `json:"rMessage"`
	Data     struct {
		Result []struct {
			ProductID    int    `json:"productId"`
			ProductURL   string `json:"productUrl"`
			ProductImage string `json:"productImage"`
		} `json:"result"`
	} `json:"data"`
}

func (pf *PostFormat) UnmarshalResponse(ctx context.Context, data []byte) ([]PartnerResp, error) {
	var body postFormatResponseBody
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	var items []PartnerResp
	for _, res := range body.Data.Result {
		items = append(items, PartnerResp{
			ProductID:    strconv.Itoa(res.ProductID),
			ProductImage: res.ProductImage,
			ProductURL:   res.ProductURL,
		})
	}

	return items, nil
}
