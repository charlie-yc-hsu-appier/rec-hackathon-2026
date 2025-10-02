package unmarshaler

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNoProducts       = errors.New("no products were returned")
	ErrInvalidProductID = errors.New("only a product with ID 0 was returned")
)

type PartnerResp struct {
	ProductID        string
	ProductURL       string
	ProductImage     string
	ProductPrice     string
	ProductSalePrice string
	ProductCurrency  string
}

//go:generate mockgen -source=./interface.go -destination=./interface_mock.go -package=unmarshaler

type Strategy interface {
	UnmarshalResponse(ctx context.Context, body []byte) ([]PartnerResp, error)
}

func newInvalidFormatError(b []byte) error {
	s := string(b)
	runes := []rune(s)

	if len(runes) > 20 {
		s = string(runes[:20]) + "..."
	}
	return fmt.Errorf("invalid format. body: %s", s)
}
