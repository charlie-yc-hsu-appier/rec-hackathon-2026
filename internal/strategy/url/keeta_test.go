package url

import (
	"rec-vendor-api/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeetaRequest(t *testing.T) {
	keeta := &KeetaRequest{
		SceneType:    "FAKE-SCENE-TYPE",
		Ver:          "0",
		ChannelToken: "FAKE-TOKEN",
	}

	tt := []struct {
		name          string
		urlPattern    config.URLPattern
		params        Params
		wantURL       string
		wantParamsMap map[string]string
	}{
		{
			name: "GIVEN all params present THEN expect full URL with all params in dictionary order",
			urlPattern: config.URLPattern{
				URL: "https://host.keeta/api/recommend",
			},
			params: Params{
				ClickID:         "FAKE-CLICK-ID",
				ClientIP:        "127.0.0.1",
				KeetaCampaignID: "FAKE-KEETA-CAMPAIGN",
				Latitude:        "67.89",
				Longitude:       "123.45",
			},
			wantURL:       "https://host.keeta/api/recommend?bizType=bType&campaignId=FAKE-KEETA-CAMPAIGN&channelToken=FAKE-TOKEN&ip=127.0.0.1&lat=67.89&lon=123.45&reqId=FAKE-CLICK-ID&sceneType=FAKE-SCENE-TYPE&ver=0",
			wantParamsMap: map[string]string{},
		},
		{
			name: "GIVEN some params empty THEN expect URL with empty values in correct order",
			urlPattern: config.URLPattern{
				URL: "https://host.keeta/api/recommend",
			},
			params: Params{
				ClickID:         "",
				ClientIP:        "",
				KeetaCampaignID: "FAKE-KEETA-CAMPAIGN",
				Latitude:        "",
				Longitude:       "56.78",
			},
			wantURL:       "https://host.keeta/api/recommend?bizType=bType&campaignId=FAKE-KEETA-CAMPAIGN&channelToken=FAKE-TOKEN&ip&lat&lon=56.78&reqId&sceneType=FAKE-SCENE-TYPE&ver=0",
			wantParamsMap: map[string]string{},
		},
		{
			name: "GIVEN special characters in params THEN expect URL encoding is correct",
			urlPattern: config.URLPattern{
				URL: "https://host.keeta/api/recommend",
			},
			params: Params{
				ClickID:         "cl ick@id",
				ClientIP:        "127.0.0.1",
				KeetaCampaignID: "camp id",
				Latitude:        "12.34",
				Longitude:       "56.78",
			},
			wantURL:       "https://host.keeta/api/recommend?bizType=bType&campaignId=camp+id&channelToken=FAKE-TOKEN&ip=127.0.0.1&lat=12.34&lon=56.78&reqId=cl+ick%40id&sceneType=FAKE-SCENE-TYPE&ver=0",
			wantParamsMap: map[string]string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gotURL, gotParamsMap, err := keeta.GenerateURL(tc.urlPattern, tc.params)
			require.NoError(t, err)
			require.Equal(t, tc.wantURL, gotURL)
			require.Equal(t, tc.wantParamsMap, gotParamsMap)
		})
	}
}
