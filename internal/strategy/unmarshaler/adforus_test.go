package unmarshaler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdforus_UnmarshalResponse(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		want        []PartnerResp
		wantedError bool
	}{
		{
			name: "GIVEN valid single product response THEN return parsed product",
			body: []byte(`[
				{
					"productId": "3288378",
					"productName": "추석 이벤트/ 세렌티 1200 모듈 수납장",
					"productPrice": 241000,
					"productUrl": "https://api.linkmine.co.kr/ck.html?app_code=zbkj6Sirtt&sid=39562&deep_link=https%3A%2F%2Flink.ohou.se%2F%40ohouse%2Faffiliate%3Fchannel%3Daffiliate"
				}
			]`),
			want: []PartnerResp{
				{
					ProductID:        "3288378",
					ProductSalePrice: "241000",
					ProductURL:       "https://api.linkmine.co.kr/ck.html?app_code=zbkj6Sirtt&sid=39562&deep_link=https%3A%2F%2Flink.ohou.se%2F%40ohouse%2Faffiliate%3Fchannel%3Daffiliate",
				},
			},
		},
		{
			name: "GIVEN multiple products response THEN return all parsed products",
			body: []byte(`[
				{
					"productId": "3288378",
					"productName": "product1",
					"productPrice": 241000,
					"productImage": "image1.jpg",
					"productUrl": "url1"
				},
				{
					"productId": "1019809",
					"productName": "product2",
					"productPrice": 15740,
					"productImage": "image2.jpg",
					"productUrl": "url2"
				}
			]`),
			want: []PartnerResp{
				{
					ProductID:        "3288378",
					ProductSalePrice: "241000",
					ProductURL:       "url1",
				},
				{
					ProductID:        "1019809",
					ProductSalePrice: "15740",
					ProductURL:       "url2",
				},
			},
		},
		{
			name:        "GIVEN empty response THEN return invalid product ID error",
			body:        []byte(`[]`),
			wantedError: true,
		},
		{
			name:        "GIVEN invalid JSON THEN return error",
			body:        []byte(`invalid json`),
			wantedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unmarshaler := &Adforus{}

			result, err := unmarshaler.UnmarshalResponse(context.Background(), tt.body)

			if tt.wantedError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}
