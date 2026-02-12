package body

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVendorPost(t *testing.T) {
	tt := []struct {
		name   string
		params Params
		want   vendorPostRequestBody
	}{
		{
			name: "GIVEN valid parameters THEN return the expected body structure",
			params: Params{
				UserID:    "TestUser123",
				ClickID:   "click-id-with-special@chars#123",
				ImgWidth:  1200,
				ImgHeight: 627,
				BundleID:  "com.example.app",
				SubID:     "sub-id-456",
			},
			want: vendorPostRequestBody{
				App: vendorPostApp{
					BundleID: "com.example.app",
				},
				Device: vendorPostDevice{
					ID:  "testuser123",
					Lmt: 0,
				},
				Imp: vendorPostImp{
					ImageSize: "1200x627",
				},
				Affiliate: vendorPostAffiliate{
					SubID: "sub-id-456",
				},
			},
		},
		{
			name: "GIVEN empty strings THEN return body with empty string values",
			params: Params{
				UserID:    "",
				ClickID:   "",
				ImgWidth:  0,
				ImgHeight: 0,
				BundleID:  "",
				SubID:     "",
			},
			want: vendorPostRequestBody{
				App: vendorPostApp{
					BundleID: "",
				},
				Device: vendorPostDevice{
					ID:  "",
					Lmt: 0,
				},
				Imp: vendorPostImp{
					ImageSize: "0x0",
				},
				Affiliate: vendorPostAffiliate{
					SubID: "",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &VendorPost{}
			got := strategy.GenerateBody(tc.params)
			require.Equal(t, tc.want, got)
		})
	}
}
