package body

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostFormat_GenerateBody(t *testing.T) {
	strategy := &PostFormat{}
	params := Params{
		UserID:    "testuser",
		ImgWidth:  100,
		ImgHeight: 200,
		SubID:     "testsubid",
		BundleID:  "com.test.app",
	}

	bodyObj := strategy.GenerateBody(params)
	bodyBytes, err := json.Marshal(bodyObj)
	assert.NoError(t, err)

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	assert.NoError(t, err)

	assert.Equal(t, "com.test.app", body["app"].(map[string]interface{})["bundleId"])
	assert.Equal(t, "testuser", body["device"].(map[string]interface{})["id"])
	assert.Equal(t, float64(0), body["device"].(map[string]interface{})["lmt"])
	assert.Equal(t, "100x200", body["imp"].(map[string]interface{})["imageSize"])
	assert.Equal(t, "testsubid", body["affiliate"].(map[string]interface{})["subId"])
}


