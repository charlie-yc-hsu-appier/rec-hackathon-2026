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
					{Key: "user_case_android", Value: "{user_id_case_by_os}"},
					{Key: "click_id", Value: "{click_id_base64}"},
					{Key: "site_domain", Value: "{web_host}"},
					{Key: "app_bundleId", Value: "{bundle_id}"},
					{Key: "imp_adType", Value: "{adtype}"},
					{Key: "partner_id", Value: "{partner_id}"},
				},
			},
			params: Params{
				UserID:    "TestUser",
				OS:        "aNDroid",
				ImgWidth:  200,
				ImgHeight: 100,
				ClickID:   "test-id",
				WebHost:   "http://example.com/query?param1=123&param2=456",
				BundleID:  "com.example.app",
				AdType:    1,
				PartnerID: "kakao_kr",
			},
			wantURL: "https://example.com/image?app_bundleId=com.example.app&click_id=dGVzdC1pZA&imp_adType=1&partner_id=kakao_kr&site_domain=http%3A%2F%2Fexample.com%2Fquery%3Fparam1%3D123%26param2%3D456&size=200x100&user=testuser&user_case_android=testuser",
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
			name: "GIVEN url with existing query parameters THEN return the expected URL with existing and parameters from URLPattern config",
			urlPattern: config.URLPattern{
				URL: "https://example.com/image/user/abc?imp_adType=1",
				Queries: []config.Query{
					{Key: "app_bundleId", Value: "com.example.app"},
				},
			},
			wantURL: "https://example.com/image/user/abc?app_bundleId=com.example.app&imp_adType=1",
		},
		{
			name: "GIVEN url with escapable characters (e.g. `space`) THEN return the expected escaped URL",
			urlPattern: config.URLPattern{
				URL: "https://example.com/image 2/user/abc",
			},
			params: Params{
				UserID:    "User",
				ImgWidth:  50,
				ImgHeight: 50,
			},
			wantURL: "https://example.com/image%202/user/abc",
		},

		{
			name: "GIVEN user_id_case_by_os macro with iOS OS THEN return uppercase user ID",
			urlPattern: config.URLPattern{
				URL: "https://api.example.com/test",
				Queries: []config.Query{
					{Key: "adid", Value: "{user_id_case_by_os}"},
				},
			},
			params: Params{
				UserID: "abc123DEF",
				OS:     "ios",
			},
			wantURL: "https://api.example.com/test?adid=ABC123DEF",
		},
		{
			name: "GIVEN user_id_case_by_os macro with empty UserID THEN return error",
			urlPattern: config.URLPattern{
				URL: "https://api.example.com/test",
				Queries: []config.Query{
					{Key: "adid", Value: "{user_id_case_by_os}"},
				},
			},
			params: Params{
				UserID: "",
				OS:     "android",
			},
			expectedErr: "UserID not provided",
		},
		{
			name: "GIVEN user_id_case_by_os macro with empty OS THEN return error",
			urlPattern: config.URLPattern{
				URL: "https://api.example.com/test",
				Queries: []config.Query{
					{Key: "adid", Value: "{user_id_case_by_os}"},
				},
			},
			params: Params{
				UserID: "TestUser123",
				OS:     "",
			},
			expectedErr: "OS not provided",
		},
		{
			name: "GIVEN user_id_case_by_os macro with unsupported OS THEN return error",
			urlPattern: config.URLPattern{
				URL: "https://api.example.com/test",
				Queries: []config.Query{
					{Key: "adid", Value: "{user_id_case_by_os}"},
				},
			},
			params: Params{
				UserID: "TestUser123",
				OS:     "web",
			},
			expectedErr: "unsupported OS: web (supported: android, ios)",
		},
		// tracking
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
		// keeta request
		{
			name: "GIVEN all params present THEN expect full URL with all params in dictionary order",
			urlPattern: config.URLPattern{
				URL: "https://host.keeta/api/recommend",
				Queries: []config.Query{
					{Key: "reqId", Value: "{click_id}"},
					{Key: "ip", Value: "{client_ip}"},
					{Key: "campaignId", Value: "{keeta_campaign_id}"},
					{Key: "lat", Value: "{latitude}"},
					{Key: "lon", Value: "{longitude}"},
					{Key: "sceneType", Value: "FAKE-SCENE-TYPE"},
					{Key: "ver", Value: "0"},
					{Key: "channelToken", Value: "FAKE-TOKEN"},
					{Key: "bizType", Value: "bType"},
				},
			},
			params: Params{
				ClickID:         "FAKE-CLICK-ID",
				ClientIP:        "127.0.0.1",
				KeetaCampaignID: "FAKE-KEETA-CAMPAIGN",
				Latitude:        "67.89",
				Longitude:       "123.45",
			},
			wantURL: "https://host.keeta/api/recommend?bizType=bType&campaignId=FAKE-KEETA-CAMPAIGN&channelToken=FAKE-TOKEN&ip=127.0.0.1&lat=67.89&lon=123.45&reqId=FAKE-CLICK-ID&sceneType=FAKE-SCENE-TYPE&ver=0",
		},
		{
			name: "GIVEN some params empty THEN expect URL with empty values in correct order",
			urlPattern: config.URLPattern{
				URL: "https://host.keeta/api/recommend",
				Queries: []config.Query{
					{Key: "reqId", Value: "{click_id}"},
					{Key: "ip", Value: "{client_ip}"},
					{Key: "campaignId", Value: "{keeta_campaign_id}"},
					{Key: "lat", Value: "{latitude}"},
					{Key: "lon", Value: "{longitude}"},
					{Key: "sceneType", Value: "FAKE-SCENE-TYPE"},
					{Key: "ver", Value: "0"},
					{Key: "channelToken", Value: "FAKE-TOKEN"},
					{Key: "bizType", Value: "bType"},
				},
			},
			params: Params{
				ClickID:         "",
				ClientIP:        "",
				KeetaCampaignID: "FAKE-KEETA-CAMPAIGN",
				Latitude:        "",
				Longitude:       "56.78",
			},
			wantURL: "https://host.keeta/api/recommend?bizType=bType&campaignId=FAKE-KEETA-CAMPAIGN&channelToken=FAKE-TOKEN&ip=&lat=&lon=56.78&reqId=&sceneType=FAKE-SCENE-TYPE&ver=0",
		},
		{
			name: "GIVEN special characters in params THEN expect URL encoding is correct",
			urlPattern: config.URLPattern{
				URL: "https://host.keeta/api/recommend",
				Queries: []config.Query{
					{Key: "reqId", Value: "{click_id}"},
					{Key: "ip", Value: "{client_ip}"},
					{Key: "campaignId", Value: "{keeta_campaign_id}"},
					{Key: "lat", Value: "{latitude}"},
					{Key: "lon", Value: "{longitude}"},
					{Key: "sceneType", Value: "FAKE-SCENE-TYPE"},
					{Key: "ver", Value: "0"},
					{Key: "channelToken", Value: "FAKE-TOKEN"},
					{Key: "bizType", Value: "bType"},
				},
			},
			params: Params{
				ClickID:         "cl ick@id",
				ClientIP:        "127.0.0.1",
				KeetaCampaignID: "camp id",
				Latitude:        "12.34",
				Longitude:       "56.78",
			},
			wantURL: "https://host.keeta/api/recommend?bizType=bType&campaignId=camp+id&channelToken=FAKE-TOKEN&ip=127.0.0.1&lat=12.34&lon=56.78&reqId=cl+ick%40id&sceneType=FAKE-SCENE-TYPE&ver=0",
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
