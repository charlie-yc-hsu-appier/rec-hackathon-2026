package unmarshaler

type CoupangPartnerResp struct {
	ProductID    int    `json:"productId"`
	ProductURL   string `json:"productUrl"`
	ProductImage string `json:"productImage"`
}
type Strategy interface {
	UnmarshalResponse(body []byte) (*[]CoupangPartnerResp, error)
}
