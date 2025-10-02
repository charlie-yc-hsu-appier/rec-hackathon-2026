package unmarshaler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeeta(t *testing.T) {
	tt := []struct {
		name        string
		input       []byte
		want        []PartnerResp
		wantedError error
	}{
		{
			name:  "GIVEN valid JSON THEN return the expected struct",
			input: []byte(`{"code":0,"msg":"success","data":{"bid":true,"items":[{"id":"123","deeplink":"https://deeplink.com/123","price":"1000","salePrice":"900","currency":"KRW"}]} }`),
			want:  []PartnerResp{{ProductID: "123", ProductURL: "https://deeplink.com/123", ProductPrice: "1000", ProductSalePrice: "900", ProductCurrency: "KRW"}},
		},
		{
			name:        "GIVEN invalid JSON THEN return an error",
			input:       []byte("invalid json"),
			wantedError: errors.New("invalid format. body: invalid json"),
		},
		{
			name:        "GIVEN error code in JSON THEN return an error",
			input:       []byte(`{"code":1,"msg":"error","data":{"bid":true,"items":[]}}`),
			wantedError: errors.New("resp code invalid. code: 1, msg: error"),
		},
		{
			name:        "GIVEN empty items THEN return an error",
			input:       []byte(`{"code":0,"msg":"success","data":{"bid":true,"items":[]}}`),
			wantedError: ErrNoProducts,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Keeta{}
			got, err := strategy.UnmarshalResponse(context.Background(), tc.input)
			if tc.wantedError != nil {
				require.Equal(t, tc.wantedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}
