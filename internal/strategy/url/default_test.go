package url

import (
	"rec-vendor-api/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	tt := []struct {
		name        string
		urlPattern  config.URLPattern
		params      Params
		wantURL     string
		expectedErr string
	}{
		{
			name: "GIVEN valid parameters THEN return the expected URL",
			urlPattern: config.URLPattern{
				URL: "https://example.com/image",
				Queries: []config.Query{
					{Key: "size", Value: "{width}x{height}"},
					{Key: "user", Value: "{user_id_lower}"},
					{Key: "click_id", Value: "{click_id_base64}"},
					{Key: "site_domain", Value: "{web_host}"},
					{Key: "app_bundleId", Value: "{bundle_id}"},
					{Key: "imp_adType", Value: "{adtype}"},
					{Key: "partner_id", Value: "{partner_id}"},
				},
			},
			params: Params{
				UserID:    "TestUser",
				ImgWidth:  200,
				ImgHeight: 100,
				ClickID:   "test-id",
				WebHost:   "http://example.com/query?param1=123&param2=456",
				BundleID:  "com.example.app",
				AdType:    1,
				PartnerID: "kakao_kr",
			},
			wantURL: "https://example.com/image?app_bundleId=com.example.app&click_id=dGVzdC1pZA&imp_adType=1&partner_id=kakao_kr&site_domain=http%3A%2F%2Fexample.com%2Fquery%3Fparam1%3D123%26param2%3D456&size=200x100&user=testuser",
		},
		{
			name: "GIVEN missing placeholders THEN return the expected URL",
			urlPattern: config.URLPattern{
				URL: "https://example.com/image/user/abc",
			},
			params: Params{
				UserID:    "User",
				ImgWidth:  50,
				ImgHeight: 50,
			},
			wantURL: "https://example.com/image/user/abc",
		},
		{
			name: "GIVEN URL with {subid} but SubID not provided THEN return error",
			urlPattern: config.URLPattern{
				URL: "https://example.com/image",
				Queries: []config.Query{
					{Key: "subid", Value: "{subid}"},
				},
			},
			params: Params{
				ImgWidth:  300,
				ImgHeight: 300,
			},
			expectedErr: "subID not provided",
		},
		{
			name: "GIVEN valid parameters THEN return the expected tracking URL",
			urlPattern: config.URLPattern{
				URL: "{product_url}",
				Queries: []config.Query{
					{Key: "click_param", Value: "test"},
					{Key: "id", Value: "{click_id_base64}"},
				},
			},
			params: Params{
				ProductURL: "https://product.com/item123",
				ClickID:    "abc123",
			},
			wantURL: "https://product.com/item123?click_param=test&id=YWJjMTIz",
		},
		{
			name: "GIVEN missing placeholders THEN return the expected tracking URL",
			urlPattern: config.URLPattern{
				URL: "{product_url}",
				Queries: []config.Query{
					{Key: "click_param", Value: "test"},
				},
			},
			params: Params{
				ProductURL: "https://product.com/item123",
				ClickID:    "abc123",
			},
			wantURL: "https://product.com/item123?click_param=test",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &Default{}
			gotURL, err := strategy.GenerateURL(tc.urlPattern, tc.params)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err.Error())
				require.Empty(t, gotURL)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantURL, gotURL)
			}
		})
	}
}
