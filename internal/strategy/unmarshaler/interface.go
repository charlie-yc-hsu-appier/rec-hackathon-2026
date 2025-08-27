package unmarshaler

type PartnerResp struct {
	ProductID    string
	ProductURL   string
	ProductImage string
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=unmarshaler

type Strategy interface {
	UnmarshalResponse(body []byte) ([]PartnerResp, error)
}
