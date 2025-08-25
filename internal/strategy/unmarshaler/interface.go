package unmarshaler

type CoupangPartnerResp struct {
	ProductID    int    `json:"productId"`
	ProductURL   string `json:"productUrl"`
	ProductImage string `json:"productImage"`
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=unmarshaler

type Strategy interface {
	UnmarshalResponse(body []byte) (*[]CoupangPartnerResp, error)
}
