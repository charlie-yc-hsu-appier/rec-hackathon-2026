package body

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplace(t *testing.T) {
	tt := []struct {
		name   string
		params Params
		want   replaceBody
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
			want: replaceBody{
				App: replaceApp{
					BundleID: "com.example.app",
				},
				Device: replaceDevice{
					ID:  "testuser123",
					Lmt: 0,
				},
				Imp: replaceImp{
					ImageSize: "1200x627",
				},
				Affiliate: replaceAffiliate{
					SubID:    "sub-id-456",
					SubParam: "Y2xpY2staWQtd2l0aC1zcGVjaWFsQGNoYXJzIzEyMw",
				},
				User: replaceUser{
					Puid: "Y2xpY2staWQtd2l0aC1zcGVjaWFsQGNoYXJzIzEyMw",
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
			want: replaceBody{
				App: replaceApp{
					BundleID: "",
				},
				Device: replaceDevice{
					ID:  "",
					Lmt: 0,
				},
				Imp: replaceImp{
					ImageSize: "0x0",
				},
				Affiliate: replaceAffiliate{
					SubID:    "",
					SubParam: "",
				},
				User: replaceUser{
					Puid: "",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Replace{}
			got := strategy.GenerateBody(tc.params)
			require.Equal(t, tc.want, got)
		})
	}
}
