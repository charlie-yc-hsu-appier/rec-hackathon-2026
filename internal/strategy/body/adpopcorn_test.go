package body

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdpopcorn(t *testing.T) {
	tt := []struct {
		name   string
		params Params
		want   adpopcornBody
	}{
		{
			name: "GIVEN valid parameters THEN return the expected body structure",
			params: Params{
				UserID:    "TestUser123",
				ClickID:   "click-id-with-special@chars#123",
				ImgWidth:  1200,
				ImgHeight: 627,
				BundleID:  "",
				SubID:     "sub-id-456",
			},
			want: adpopcornBody{
				App: adpopcornApp{
					BundleID: "",
				},
				Device: adpopcornDevice{
					ID:  "testuser123",
					Lmt: 0,
				},
				Imp: adpopcornImp{
					ImageSize: "1200x627",
				},
				Affiliate: adpopcornAffiliate{
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
			want: adpopcornBody{
				App: adpopcornApp{
					BundleID: "",
				},
				Device: adpopcornDevice{
					ID:  "",
					Lmt: 0,
				},
				Imp: adpopcornImp{
					ImageSize: "0x0",
				},
				Affiliate: adpopcornAffiliate{
					SubID: "",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Adpopcorn{}
			got := strategy.GenerateBody(tc.params)
			require.Equal(t, tc.want, got)
		})
	}
}
