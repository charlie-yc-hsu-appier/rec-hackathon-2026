package unmarshaler

import (
	"context"
	"encoding/json"
	"strconv"
)

type VendorPost struct{}

type vendorPostResponse struct {
	RCode    string `json:"rCode"`
	RMessage string `json:"rMessage"`
	Data     data   `json:"data"`
}

type data struct {
	Result []vendorPostResult `json:"result"`
}

type vendorPostResult struct {
	ProductID    int    `json:"productId"`
	ProductURL   string `json:"productUrl"`
	ProductImage string `json:"productImage"`
}

func (u *VendorPost) UnmarshalResponse(ctx context.Context, body []byte) ([]PartnerResp, error) {
	var resp vendorPostResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, newInvalidFormatError(body)
	}

	var products []PartnerResp
	for _, item := range resp.Data.Result {
		products = append(products, PartnerResp{
			ProductID:    strconv.Itoa(item.ProductID),
			ProductURL:   item.ProductURL,
			ProductImage: item.ProductImage,
		})
	}
	return products, nil
}
