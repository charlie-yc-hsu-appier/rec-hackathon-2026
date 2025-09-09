package requester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeeta(t *testing.T) {
	keeta := &Keeta{
		SceneType:    "FAKE-SCENE-TYPE",
		Ver:          "0",
		ChannelToken: "FAKE-TOKEN",
	}

	tt := []struct {
		name   string
		params Params
		want   string
	}{
		{
			name: "GIVEN all params present THEN expect full URL with all params in dictionary order",
			params: Params{
				ClickID:         "FAKE-CLICK-ID",
				ClientIP:        "127.0.0.1",
				KeetaCampaignID: "FAKE-KEETA-CAMPAIGN",
				Latitude:        "67.89",
				Longitude:       "123.45",
				RequestURL:      "https://host.keeta/api/recommend",
			},
			want: "https://host.keeta/api/recommend?bizType=bType&campaignId=FAKE-KEETA-CAMPAIGN&channelToken=FAKE-TOKEN&ip=127.0.0.1&lat=67.89&lon=123.45&reqId=FAKE-CLICK-ID&sceneType=FAKE-SCENE-TYPE&ver=0",
		},
		{
			name: "GIVEN some params empty THEN expect URL with empty values in correct order",
			params: Params{
				ClickID:         "",
				ClientIP:        "",
				KeetaCampaignID: "FAKE-KEETA-CAMPAIGN",
				Latitude:        "",
				Longitude:       "56.78",
				RequestURL:      "https://host.keeta/api/recommend",
			},
			want: "https://host.keeta/api/recommend?bizType=bType&campaignId=FAKE-KEETA-CAMPAIGN&channelToken=FAKE-TOKEN&ip&lat&lon=56.78&reqId&sceneType=FAKE-SCENE-TYPE&ver=0",
		},
		{
			name: "GIVEN special characters in params THEN expect URL encoding is correct",
			params: Params{
				ClickID:         "cl ick@id",
				ClientIP:        "127.0.0.1",
				KeetaCampaignID: "camp id",
				Latitude:        "12.34",
				Longitude:       "56.78",
				RequestURL:      "https://host.keeta/api/recommend",
			},
			want: "https://host.keeta/api/recommend?bizType=bType&campaignId=camp+id&channelToken=FAKE-TOKEN&ip=127.0.0.1&lat=12.34&lon=56.78&reqId=cl+ick%40id&sceneType=FAKE-SCENE-TYPE&ver=0",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := keeta.GenerateRequestURL(tc.params)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
