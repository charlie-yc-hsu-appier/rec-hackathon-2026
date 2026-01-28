package header

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestKeetaHeader(t *testing.T) {
	mockClock := NewMockClock(gomock.NewController(t))
	mockClock.EXPECT().getCurrentMilliTimestamp().Return("1734540677921")

	s := &KeetaHeader{
		SCaApp:    "FAKE-APP",
		SCaSecret: "FAKE-SECRET",
		Clock:     mockClock,
	}
	params := Params{
		UserID:     "FAKE-USER",
		RequestURL: "https://host.keeta/api/recommend?bizType=bType&campaignId=FAKE-KEETA-CAMPAIGN&channelToken=FAKE-TOKEN&ip=127.0.0.1&lat=67.89&lon=123.45&reqId=FAKE-CLICK-ID&sceneType=FAKE-SCENE-TYPE&ver=0",
	}
	wantHeaders := map[string]string{
		"Idfa":                   "FAKE-USER",
		"S-Ca-App":               "FAKE-APP",
		"S-Ca-Signature":         "19Bi+t4CYIidoRpfonYvXk22K/PUW2/AjAbe8VkeD7I=",
		"S-Ca-Signature-Headers": "Idfa,S-Ca-App,S-Ca-Timestamp",
		"S-Ca-Timestamp":         "1734540677921",
	}

	got := s.GenerateHeaders(params)
	require.Equal(t, wantHeaders, got)
}

// this test case implements the example signature described in Keeta Real-time DPA API document
func TestGenKeetaSignature(t *testing.T) {
	method := "GET"
	secret := "test"
	headers := map[string]string{
		"IDFA":           "asdfghjkl",
		"S-Ca-App":       "scaapp",
		"S-Ca-Timestamp": "1733906197000",
	}
	pathWithQuery := "/api/rl-recommend?bizType=bType&campaignId=123123&channelToken=cToken&ip=127.0.0.1&lat=22.324091&lon=114.254329&reqId=123321&sceneType=testType&ver=1"

	result := genKeetaSignature(secret, method, headers, pathWithQuery)
	wantedSignature := "z7eJqzEMV/cXeUAnpOZEl1IJvBDbgt0cLY3npp1k7kU="

	require.Equal(t, wantedSignature, result)
}
