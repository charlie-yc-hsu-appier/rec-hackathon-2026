package requester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	tt := []struct {
		name   string
		params Params
		want   string
	}{
		{
			name: "GIVEN valid parameters THEN return the expected URL",
			params: Params{
				RequestURL: "https://example.com/image?size={width}x{height}&user={user_id_lower}&click_id={click_id_base64}&site_domain={web_host}&app_bundleId={bundle_id}&imp_adType={adtype}",
				UserID:     "TestUser",
				ImgWidth:   200,
				ImgHeight:  100,
				ClickID:    "test-id",
				WebHost:    "example.com",
				BundleID:   "com.example.app",
				AdType:     1,
			},
			want: "https://example.com/image?size=200x100&user=testuser&click_id=dGVzdC1pZA&site_domain=example.com&app_bundleId=com.example.app&imp_adType=1",
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Default{}
			got, err := strategy.GenerateRequestURL(tc.params)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
