package requester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	tt := []struct {
		name        string
		params      Params
		want        string
		expectedErr string
	}{
		{
			name: "GIVEN valid parameters THEN return the expected URL",
			params: Params{
				RequestURL: "https://example.com/image?size={width}x{height}&user={user_id_lower}&click_id={click_id_base64}&site_domain={web_host}&app_bundleId={bundle_id}&imp_adType={adtype}&partner_id={partner_id}",
				UserID:     "TestUser",
				ImgWidth:   200,
				ImgHeight:  100,
				ClickID:    "test-id",
				WebHost:    "http://example.com/query?param1=123&param2=456",
				BundleID:   "com.example.app",
				AdType:     1,
				PartnerID:  "kakao_kr",
			},
			want: "https://example.com/image?size=200x100&user=testuser&click_id=dGVzdC1pZA&site_domain=http%3A%2F%2Fexample.com%2Fquery%3Fparam1%3D123%26param2%3D456&app_bundleId=com.example.app&imp_adType=1&partner_id=kakao_kr",
		},
		{
			name: "GIVEN missing placeholders THEN return the expected URL",
			params: Params{
				RequestURL: "https://example.com/image/user/abc",
				UserID:     "User",
				ImgWidth:   50,
				ImgHeight:  50,
			},
			want: "https://example.com/image/user/abc",
		},
		{
			name: "GIVEN URL with {subid} but SubID not provided THEN return error",
			params: Params{
				RequestURL: "https://example.com/image?subid={subid}",
				ImgWidth:   300,
				ImgHeight:  300,
			},
			expectedErr: "subID not provided (image size: 300x300)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Default{}
			got, err := strategy.GenerateRequestURL(tc.params)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err.Error())
				require.Empty(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}
