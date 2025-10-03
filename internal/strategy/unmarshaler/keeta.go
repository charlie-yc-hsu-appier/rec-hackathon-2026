package unmarshaler

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type keetaResp struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data keetaData `json:"data"`
}

type keetaData struct {
	Bid   bool        `json:"bid"`
	Items []keetaItem `json:"items"`
}

type keetaItem struct {
	Id        string `json:"id"`
	Deeplink  string `json:"deeplink"`
	Price     string `json:"price"`
	SalePrice string `json:"salePrice"`
	Currency  string `json:"currency"`
}

type Keeta struct{}

func (s *Keeta) UnmarshalResponse(ctx context.Context, body []byte) ([]PartnerResp, error) {
	var resp keetaResp
	if err := json.Unmarshal(body, &resp); err != nil {
		log.WithContext(ctx).Errorf("fail to unmarshal response body: %s", string(body))
		return nil, newInvalidFormatError(body)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("resp code invalid. code: %d, msg: %s", resp.Code, resp.Msg)
	}
	if len(resp.Data.Items) == 0 {
		return nil, ErrNoProducts
	}

	res := make([]PartnerResp, 0, len(resp.Data.Items))
	for _, item := range resp.Data.Items {
		res = append(res, PartnerResp{
			ProductID:        item.Id,
			ProductURL:       item.Deeplink,
			ProductPrice:     item.Price,
			ProductSalePrice: item.SalePrice,
			ProductCurrency:  item.Currency,
		})
	}
	return res, nil
}
